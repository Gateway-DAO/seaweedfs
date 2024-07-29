.PHONY: image
images:
	@./docker_build.sh -t local

images-testnet:
	@./docker_build.sh -t testnet

.PHONY: kafka kafka-logs kafka-stop
KAFKA_COMPOSE_FILES := -f docker-compose.kafka.yml
kafka: docker-compose.kafka.yml
	docker compose $(KAFKA_COMPOSE_FILES) up -d --build

kafka-logs: docker-compose.kafka.yml
	docker compose $(KAFKA_COMPOSE_FILES) logs -f

kafka-stop:
	docker compose $(KAFKA_COMPOSE_FILES) down -v

benchmark:
	docker compose -f docker-compose.benchmark.yml up -d --build benchmark
