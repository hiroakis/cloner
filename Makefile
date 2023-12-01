BUILD=go build

DOCKER_BUILD=docker build
DOCKER_BUILD_OPTS=--no-cache

DOCKER_RMI=docker rmi -f

DESTDIR=./bin
TAG=koudaiii/cloner:latest

.PHONY: build
build:
	CGO_ENABLED=0 $(BUILD) -o $(DESTDIR)/cloner -ldflags "-s -w"

.PHONY: docker_image
docker_image: clean linux
	$(DOCKER_BUILD) -t $(TAG) . $(DOCKER_BUILD_OPTS)

.PHONY: clean
clean:
	$(DOCKER_RMI) -f $(TAG)
