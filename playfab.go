package playfab

import "github.com/df-mc/go-playfab/internal"

// Result represents a successful response in PlayFab API.
// Make sure to specify the T generic type to whatever you want in the Data.
type Result[T any] internal.Result[T]

// Error represents an error included in the response body.
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
