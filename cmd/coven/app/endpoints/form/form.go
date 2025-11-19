package form

import (
	"ciao-admin/internal/server/router/endpoint"
	"path"
)

const FormPrefix = "/form"
const characterPathPart = "character"

func GetFormEndpoints() []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:        path.Join(FormPrefix, "login"),
			Methods:     []string{"POST"},
			Secure:      false,
			HandlerFunc: loginHandleFunc,
		},
		{
			Path:        path.Join(FormPrefix, characterPathPart),
			Methods:     []string{"POST", "GET", "PUT", "DELETE"},
			Secure:      true,
			HandlerFunc: characterHadleFunc,
		},
		{
			Path:        path.Join(FormPrefix, characterPathPart, "image"),
			Methods:     []string{"POST", "GET", "DELETE"},
			Secure:      true,
			HandlerFunc: imgPoolUploadFileFunc,
		},
	}
}
