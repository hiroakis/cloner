
.PHONY: all clean deps build

all: clean build

deps:
	go mod tidy

build: deps
		go build

clean:
		rm -f cloner
