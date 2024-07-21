package errcode

const (
	ErrValidationErrors        = "VALIDATION_ERRORS"
	ErrInternalError           = "INTERNAL_ERROR"
	ErrInvalidDataSignature    = "INVALID_DATA_SIGNATURE"
	ErrInvalidRiskCheck        = "INVALID_RISK_CHECK"
	ErrUserNotFound            = "USER_NOT_FOUND"
	ErrUserAlreadyExist        = "USER_ALREADY_EXIST"
	ErrWalletAlreadyExist      = "WALLET_ALREADY_EXIST"
	ErrWalletNotFound          = "WALLET_NOT_FOUND"
	ErrTransactionAlreadyExist = "TRANSACTION_ALREADY_EXIST"
	ErrEkycAlreadyExist        = "EKYC_ALREADY_EXIST"
	ErrRiskConfigAlreadyExist  = "RISK_CONFIG_ALREADY_EXIST"
	ErrRiskConfigNotFound      = "RISK_CONFIG_NOT_FOUND"
	ErrWalletTypeNotDefined    = "WALLET_TYPE_NOT_DEFINED"
	ErrNotAllowAccess          = "NOT_ALLOW_ACCESS"
)
