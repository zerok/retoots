all: bin/retoots

bin:
	mkdir -p bin

clean:
	rm -rf bin

bin/retoots: $(shell find . -name '*.go') go.mod bin
	cd cmd/retoots && go build -o ../../$@

test:
	go test ./... -v

.PHONY: clean all test
