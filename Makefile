.PHONY: generate
generate:
	go generate ./internal/dbschema

PHONY: docker-build
docker-build:
	docker build -t local/trades:local -f trades-api.Dockerfile .
	docker build -t local/poll:local -f poller.Dockerfile .
	docker build -t local/report:local -f report.Dockerfile .