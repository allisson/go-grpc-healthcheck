.PHONY: lint test

lint:
	if [ ! -f ./bin/golangci-lint ] ; \
	then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.18.0; \
	fi;
	./bin/golangci-lint run

test: lint
	go test -covermode=count -coverprofile=count.out -v ./...
