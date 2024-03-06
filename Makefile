.PHONY: build
# build
build:
	mkdir -p bin/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./bin/app ./cmd

.PHONY: docker
# docker
docker:
	make build && \
	docker build -t smalls0098/proxydown:latest . && \
	docker push smalls0098/proxydown:latest