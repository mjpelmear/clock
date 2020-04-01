
GO = go

format:
	$(GO) format

vet:
	$(GO) vet

test:
	$(GO) test

coverage:
	$(GO) test -covermode=count -coverprofile=coverage.out && \
	$(GO) tool cover -func=coverage.out && \
	$(GO) tool cover -html=coverage.out -o coverage.html

