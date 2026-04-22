package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

type ExportService struct {
	txnRepo     *repository.TransactionRepository
	accountRepo *repository.AccountRepository
}

func NewExportService(txnRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository) *ExportService {
	return &ExportService{txnRepo: txnRepo, accountRepo: accountRepo}
}

func (s *ExportService) ExportTransactions(userID uint64, startDate, endDate, format string) ([]byte, string, error) {
	filter := repository.TransactionFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	txns, err := s.txnRepo.List(userID, filter, 0, 10000)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportTransactionsCSV(txns)
	case "json":
		return s.exportTransactionsJSON(txns)
	default:
		return s.exportTransactionsCSV(txns)
	}
}

func (s *ExportService) ExportAccounts(userID uint64, format string) ([]byte, string, error) {
	accounts, err := s.accountRepo.List(userID, "", "", "name", 0, 1000)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportAccountsCSV(accounts)
	case "json":
		return s.exportAccountsJSON(accounts)
	default:
		return s.exportAccountsCSV(accounts)
	}
}

func (s *ExportService) exportTransactionsCSV(txns []model.Transaction) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"date", "type", "description", "amount", "source_account", "destination_account", "category", "tags", "notes"}
	writer.Write(header)

	for _, txn := range txns {
		sourceName := ""
		if txn.Source.ID != 0 {
			sourceName = txn.Source.Name
		}
		destName := ""
		if txn.Destination != nil && txn.Destination.ID != 0 {
			destName = txn.Destination.Name
		}
		categoryName := ""
		if txn.CategoryID != nil && txn.Category != nil {
			categoryName = txn.Category.Name
		}
		tags := ""
		for i, t := range txn.Tags {
			if i > 0 {
				tags += ","
			}
			tags += t.Name
		}

		row := []string{
			txn.Date.Format("2006-01-02"),
			string(txn.Type),
			txn.Description,
			txn.Amount.StringFixed(4),
			sourceName,
			destName,
			categoryName,
			tags,
			txn.Notes,
		}
		writer.Write(row)
	}

	writer.Flush()
	filename := fmt.Sprintf("transactions_%s.csv", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

func (s *ExportService) exportTransactionsJSON(txns []model.Transaction) ([]byte, string, error) {
	type txnExport struct {
		Date              string `json:"date"`
		Type              string `json:"type"`
		Description       string `json:"description"`
		Amount            string `json:"amount"`
		SourceAccount     string `json:"source_account"`
		DestinationAccount string `json:"destination_account"`
		Category          string `json:"category"`
		Tags              string `json:"tags"`
		Notes             string `json:"notes"`
	}

	items := make([]txnExport, 0, len(txns))
	for _, txn := range txns {
		sourceName := ""
		if txn.Source.ID != 0 {
			sourceName = txn.Source.Name
		}
		destName := ""
		if txn.Destination != nil && txn.Destination.ID != 0 {
			destName = txn.Destination.Name
		}
		categoryName := ""
		if txn.CategoryID != nil && txn.Category != nil {
			categoryName = txn.Category.Name
		}
		tags := ""
		for i, t := range txn.Tags {
			if i > 0 {
				tags += ","
			}
			tags += t.Name
		}

		items = append(items, txnExport{
			Date:              txn.Date.Format("2006-01-02"),
			Type:              string(txn.Type),
			Description:       txn.Description,
			Amount:            txn.Amount.StringFixed(4),
			SourceAccount:     sourceName,
			DestinationAccount: destName,
			Category:          categoryName,
			Tags:              tags,
			Notes:             txn.Notes,
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("transactions_%s.json", time.Now().Format("20060102"))
	return data, filename, nil
}

func (s *ExportService) exportAccountsCSV(accounts []model.Account) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"name", "type", "currency", "initial_balance", "current_balance", "is_virtual", "notes"}
	writer.Write(header)

	for _, a := range accounts {
		currencyCode := ""
		if a.Currency.ID != 0 {
			currencyCode = a.Currency.Code
		}
		row := []string{
			a.Name,
			string(a.Type),
			currencyCode,
			a.InitialBalance.StringFixed(4),
			a.CurrentBalance.StringFixed(4),
			fmt.Sprintf("%v", a.IsVirtual),
			a.Notes,
		}
		writer.Write(row)
	}

	writer.Flush()
	filename := fmt.Sprintf("accounts_%s.csv", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

func (s *ExportService) exportAccountsJSON(accounts []model.Account) ([]byte, string, error) {
	type accountExport struct {
		Name           string `json:"name"`
		Type           string `json:"type"`
		Currency       string `json:"currency"`
		InitialBalance string `json:"initial_balance"`
		CurrentBalance string `json:"current_balance"`
		IsVirtual      bool   `json:"is_virtual"`
		Notes          string `json:"notes"`
	}

	items := make([]accountExport, 0, len(accounts))
	for _, a := range accounts {
		currencyCode := ""
		if a.Currency.ID != 0 {
			currencyCode = a.Currency.Code
		}
		items = append(items, accountExport{
			Name:           a.Name,
			Type:           string(a.Type),
			Currency:       currencyCode,
			InitialBalance: a.InitialBalance.StringFixed(4),
			CurrentBalance: a.CurrentBalance.StringFixed(4),
			IsVirtual:      a.IsVirtual,
			Notes:          a.Notes,
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("accounts_%s.json", time.Now().Format("20060102"))
	return data, filename, nil
}
