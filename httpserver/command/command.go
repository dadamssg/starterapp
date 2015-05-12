package command

import (
	"github.com/dadamssg/jsonhttp"
	commandpkg "github.com/dadamssg/starterapp/app/command"
	"net/http"
)

func SendCommandErrors(w http.ResponseWriter, command commandpkg.Errorable) {
	var errors []jsonhttp.ResponseError
	for _, err := range command.GetErrors() {
		errors = append(errors, jsonhttp.ResponseError{Code: err.Code, Error: err.Err.Error()})
	}
	jsonhttp.SendErrors(w, errors)
}
