package playfab

import "github.com/df-mc/go-playfab/internal"

type Result[T any] internal.Result[T]

type Error = internal.Error

const (
	ErrorCodeEncryptionKeyMissing                  = 1290
	ErrorCodeEvaluationModePlayerCountExceeded     = 1490
	ErrorCodeExpiredXboxLiveToken                  = 1189
	ErrorCodeInvalidXboxLiveToken                  = 1188
	ErrorCodeRequestViewConstraintParamsNotAllowed = 1303
	ErrorCodeSignedRequestNotAllowed               = 1302
	ErrorCodeXboxInaccessible                      = 1339
	ErrorCodeXboxRejectedXSTSExchangeRequest       = 1343
	ErrorCodeXboxXASSExchangeFailure               = 1306
)

const (
	ErrorCodeDatabaseThroughputExceeded = 1113
	ErrorCodeItemNotFound               = 1047
	ErrorCodeNotImplemented             = 1515
)
