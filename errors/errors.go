package errors

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
	ErrorTypeNotFound       = ErrorType{"not-found"}
)

type Code int
type Component string

const (
	UnknownCode      Code      = 1000
	UnknownComponent Component = "Unknown"
)

type ApiError struct {
	code      Code
	component Component
	errorType ErrorType
	err       error
}

func (a ApiError) Error() string {
	return a.err.Error()
}

func (a ApiError) Code() Code {
	return a.code
}

func (a ApiError) Component() Component {
	return a.component
}

func (a ApiError) ErrorType() ErrorType {
	return a.errorType
}

func NewUnknownError(code Code, component Component, err error) ApiError {
	return ApiError{
		code:      code,
		component: component,
		errorType: ErrorTypeUnknown,
		err:       err,
	}
}

func NewAuthorizationError(code Code, component Component, err error) ApiError {
	return ApiError{
		code:      code,
		component: component,
		errorType: ErrorTypeAuthorization,
		err:       err,
	}
}

func NewIncorrectInputError(code Code, component Component, err error) ApiError {
	return ApiError{
		code:      code,
		component: component,
		errorType: ErrorTypeIncorrectInput,
		err:       err,
	}
}

func NewNotFoundError(code Code, component Component, err error) ApiError {
	return ApiError{
		code:      code,
		component: component,
		errorType: ErrorTypeNotFound,
		err:       err,
	}
}
