package form

import (
	"ciao-admin/internal/server/router/endpoint"
	"path"
)

const FormPrefix = "/form"

func GetFormEndpoints() []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:        path.Join(FormPrefix, "login"),
			Methods:     []string{"POST"},
			Secure:      false,
			HandlerFunc: loginHandleFunc,
		},
		{
			Path:        path.Join(FormPrefix, "character"),
			Methods:     []string{"POST", "GET", "PUT", "DELETE"},
			Secure:      true,
			HandlerFunc: characterHadleFunc,
		},
	}
}
