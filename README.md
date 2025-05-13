local: go run ./cmd/dev
producer: go run ./cmd/producer
consumer: go run ./cmd/consumer
kafka: docker compose up
ui: localhost:8080

### producer

- we currently have one producer. its a golang loop with a 10 second sleep timer that goes through all files in a specified folder and sends a message to the topic

### consumers and consumer groups

- we can spin up as many consumers as we want in the consumer group. but only one consumer will be picked for one partition
- e.g. we want to see more than one consumer from the consumer group running, we should set >1 partition for the topic

### Goal: See how consumers in the same group divide work using partitions.

- how:

  - create topic with 2 partitions
  - start 3 consumers for the same consumer group
  - run the producer

- note:

  - when you stop kafka, make sure you stop the consumers too. otherwise they don't get assigned to partitions properly somehow

TODO:

- create a new consumer group -> maybe notifications
- find out what else can kafka do that is worth exploring
- maybe try big loads
- how do i DDD? the encapsulation of items -> ask chatgpt
