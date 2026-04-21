package service

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

const importDir = "tmp/imports"

type ImportService struct {
	txnRepo    *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	db         *gorm.DB
}

func NewImportService(txnRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository, db *gorm.DB) *ImportService {
	return &ImportService{
		txnRepo:    txnRepo,
		accountRepo: accountRepo,
		db:         db,
	}
}

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

		// Map fields
		dateStr := s.getMappedValue(rowData, mapping, "date")
		amountStr := s.getMappedValue(rowData, mapping, "amount")
		description := s.getMappedValue(rowData, mapping, "description")
		txnType := s.getMappedValue(rowData, mapping, "type")

		if dateStr == "" || amountStr == "" {
			result.Skipped++
			continue
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			result.Failed++
			continue
		}

		amount, err := decimal.NewFromString(amountStr)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			result.Failed++
			continue
		}

		if txnType == "" {
			txnType = "withdrawal"
		}

		// Get source account (use first account as default)
		accounts, err := s.accountRepo.List(userID, "", "", "name", 0, 1)
		if err != nil || len(accounts) == 0 {
			result.Failed++
			continue
		}

		txn := &model.Transaction{
			UserID:      userID,
			Type:        model.TransactionType(txnType),
			Date:        date,
			Description: description,
			Amount:      amount,
			SourceID:    accounts[0].ID,
		}

		if err := s.txnRepo.Create(txn, nil); err != nil {
			result.Failed++
			continue
		}

		result.Success++
	}

	return result, nil
}

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
