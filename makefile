ADMIN_APP_NAME=ciao_admin
SYS_ADMIN_APP_NAME=ciao_sys_admin
VERSION=0.0.1
OUT_DIR=_bin
GIT_SHORT_SHA=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags="-X main.Version=$(VERSION) -X main.GitShorSha=$(GIT_SHORT_SHA)"

.PHONY: admin-win-x64 sys-admin-win-x64

admin-win-x64:
	set GOOS=windows&& set GOARCH=amd64&& go build $(LDFLAGS) -o ./$(OUT_DIR)/win-x64/$(ADMIN_APP_NAME).exe ./cmd/admin/main.go

sys-admin-win-x64:
	set GOOS=windows&& set GOARCH=amd64&& go build $(LDFLAGS) -o ./$(OUT_DIR)/win-x64/$(SYS_ADMIN_APP_NAME).exe ./cmd/sys_admin/main.go
