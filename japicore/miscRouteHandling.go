package japicore

import (
	"fmt"
	"github.com/JackalLabs/jackalapi/jutils"
	"github.com/uptrace/bunrouter"
	"net/http"
)

var (
	Version = "v0.0.0"
	Module  = "Jackal API Core"
)

func Handler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

func RouteNotFoundHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		warning := fmt.Sprintf("%s is not an availble route", req.URL.Path)
		return jutils.ProcessCustomHttpError("MethodNotAllowedHandler", warning, 404, w)
	}
}

func MethodNotAllowedHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		warning := fmt.Sprintf("%s method not availble for \"%s\"", req.URL.Path, req.Method)
		return jutils.ProcessCustomHttpError("MethodNotAllowedHandler", warning, 405, w)
	}
}

func VersionHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		message := createJsonResponse("")
		condensedWriteJSON(w, message)
		return nil
	}
}
