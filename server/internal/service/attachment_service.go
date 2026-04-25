// Package service 业务逻辑层，实现核心业务逻辑
// AttachmentService 附件业务逻辑，处理文件上传、存储和下载
package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

const maxFileSize = 10 * 1024 * 1024 // 10MB

var allowedMimeTypes = map[string]bool{
	"image/jpeg":     true,
	"image/png":      true,
	"image/gif":      true,
	"application/pdf": true,
	"text/csv":       true,
	"application/x-ofx": true,
	"application/x-qif": true,
}

var allowedAttachableTypes = map[string]bool{
	"transaction": true,
	"account":     true,
	"bill":        true,
	"budget":      true,
	"piggy_bank":  true,
}

// AttachmentService 附件服务
// 负责处理文件上传、存储、下载和删除等业务逻辑
// 依赖attachmentRepo进行附件元数据访问，依赖storagePath指定文件存储路径
// 安全限制：文件大小上限10MB，仅允许特定MIME类型，仅允许关联到特定实体类型
type AttachmentService struct {
	attachmentRepo *repository.AttachmentRepository // 附件数据访问对象
	storagePath    string                           // 文件存储根路径
}

// NewAttachmentService 创建附件服务实例
// 参数：
//   - attachmentRepo: 附件数据访问对象
//   - storagePath: 文件存储根路径
// 返回：
//   - *AttachmentService: 附件服务实例
func NewAttachmentService(attachmentRepo *repository.AttachmentRepository, storagePath string) *AttachmentService {
	return &AttachmentService{attachmentRepo: attachmentRepo, storagePath: storagePath}
}

// Upload 上传附件
// 业务流程：
// 1. 验证关联实体类型是否允许（transaction/account/bill/budget/piggy_bank）
// 2. 验证文件大小不超过10MB
// 3. 验证MIME类型是否允许（图片/PDF/CSV/OFX/QIF）
// 4. 生成UUID文件名，按 用户ID/实体类型/实体ID 目录结构存储
// 5. 保存文件到磁盘，创建附件元数据记录
// 参数：
//   - userID: 用户ID
//   - attachableType: 关联实体类型（如"transaction"）
//   - attachableID: 关联实体ID
//   - fileHeader: 上传文件头信息
// 返回：
//   - *response.AttachmentResp: 附件信息
//   - error: 错误信息
func (s *AttachmentService) Upload(userID uint64, attachableType string, attachableID uint64, fileHeader *multipart.FileHeader) (*response.AttachmentResp, error) {
	// 步骤1：验证关联实体类型是否允许（transaction/account/bill/budget/piggy_bank）
	if !allowedAttachableTypes[attachableType] {
		return nil, errcode.ErrBadRequest
	}

	// 步骤2：验证文件大小不超过10MB
	if fileHeader.Size > maxFileSize {
		return nil, errcode.ErrBadRequest
	}

	// 步骤3：验证MIME类型是否允许（图片/PDF/CSV/OFX/QIF）
	mime := fileHeader.Header.Get("Content-Type")
	if !allowedMimeTypes[mime] {
		return nil, errcode.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, errcode.ErrInternal
	}
	defer file.Close()

	// 步骤4：生成UUID文件名，按 用户ID/实体类型/实体ID 目录结构存储
	ext := filepath.Ext(fileHeader.Filename)
	fileUUID := uuid.New().String() + ext
	relativePath := filepath.Join(fmt.Sprintf("%d", userID), attachableType, fmt.Sprintf("%d", attachableID), fileUUID)
	fullPath := filepath.Join(s.storagePath, relativePath)

	// 步骤5：保存文件到磁盘
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, errcode.ErrInternal
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		os.Remove(fullPath)
		return nil, errcode.ErrInternal
	}

	// 步骤6：创建附件元数据记录到数据库
	attachment := &model.Attachment{
		UserID:         userID,
		AttachableType: attachableType,
		AttachableID:   attachableID,
		Filename:       fileHeader.Filename,
		Mime:           mime,
		Size:           fileHeader.Size,
		Path:           relativePath,
	}

	if err := s.attachmentRepo.Create(attachment); err != nil {
		os.Remove(fullPath)
		return nil, errcode.ErrInternal
	}

	return s.toResp(attachment), nil
}

// Download 获取附件文件路径用于下载
// 参数：
//   - userID: 用户ID
//   - id: 附件ID
// 返回：
//   - string: 文件绝对路径
//   - error: 错误信息
func (s *AttachmentService) Download(userID, id uint64) (string, error) {
	attachment, err := s.attachmentRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errcode.ErrNotFound
		}
		return "", errcode.ErrInternal
	}

	fullPath := filepath.Join(s.storagePath, attachment.Path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errcode.ErrNotFound
	}

	return fullPath, nil
}

// AttachmentViewResult 附件查看结果，包含文件路径、文件名和MIME类型
type AttachmentViewResult struct {
	FilePath string
	Filename string
	Mime     string
}

// View 获取附件信息用于在线查看（含文件路径、文件名、MIME类型）
// 参数：
//   - userID: 用户ID
//   - id: 附件ID
// 返回：
//   - *AttachmentViewResult: 附件查看结果
//   - error: 错误信息
func (s *AttachmentService) View(userID, id uint64) (*AttachmentViewResult, error) {
	attachment, err := s.attachmentRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	fullPath := filepath.Join(s.storagePath, attachment.Path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errcode.ErrNotFound
	}

	return &AttachmentViewResult{
		FilePath: fullPath,
		Filename: attachment.Filename,
		Mime:     attachment.Mime,
	}, nil
}

// List 获取指定实体的附件列表
// 参数：
//   - userID: 用户ID
//   - attachableType: 关联实体类型
//   - attachableID: 关联实体ID
// 返回：
//   - []response.AttachmentResp: 附件列表
//   - error: 错误信息
func (s *AttachmentService) List(userID uint64, attachableType string, attachableID uint64) ([]response.AttachmentResp, error) {
	attachments, err := s.attachmentRepo.List(userID, attachableType, attachableID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.AttachmentResp, 0, len(attachments))
	for _, a := range attachments {
		items = append(items, *s.toResp(&a))
	}
	return items, nil
}

// Delete 删除附件
// 同时删除磁盘上的文件和数据库中的元数据记录
// 参数：
//   - userID: 用户ID
//   - id: 附件ID
// 返回：
//   - error: 错误信息
func (s *AttachmentService) Delete(userID, id uint64) error {
	attachment, err := s.attachmentRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	fullPath := filepath.Join(s.storagePath, attachment.Path)
	os.Remove(fullPath)

	return s.attachmentRepo.Delete(id, userID)
}

// toResp 将附件模型转换为响应对象
func (s *AttachmentService) toResp(a *model.Attachment) *response.AttachmentResp {
	return &response.AttachmentResp{
		ID:             a.ID,
		AttachableType: a.AttachableType,
		AttachableID:   a.AttachableID,
		Filename:       a.Filename,
		Mime:           a.Mime,
		Size:           a.Size,
		CreatedAt:      a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
