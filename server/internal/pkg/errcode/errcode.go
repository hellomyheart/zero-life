// Package errcode 提供统一的业务错误码定义
// 错误码按模块分段管理：
//   - 0: 成功
//   - 400-499: HTTP标准错误
//   - 1xxxx: 认证错误
//   - 2xxxx: 账户错误
//   - 3xxxx: 交易错误
//   - 4xxxx: 分类错误
//   - 5xxxx: 标签错误
//   - 6xxxx: 预算错误
//   - 7xxxx: 账单错误
//   - 8xxxx: 规则错误
//   - 9xxxx: 导入错误
//   - 10xxxx: 定期交易错误
//   - 11xxxx: Webhook错误
//   - 12xxxx: 对象组错误
//   - 13xxxx: 交易链接错误
//   - 14xxxx: 偏好设置错误
//   - 15xxxx: 对账错误
//   - 16xxxx: MFA错误
package errcode

// Error 业务错误结构
// 包含错误码和错误消息，实现error接口
type Error struct {
	Code    int    `json:"code"`    // 错误码，用于程序判断
	Message string `json:"message"` // 错误消息，用于展示给用户
}

// Error 实现error接口，返回错误消息
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
	ErrTagTooDeep      = &Error{Code: 50003, Message: "tag hierarchy exceeds 2 levels"}

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

	// 定期交易/Recurrence错误 10xxxx
	ErrRecurrenceAmountInvalid = &Error{Code: 100001, Message: "recurrence amount must be greater than zero"}
	ErrRecurrenceInactive      = &Error{Code: 100002, Message: "recurrence is inactive"}
	ErrRecurringTransactionInvalid = &Error{Code: 100003, Message: "invalid recurring transaction"}
	ErrRecurringTransactionExpired = &Error{Code: 100004, Message: "recurring transaction has expired"}

	// Webhook错误 11xxxx
	ErrWebhookURLInvalid     = &Error{Code: 110001, Message: "invalid webhook URL"}
	ErrWebhookTriggerInvalid = &Error{Code: 110002, Message: "invalid webhook trigger"}
	ErrWebhookDeliveryFail   = &Error{Code: 110003, Message: "webhook delivery failed"}

	// 对象组错误 12xxxx
	ErrObjectGroupNameExists = &Error{Code: 120001, Message: "object group name already exists"}

	// 交易链接错误 13xxxx
	ErrTransactionLinkInvalid = &Error{Code: 130001, Message: "invalid transaction link"}
	ErrTransactionLinkExists  = &Error{Code: 130002, Message: "transaction link already exists"}

	// 偏好设置错误 14xxxx
	ErrPreferenceInvalid = &Error{Code: 140001, Message: "invalid preference key"}

	// 对账错误 15xxxx
	ErrReconciliationInvalid = &Error{Code: 150001, Message: "invalid reconciliation"}
	ErrReconciliationNotOpen = &Error{Code: 150002, Message: "reconciliation is not open"}

	// MFA错误 16xxxx
	ErrMFAAlreadyEnabled  = &Error{Code: 160001, Message: "MFA is already enabled"}
	ErrMFASetupNotFound   = &Error{Code: 160002, Message: "MFA setup not found or expired"}
	ErrMFAInvalidCode     = &Error{Code: 160003, Message: "invalid MFA code"}
	ErrMFANotEnabled      = &Error{Code: 160004, Message: "MFA is not enabled"}
	ErrMFAInvalidToken    = &Error{Code: 160005, Message: "invalid MFA token"}
)

// WithMessage 创建具有自定义消息的错误
// 保留原始错误码，替换错误消息
// 参数：
//   - e: 原始错误
//   - msg: 自定义错误消息
// 返回：
//   - *Error: 新的错误实例
func WithMessage(e *Error, msg string) *Error {
	return &Error{Code: e.Code, Message: msg}
}
