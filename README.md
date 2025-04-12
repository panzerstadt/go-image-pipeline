producer: go run ./cmd/producer
consumer: go run ./cmd/consumer
kafka: docker compose up
ui: localhost:8080
