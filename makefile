run: build
	@./bot

build:
	@go build -o bot .

tidy:
	@go mod tidy

unzip:
	@unzip comments.txt.zip