.PHONY: e2e
e2e:
	docker compose -f docker-compose.e2e.yaml up -d postgres flyway app
	docker compose -f docker-compose.e2e.yaml run --rm e2e
	code=$$?; docker compose -f docker-compose.e2e.yaml down -v; exit $$code

.PHONY: e2e-local
e2e-local:
	APP_URL="http://localhost:8080" \
	go test -race -v ./tests/e2e