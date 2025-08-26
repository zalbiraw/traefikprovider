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
	docker-compose -f test/docker-compose.yml up -d

traefik-restart:
	@echo "Restarting Traefik and test services..."
	docker-compose -f test/docker-compose.yml restart

traefik-down:
	@echo "Stopping Traefik and test services..."
	docker-compose -f test/docker-compose.yml down

traefik-logs:
	@echo "Showing Traefik logs..."
	docker-compose -f test/docker-compose.yml logs -f traefik
