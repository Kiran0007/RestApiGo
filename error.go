package main

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrorInternalServerError(err error) *HttpError {
	fmt.Println(err.Error())
	return &HttpError{"Internal Server Error", http.StatusInternalServerError}
}

func NewErrorParamMissing(param string) *HttpError {
	message := fmt.Sprintf("Required param \"%s\" missing", param)
	return &HttpError{message, http.StatusBadRequest}
}

func NewErrorParamEmpty(param string) *HttpError {
	message := fmt.Sprintf("param \"%s\" is empty", param)
	return &HttpError{message, http.StatusBadRequest}
}

func NewErrorBadRequestError(message string) *HttpError {
	return &HttpError{message, http.StatusBadRequest}
}

func NewErrorUnauthorizedError(message string) *HttpError {
	return &HttpError{message, http.StatusUnauthorized}
}

func NewErrorBodyParseError(err error) *HttpError {
	message := fmt.Sprintf("Error while parsing request body : \"%s\"", err.Error())
	return &HttpError{message, http.StatusBadRequest}
}
