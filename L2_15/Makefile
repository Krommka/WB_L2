all: build

build:
	go build -o mysh main.go


test: build
	sh script.sh
	go test -v -race L2_15/minishell

vet:
	go vet -v ./...

lint:
	golint ./...

clean:
	rm -f *.txt mysh