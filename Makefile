PROJECTNAME=$(shell basename "$(PWD)")

.PHONY: run build compile test dep clean
all: test compile

run:
	go run cmd/go-auth/main.go

build: dep
	go build -o bin/$(PROJECTNAME) cmd/go-auth/main.go

compile: dep
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o bin/$(PROJECTNAME)-linux-x86_64 cmd/go-auth/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(PROJECTNAME)-win-x86_64.exe cmd/go-auth/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(PROJECTNAME)-darwin-x86_64 cmd/go-auth/main.go

test:
	go test ./...

dep:
	go mod tidy -v
	go mod vendor

clean:
	go mod tidy
	go clean
	rm -rf bin/*

