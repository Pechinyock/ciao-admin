package webui

import "io"

type WebUIElemet interface {
	Render(data any) (io.Writer, error)
}
