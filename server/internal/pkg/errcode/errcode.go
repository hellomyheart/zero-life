package errcode

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

var (
	// 通用错误
	ErrSuccess       = &Error{Code: 0, Message: "success"}
	ErrBadRequest    = &Error{Code: 400, Message: "bad request"}
	ErrUnauthorized  = &Error{Code: 401, Message: "unauthorized"}
	ErrForbidden     = &Error{Code: 403, Message: "forbidden"}
	ErrNotFound      = &Error{Code: 404, Message: "not found"}
	ErrInternal      = &Error{Code: 500, Message: "internal server error"}
	ErrTooManyReqs   = &Error{Code: 429, Message: "too many requests"}

	// 认证错误 1xxxx
	ErrEmailExists      = &Error{Code: 10001, Message: "email already exists"}
	ErrInvalidCredential = &Error{Code: 10002, Message: "invalid email or password"}
	ErrAccountLocked    = &Error{Code: 10003, Message: "account is locked, try again later"}
	ErrInvalidToken     = &Error{Code: 10004, Message: "invalid or expired token"}
	ErrPasswordTooShort = &Error{Code: 10005, Message: "password must be at least 8 characters"}
	ErrOldPasswordWrong = &Error{Code: 10006, Message: "old password is incorrect"}

	// 账户错误 2xxxx
	ErrAccountNameExists = &Error{Code: 20001, Message: "account name already exists under this type"}
	ErrAccountTypeChange = &Error{Code: 20002, Message: "account type cannot be changed"}
	ErrAccountHasTxns    = &Error{Code: 20003, Message: "account has associated transactions"}
	ErrCurrencyInUse     = &Error{Code: 20004, Message: "currency is in use and cannot be disabled"}
	ErrDefaultCurrency   = &Error{Code: 20005, Message: "default currency cannot be disabled"}

	// 交易错误 3xxxx
	ErrInvalidAmount     = &Error{Code: 30001, Message: "amount must be greater than zero"}
	ErrSplitAmountMismatch = &Error{Code: 30002, Message: "split amounts must equal total amount"}
	ErrInvalidTxnType    = &Error{Code: 30003, Message: "invalid transaction type"}

	// 分类错误 4xxxx
	ErrCategoryNameExists = &Error{Code: 40001, Message: "category name already exists"}
	ErrCategoryTooDeep   = &Error{Code: 40002, Message: "category hierarchy exceeds 2 levels"}

	// 标签错误 5xxxx
	ErrTagNameExists   = &Error{Code: 50001, Message: "tag name already exists"}
	ErrTagAlreadyAdded = &Error{Code: 50002, Message: "tag already added to transaction"}

	// 预算错误 6xxxx
	ErrBudgetAmountInvalid = &Error{Code: 60001, Message: "budget amount must be greater than zero"}

	// 账单错误 7xxxx
	ErrBillAmountInvalid = &Error{Code: 70001, Message: "bill amount must be greater than zero"}

	// 规则错误 8xxxx
	ErrRuleConditionInvalid = &Error{Code: 80001, Message: "invalid rule condition"}
	ErrRuleActionInvalid    = &Error{Code: 80002, Message: "invalid rule action"}
	ErrRuleGroupHasRules    = &Error{Code: 80003, Message: "rule group has associated rules"}
	ErrRuleGroupInactive    = &Error{Code: 80004, Message: "rule group is inactive"}

	// 导入错误 9xxxx
	ErrImportFileInvalid = &Error{Code: 90001, Message: "invalid import file"}
	ErrImportParseFail   = &Error{Code: 90002, Message: "failed to parse import file"}

	// 循环交易错误 10xxxx
	ErrRecurrenceAmountInvalid = &Error{Code: 100001, Message: "recurrence amount must be greater than zero"}
	ErrRecurrenceInactive      = &Error{Code: 100002, Message: "recurrence is inactive"}

	// Webhook错误 11xxxx
	ErrWebhookURLInvalid = &Error{Code: 110001, Message: "webhook URL must start with https://"}
	ErrWebhookNotFound   = &Error{Code: 110002, Message: "webhook not found"}

	// 对账错误 12xxxx
	ErrReconciliationInvalidDate = &Error{Code: 120001, Message: "reconciliation date range is invalid"}
	ErrReconciliationBalanceInvalid = &Error{Code: 120002, Message: "submitted balance is invalid"}

	// 批量操作错误 13xxxx
	ErrBulkEmptyIDs    = &Error{Code: 130001, Message: "transaction IDs list cannot be empty"}
	ErrBulkConvertFail = &Error{Code: 130002, Message: "failed to convert transaction type"}

	// 交易链接错误 14xxxx
	ErrTxnLinkDuplicate = &Error{Code: 140001, Message: "transaction link already exists"}
	ErrLinkTypeNameExists = &Error{Code: 140002, Message: "link type name already exists"}

	// 偏好设置错误 15xxxx
	ErrPreferenceNameInvalid = &Error{Code: 150001, Message: "preference name is invalid"}
)

func WithMessage(e *Error, msg string) *Error {
	return &Error{Code: e.Code, Message: msg}
}
