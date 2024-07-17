.PHONY: default

default:
	rm -fv bin/*
	GOOS=darwin GOARCH=arm64 go build -o bin/ebsc-arm64-0.1.0 -ldflags "-X main.ver=0.1.0 -X 'main.build=`date +%Y%m%d`'" ./cmd/ebsc
	GOOS=darwin GOARCH=amd64 go build -o bin/ebsc-amd64-0.1.0 -ldflags "-X main.ver=0.1.0 -X 'main.build=`date +%Y%m%d`'" ./cmd/ebsc
