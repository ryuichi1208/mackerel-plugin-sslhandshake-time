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

.PHONY:
test:
	go test -v ./...


.PHONY:
tidy:
	go mod tidy
	go mod verify
