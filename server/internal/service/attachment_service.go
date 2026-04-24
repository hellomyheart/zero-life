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

type AttachmentService struct {
	attachmentRepo *repository.AttachmentRepository
	storagePath    string
}

func NewAttachmentService(attachmentRepo *repository.AttachmentRepository, storagePath string) *AttachmentService {
	return &AttachmentService{attachmentRepo: attachmentRepo, storagePath: storagePath}
}

func (s *AttachmentService) Upload(userID uint64, attachableType string, attachableID uint64, fileHeader *multipart.FileHeader) (*response.AttachmentResp, error) {
	if !allowedAttachableTypes[attachableType] {
		return nil, errcode.ErrBadRequest
	}

	if fileHeader.Size > maxFileSize {
		return nil, errcode.ErrBadRequest
	}

	mime := fileHeader.Header.Get("Content-Type")
	if !allowedMimeTypes[mime] {
		return nil, errcode.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, errcode.ErrInternal
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	fileUUID := uuid.New().String() + ext
	relativePath := filepath.Join(fmt.Sprintf("%d", userID), attachableType, fmt.Sprintf("%d", attachableID), fileUUID)
	fullPath := filepath.Join(s.storagePath, relativePath)

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

type AttachmentViewResult struct {
	FilePath string
	Filename string
	Mime     string
}

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
