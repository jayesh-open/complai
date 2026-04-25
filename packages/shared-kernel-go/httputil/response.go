package httputil

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/complai/complai/packages/shared-kernel-go/errors"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := SuccessResponse{Data: data}
	_ = json.NewEncoder(w).Encode(resp)
}

func Error(w http.ResponseWriter, err error) {
	status, body := mapDomainError(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: body})
}

func mapDomainError(err error) (int, ErrorBody) {
	switch e := err.(type) {
	case *domainerrors.NotFoundError:
		return http.StatusNotFound, ErrorBody{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	case *domainerrors.ConflictError:
		return http.StatusConflict, ErrorBody{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	case *domainerrors.ValidationError:
		return http.StatusUnprocessableEntity, ErrorBody{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	case *domainerrors.AuthorizationError:
		return http.StatusForbidden, ErrorBody{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	case *domainerrors.ProviderError:
		return http.StatusBadGateway, ErrorBody{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	case *domainerrors.InternalError:
		return http.StatusInternalServerError, ErrorBody{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	default:
		return http.StatusInternalServerError, ErrorBody{
			Code:    "INTERNAL_ERROR",
			Message: "an unexpected error occurred",
		}
	}
}
