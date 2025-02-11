package apperror

import "fmt"

var (
	// Common
	ErrAuthDatabaseError = &AuthError{Code: DatabaseError, Message: "registration failed"} // 500
	ErrAuthInternalError = &AuthError{Code: InternalError, Message: "registration failed"}

	// Registration
	ErrAuthUserEmailExists = &AuthError{Code: DuplicateError, Message: "email already used"}  // 409
	ErrAuthPassword        = &AuthError{Code: BadRequest, Message: "invalid password format"} // 400
	ErrAuthEmailFormat     = &AuthError{Code: BadRequest, Message: "invalid email format"}    // 400
	ErrAuthEmailDomain     = &AuthError{Code: BadRequest, Message: "invalid email domain"}    // 400

	// Login
	ErrAuthUserNotFound = &AuthError{Code: NotFound, Message: "invalid username or password"} // 400

	// Refresh tokens
	// ErrUserInvalidRefreshToken = &SessionError{}
)

type AuthError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Err     error    `json:"-"`
	Details []string `json:"details,omitempty"`
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *AuthError) Internal(err error) *AuthError {
	e.Err = err
	return e
}

func (e *AuthError) Msg(msg string) *AuthError {
	e.Message = msg
	return e
}

func (e *AuthError) SetDetails(d []string) *AuthError {
	e.Details = d
	return e
}
