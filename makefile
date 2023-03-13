.PHONY: all
TARGET := bot
GOENV := GOOS=linux GOARCH=amd64
GOMACENV := GOOS=darwin GOARCH=amd64
all:
	CGO_ENABLED=0 ${GOENV} go build -o ./bin/${TARGET}

clean:
	rm -rf trpc.log coverage.out
	rm -rf ${TARGET}
format:
	gofmt -w .
	goimports -w .
	golint ./...
test:
	go test --cover -gcflags=-l ./...

1:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    	go build -v -o ./bin/bot1

2:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    	go build -v -o ./bin/bot2

3:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    	go build -v -o ./bin/bot3

4:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    	go build -v -o ./bin/bot4

5:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    	go build -v -o ./bin/bot5