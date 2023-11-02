package japicore

import (
	"fmt"
	"net/http"

	"github.com/JackalLabs/jutils"
	"github.com/uptrace/bunrouter"
)

var (
	Version = "v0.0.0"
	Module  = "Jackal API Core"
)

func (j JApiCore) Handler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

func (j JApiCore) RouteNotFoundHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		warning := fmt.Sprintf("%s is not an availble route", req.URL.Path)
		return jutils.ProcessCustomHttpError("MethodNotAllowedHandler", warning, 404, w)
	}
}

func (j JApiCore) MethodNotAllowedHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		warning := fmt.Sprintf("%s method not availble for \"%s\"", req.URL.Path, req.Method)
		return jutils.ProcessCustomHttpError("MethodNotAllowedHandler", warning, 405, w)
	}
}

func (j JApiCore) VersionHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		message := createJsonResponse("", []string{})
		jutils.SimpleWriteJSON(w, message)
		return nil
	}
}
