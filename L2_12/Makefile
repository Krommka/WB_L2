all: build

build:
	go build -o my_grep *.go


test: build
	sh script.sh


clean:
	rm -f *.txt my_grep