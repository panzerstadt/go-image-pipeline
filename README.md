local: go run ./cmd/dev
producer: go run ./cmd/producer
consumer: go run ./cmd/consumer
kafka: docker compose up
ui: localhost:8080

TODO:

- test my undersatnding of consumers and how they work within a consumergroup
- test partitions

THEN:

- create a new consumer group -> maybe notifications
- find out what else can kafka do that is worth exploring
- maybe try big loads
- how do i DDD? the encapsulation of items -> ask chatgpt
