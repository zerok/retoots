all: bin/retoots

bin:
	mkdir -p bin

clean:
	rm -rf bin

bin/retoots: $(shell find . -name '*.go') go.mod bin
	cd cmd/retoots && go build -o ../../$@

test:
	go test ./... -v

run-docs:
	docker run --rm \
		-v $(PWD):/data \
		-p 8000:8000 \
		zerok/mkdocs:latest serve -a 0.0.0.0:8000

.PHONY: clean all test run-docs
