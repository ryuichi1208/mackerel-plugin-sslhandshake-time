NAME := check-sslhandshake-time

.PHONY:
build:
	go build -ldflags '-w -s' -o $(NAME) .

.PHONY:
clean:
	rm -f $(NAME)

.PHONY:
lint:
	go fmt ./...
	go vet ./...
