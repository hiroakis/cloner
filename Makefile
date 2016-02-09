
.PHONY: all clean deps build

all: clean build

deps:
		go get -d -v ./...
		go get golang.org/x/oauth2
		go get github.com/google/go-github/github

build: deps
		go build

clean:
		rm -f cloner
