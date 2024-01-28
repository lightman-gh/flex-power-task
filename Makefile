.PHONY: generate
generate:
	go generate ./internal/dbschema

PHONY: docker-build
docker-build:
	docker build -t local/trades:local -f trades-api.Dockerfile .