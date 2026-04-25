// Package service 业务逻辑层，实现核心业务逻辑
// ImportService 数据导入业务逻辑，处理CSV文件解析和映射转换
package service

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

const importDir = "tmp/imports"

// ImportService 数据导入服务
// 负责处理CSV文件的上传、解析预览和导入执行
// 依赖txnService创建交易（确保余额更新、规则触发），依赖accountRepo获取默认账户，依赖db进行数据库操作
type ImportService struct {
	txnService *TransactionService                // 交易服务，用于创建导入的交易
	accountRepo *repository.AccountRepository     // 账户数据访问对象，用于获取默认源账户
	db         *gorm.DB                           // 数据库连接
}

// NewImportService 创建数据导入服务实例
func NewImportService(txnService *TransactionService, accountRepo *repository.AccountRepository, db *gorm.DB) *ImportService {
	return &ImportService{
		txnService: txnService,
		accountRepo: accountRepo,
		db:         db,
	}
}

// Upload 上传CSV文件到临时目录
// 生成唯一文件ID，保存到 tmp/imports/ 目录
// 参数：
//   - filename: 原始文件名
//   - reader: 文件内容读取器
// 返回：
//   - *response.ImportUploadResp: 上传结果（含文件ID）
//   - error: 错误信息
func (s *ImportService) Upload(filename string, reader io.Reader) (*response.ImportUploadResp, error) {
	// Ensure directory exists
	if err := os.MkdirAll(importDir, 0755); err != nil {
		return nil, errcode.ErrInternal
	}

	fileID := uuid.New().String()
	filePath := filepath.Join(importDir, fileID+".csv")

	f, err := os.Create(filePath)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.ImportUploadResp{FileID: fileID}, nil
}

// Parse 解析CSV文件并生成预览
// 根据字段映射验证每行数据的有效性（日期格式、金额格式）
// 参数：
//   - fileID: 上传时返回的文件ID
//   - mapping: 字段映射（CSV列名 -> 系统字段名，如 "交易日期" -> "date"）
// 返回：
//   - *response.ImportPreviewResp: 预览结果（含总行数、有效行数、无效行数、每行验证结果）
//   - error: 错误信息
func (s *ImportService) Parse(fileID string, mapping map[string]string) (*response.ImportPreviewResp, error) {
	filePath := filepath.Join(importDir, fileID+".csv")
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errcode.ErrImportFileInvalid
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errcode.ErrImportParseFail
	}

	if len(records) == 0 {
		return nil, errcode.ErrImportParseFail
	}

	// First row is header
	headers := records[0]
	rows := records[1:]

	preview := &response.ImportPreviewResp{
		Total: len(rows),
		Rows:  make([]response.ImportRowResp, 0, len(rows)),
	}

	for i, row := range rows {
		rowResp := response.ImportRowResp{
			Index: i + 1,
			Data:  make(map[string]string),
		}

		for j, val := range row {
			if j < len(headers) {
				rowResp.Data[headers[j]] = val
			}
		}

		// Validate mapped fields
		rowResp.IsValid = true
		rowResp.Errors = []string{}

		if dateField, ok := mapping["date"]; ok {
			if val, exists := rowResp.Data[dateField]; exists {
				if _, err := time.Parse("2006-01-02", val); err != nil {
					rowResp.IsValid = false
					rowResp.Errors = append(rowResp.Errors, "invalid date format")
				}
			} else {
				rowResp.IsValid = false
				rowResp.Errors = append(rowResp.Errors, "date field not found")
			}
		}

		if amountField, ok := mapping["amount"]; ok {
			if val, exists := rowResp.Data[amountField]; exists {
				if _, err := decimal.NewFromString(val); err != nil {
					rowResp.IsValid = false
					rowResp.Errors = append(rowResp.Errors, "invalid amount")
				}
			} else {
				rowResp.IsValid = false
				rowResp.Errors = append(rowResp.Errors, "amount field not found")
			}
		}

		if rowResp.IsValid {
			preview.Valid++
		} else {
			preview.Invalid++
		}

		preview.Rows = append(preview.Rows, rowResp)
	}

	return preview, nil
}

// Execute 执行数据导入
// 逐行解析CSV数据，根据字段映射提取日期、金额、描述等，通过txnService.Create创建交易
// 如果未指定交易类型，默认为withdrawal（支出）；如果未指定源账户，使用用户的第一个账户
// 参数：
//   - userID: 用户ID
//   - fileID: 上传时返回的文件ID
//   - mapping: 字段映射
// 返回：
//   - *response.ImportResultResp: 导入结果（含总数、成功数、失败数、跳过数）
//   - error: 错误信息
func (s *ImportService) Execute(userID uint64, fileID string, mapping map[string]string) (*response.ImportResultResp, error) {
	filePath := filepath.Join(importDir, fileID+".csv")
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errcode.ErrImportFileInvalid
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errcode.ErrImportParseFail
	}

	if len(records) == 0 {
		return nil, errcode.ErrImportParseFail
	}

	rows := records[1:]
	result := &response.ImportResultResp{
		Total: len(rows),
	}

	for _, row := range rows {
		headers := records[0]
		rowData := make(map[string]string)
		for j, val := range row {
			if j < len(headers) {
				rowData[headers[j]] = val
			}
		}

		// 步骤1：根据字段映射提取日期、金额、描述和交易类型
		dateStr := s.getMappedValue(rowData, mapping, "date")
		amountStr := s.getMappedValue(rowData, mapping, "amount")
		description := s.getMappedValue(rowData, mapping, "description")
		txnType := s.getMappedValue(rowData, mapping, "type")

		// 步骤2：跳过缺少日期或金额的行
		if dateStr == "" || amountStr == "" {
			result.Skipped++
			continue
		}

		// 步骤3：验证日期格式
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			result.Failed++
			continue
		}

		// 步骤4：验证金额格式和有效性（必须大于0）
		amountVal, err := decimal.NewFromString(amountStr)
		if err != nil || amountVal.LessThanOrEqual(decimal.Zero) {
			result.Failed++
			continue
		}

		// 步骤5：如果未指定交易类型，默认为支出（withdrawal）
		if txnType == "" {
			txnType = "withdrawal"
		}

		// 步骤6：获取用户的第一个账户作为默认源账户
		accounts, err := s.accountRepo.List(userID, "", "", "name", 0, 1)
		if err != nil || len(accounts) == 0 {
			result.Failed++
			continue
		}

		txn := &request.CreateTransactionReq{
			Type:        txnType,
			Date:        dateStr,
			Description: description,
			Amount:      amountStr,
			SourceID:    accounts[0].ID,
		}

		// 步骤7：通过txnService.Create创建交易（确保余额更新、规则触发和Webhook通知）
		if _, err := s.txnService.Create(userID, txn); err != nil {
			result.Failed++
			continue
		}

		result.Success++
	}

	return result, nil
}

// getMappedValue 根据字段映射从行数据中获取对应值
// 优先使用映射关系，如果映射不存在则尝试直接用字段名匹配
func (s *ImportService) getMappedValue(rowData map[string]string, mapping map[string]string, field string) string {
	if col, ok := mapping[field]; ok {
		return rowData[col]
	}
	// Try direct field name
	if val, ok := rowData[field]; ok {
		return val
	}
	return ""
}
