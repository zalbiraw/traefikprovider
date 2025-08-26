.PHONY: lint test vendor clean traefik-up traefik-down traefik-logs preview-test

export GO111MODULE=on

default: lint test

lint:
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	yaegi test .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor

# Traefik Docker commands
traefik-up:
	@echo "Starting Traefik v3.5 with test services..."
	docker-compose -f docker-compose.test.yml up -d

traefik-down:
	@echo "Stopping Traefik and test services..."
	docker-compose -f docker-compose.test.yml down

traefik-logs:
	@echo "Showing Traefik logs..."
	docker-compose -f docker-compose.test.yml logs -f traefik

# Run shell-based integration checks using the preview docker-compose setup
preview-test:
	@echo "Running shell integration checks against cmd/preview/docker-compose.yml..."
	bash ./cmd/preview/test.sh
