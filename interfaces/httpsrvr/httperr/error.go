package httperr

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/nononsensecode/base/errors"
	"github.com/nononsensecode/base/logs"
	"github.com/sirupsen/logrus"
)

func RespondWithApiError(err error, w http.ResponseWriter, r *http.Request) {
	apiError, ok := err.(errors.ApiError)
	if !ok {
		InternalServerError(int(errors.UnknownCode), string(errors.UnknownComponent), err, w, r)
		return
	}

	switch apiError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorised(int(apiError.Code()), string(apiError.Component()), apiError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(int(apiError.Code()), string(apiError.Component()), apiError, w, r)
	case errors.ErrorTypeNotFound:
		NotFound(int(apiError.Code()), string(apiError.Component()), apiError, w, r)
	default:
		InternalServerError(int(apiError.Code()), string(apiError.Component()), apiError, w, r)
	}
}

func InternalServerError(code int, component string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, code, component, w, r, "internal server error", http.StatusInternalServerError)
}

func Unauthorised(code int, component string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, code, component, w, r, err.Error(), http.StatusUnauthorized)
}

func BadRequest(code int, component string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, code, component, w, r, err.Error(), http.StatusBadRequest)
}

func NotFound(code int, component string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, code, component, w, r, err.Error(), http.StatusNotFound)
}

func httpRespondWithError(err error, code int, component string, w http.ResponseWriter, r *http.Request, logMsg string, status int) {
	logs.GetLogEntry(r).WithError(err).WithFields(logrus.Fields{
		"code":      code,
		"component": component,
	}).Warn(logMsg)
	resp := ErrorResponse{code, component, logMsg, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Code       int
	Component  string
	Msg        string
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
