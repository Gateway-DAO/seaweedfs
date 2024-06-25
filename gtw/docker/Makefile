.PHONY: image
images:
	docker build -t gateway/edv:latest -f edv.Dockerfile ../..
	docker build -t gateway/master:latest -f master.Dockerfile ../..

.PHONY: kafka kafka-logs kafka-stop
KAFKA_COMPOSE_FILES := -f docker-compose.kafka.yml
kafka: docker-compose.kafka.yml
	docker compose $(KAFKA_COMPOSE_FILES) up -d --build

kafka-logs: docker-compose.kafka.yml
	docker compose $(KAFKA_COMPOSE_FILES) logs -f

kafka-stop:
	docker compose $(KAFKA_COMPOSE_FILES) down -v
