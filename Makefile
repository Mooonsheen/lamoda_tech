build:
	go build ./app/cmd ;

run: build
	go run cmd ;

up: run
