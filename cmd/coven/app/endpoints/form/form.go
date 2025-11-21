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
			Path:        path.Join(FormPrefix, "card"),
			Methods:     []string{"POST", "GET", "PUT", "DELETE"},
			Secure:      true,
			HandlerFunc: cardHandleFunc,
		},
		{
			Path:        path.Join(FormPrefix, "image"),
			Methods:     []string{"POST", "GET", "DELETE"},
			Secure:      true,
			HandlerFunc: imgPoolUploadFileFunc,
		},
	}
}
