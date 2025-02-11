package apperror

import "fmt"

const (
	TokenExpired      = "TOKEN_EXPIRED"
	SecurityViolation = "SECURITY_ERROR"
)

var (
	ErrSessionDatabaseError            = &SessionError{Code: DatabaseError, Message: "database error"}
	ErrSessionDatabaseTransactionError = &SessionError{Code: DatabaseError, Message: "database transaction error"}
	ErrSessionInternalError            = &SessionError{Code: InternalError, Message: "internal error"}
	ErrSessionExpired                  = &SessionError{Code: TokenExpired, Message: "invalid session"}
	ErrSessionInvalidToken             = &SessionError{Code: TokenExpired, Message: "invalid token"}
	ErrSessionSecurityViolation        = &SessionError{Code: SecurityViolation, Message: "invalid session"}
	ErrSessionNotFound                 = &SessionError{Code: NotFound, Message: "session not found"}
)

type SessionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *SessionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *SessionError) Internal(err error) *SessionError {
	e.Err = err
	return e
}

func (e *SessionError) Msg(msg string) *SessionError {
	e.Message = msg
	return e
}
