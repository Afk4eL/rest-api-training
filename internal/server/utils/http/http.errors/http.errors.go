package httperrors

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound  = errors.New("user doesn`t exist")
	ErrWrongPassword = errors.New("incorrect password")
	ErrNotFound      = errors.New("record not found")
	ErrNoCookie      = errors.New("no required cookie")
	ErrBadRequest    = errors.New("bad request")
	ErrNoToken       = errors.New("no csrf token")

	ErrInternalServerError = errors.New("internal server error")
)

//TODO: Refactor to
//Error err
//Causes any

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
	Causes string `json:"causes,omitempty"`
}

func SetErrResponse(w http.ResponseWriter, r *http.Request, err error) {
	er := ParseError(err)

	render.Status(r, er.Status)
	render.JSON(w, r, er)
}

func NewBadRequestError(causes string) ErrorResponse {
	return ErrorResponse{
		Status: http.StatusBadRequest,
		Error:  ErrBadRequest.Error(),
		Causes: causes,
	}
}

func NewNotFoundError(causes string) ErrorResponse {
	return ErrorResponse{
		Status: http.StatusNotFound,
		Error:  ErrNotFound.Error(),
		Causes: causes,
	}
}

func NewUnauthoeizedError(causes string) ErrorResponse {
	return ErrorResponse{
		Status: http.StatusUnauthorized,
		Error:  ErrNoCookie.Error(),
		Causes: causes,
	}
}

func NewInternalServerError(causes string) ErrorResponse {
	return ErrorResponse{
		Status: http.StatusInternalServerError,
		Error:  ErrInternalServerError.Error(),
		Causes: causes,
	}
}

func ParseError(err error) ErrorResponse {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewNotFoundError(err.Error())
	case strings.Contains(strings.ToLower(err.Error()), "cookie"):
		return NewUnauthoeizedError(err.Error())
	case strings.Contains(strings.ToLower(err.Error()), "bcrypt"):
		return NewUnauthoeizedError(ErrWrongPassword.Error())
	case strings.Contains(strings.ToLower(err.Error()), "csrf"):
		return NewUnauthoeizedError(ErrNoToken.Error())
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		return ValidationError(err.(validator.ValidationErrors))

	default:
		log.Println(err)
		return NewInternalServerError(err.Error())
	}
}

func ValidationError(errs validator.ValidationErrors) ErrorResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return ErrorResponse{
		Status: http.StatusBadRequest,
		Error:  ErrBadRequest.Error(),
		Causes: strings.Join(errMsgs, ", "),
	}
}
