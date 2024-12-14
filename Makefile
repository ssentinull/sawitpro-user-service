

.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate
	go mod tidy
	go mod vendor

run:
	go run cmd/*.go

test:
	go test -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	rm -rf ./generated
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

generate_mocks: 
	@rm -rf ./mocks
	@mockgen -destination=mocks/repository.go -source=repository/interfaces.go -package=mocks RepositoryInterface
	@mockgen -destination=mocks/usecase.go -source=usecase/interfaces.go -package=mocks UsecaseInterface
	@mockgen -destination=mocks/utils/crypt.go -source=utils/crypt.go -package=mocks CryptInterface
	@mockgen -destination=mocks/utils/auth.go -source=utils/auth.go -package=mocks AuthInterface
