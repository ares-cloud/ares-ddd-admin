GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
	ERROR_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *_error.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
	ERROR_PROTO_FILES=$(shell find . api -name *_error.proto)
endif


.PHONY: swag_admin
# swag_admin
swag_admin:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g main.go -d ./cmd/admin,./internal/base/interfaces/rest,./internal/monitoring/interfaces/rest --parseDependency --parseInternal -o docs/admin

.PHONY: swag_app
# swag_app
swag_app:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init   -d ./apps/app,./pkg/file --parseDependency --parseInternal -o ./apps/app/docs

.PHONY: wire_admin
# wire_admin
wire_admin:
	cd cmd/admin/  && wire

.PHONY: wire_app
# wire_app
wire_app:
	cd apps/app/ && wire


.PHONY: build_admin
# build
build_admin:
	make wire_admin;
	make swag_admin;
	mkdir -p bin/ && GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static" -s -w -X main.Version=$(VERSION)' -o ./bin/c10-admin ./apps/admin
	chmod +x ./bin/c10-admin

.PHONY: build_app
# build
build_app:
	rm -rf ./bin/app;
	make wire_app;
	make swag_app;
	mkdir -p bin/ && GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static" -s -w -X main.Version=$(VERSION)' -o ./bin/c10-app ./apps/app
	chmod +x ./bin/c10-app
