all: build

build:
	go build -o my_cut *.go


test: build
	sh script.sh


clean:
	rm -f *.txt test_input.csv my_cut