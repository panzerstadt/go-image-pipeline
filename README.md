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

### producing messages with custom partition keys

```
msg := &sarama.ProducerMessage{
    Topic: "your-topic",
    Key:   sarama.StringEncoder("user123"), // This determines the partition
    Value: sarama.StringEncoder("some image job payload"),
}

partition, offset, err := producer.SendMessage(msg)
if err != nil {
    log.Fatalf("failed to send message: %v", err)
}

log.Printf("Message sent to partition %d at offset %d", partition, offset)
```

#### create a new consumer group -> maybe notifications

Try this:
• Reuse the same topic
• Spin up a new consumer group (e.g., notification-worker) and analytics worker
• It should receive all messages from offset 0, even if your image-processing group already did
• This simulates a fan-out architecture

TODO:

#### find out what else can kafka do that is worth exploring

Here are some spicy ideas:
• Dead Letter Queues (DLQs) – catch poison messages that keep failing
• Delayed delivery via scheduled retries (not native, but doable with workarounds)
• Stream processing using something like Kafka Streams or Faust (Python)
• Log compaction vs normal retention
• Message headers for metadata transport

#### maybe try big loads

Goal: Stress test for throughput and backpressure

Try this:
• Bulk-produce 10,000+ messages quickly
• Monitor how your consumers handle lag (check ConsumerLag via Kafka UI or metrics)
• Optional: add artificial delay in consumer to simulate heavy processing

#### how do i DDD? the encapsulation of items -> ask chatgpt

This is a big one. At a high level:
• Entities: Have identity (e.g., ImageJob) and mutable state
• Value Objects: No identity, just data (e.g., Resolution, Dimensions)
• Aggregates: Cluster of objects treated as a unit (e.g., an ImageProcessingWorkflow)
• Repositories: Abstract away data access
• Services: Contain business logic that spans multiple aggregates

For your image pipeline, you might:
• Turn a ResizeJob into an aggregate
• Have JobStatus as a value object
• Use a repository interface to load/save job data (e.g., from Kafka state, S3, or DB)

💡 Breakdown
• imagejob: your domain model—aggregates, value objects, behavior
• processing: application layer—services that orchestrate the domain
• repository: infrastructure layer—Kafka or FS adapters for persistence
• cmd: executables (Go’s standard approach)
• internal: keeps implementation details encapsulated, Go idiomatic

```
go-image-pipeline/
├── cmd/
│   ├── consumer/               # Entry point for consumer service
│   ├── producer/               # Entry point for producer
│   └── dev/                    # Dev/testing runner
├── internal/
│   ├── imagejob/               # Aggregate: ResizeJob, JobState transitions
│   │   ├── job.go              # Entity definition (ID, state, etc.)
│   │   └── status.go           # Value object: JobStatus, e.g. Pending, Done
│   ├── processing/             # Service: contains business logic (e.g., ResizeImage)
│   │   └── service.go
│   ├── repository/             # Interfaces and adapters for persistence
│   │   ├── interface.go        # Repository interfaces
│   │   ├── kafka_repo.go       # Implementation: writes/reads from Kafka
│   │   └── fs_repo.go          # (Optional) Implementation: reads from local FS
│   └── shared/                 # Utilities, domain-wide constants, error types
│       └── logger.go
├── proto/                      # Protobuf definitions
├── ui/                         # Optional: frontend or dashboard
├── scripts/                    # Local tooling, CLI commands
├── README.md
└── go.mod
```

#### DDD

✅ 3. Separate domain logic from coordination + plumbing
• Domain logic: rules, decisions, calculations
• Application logic: orchestration, calling services, scheduling
• Infrastructure: Kafka, DB, HTTP, Docker, etc.

For instance:
• The rule “only resize images wider than 800px” → domain
• The act of pulling that job from Kafka → infra

#### domain: Image Processing Workflow Management

Which means your domain model includes:
• Jobs (entities)
• Status/state transitions (value objects or enums)
• Dimensions, file types (value objects)
• Processing decisions and constraints (business rules)

Everything else—how you send a job through Kafka, where you store images, even the UI—is supporting infrastructure.

### DDD

- if the domain is about image processing workflow mgmt
- then the domain object is not the jpeg, but the imageJob

📚 Core DDD Vocabulary

These are terms you’ll see in DDD-oriented systems:

🧱 Domain Layer
• Entity – has identity (e.g. ImageJob, User)
• Value Object – no identity, immutable (e.g. Dimensions, Resolution)
• Aggregate – cluster of entities/VOs with rules (e.g. ResizeJob with status and constraints)
• Domain Service – stateless service with domain logic across multiple entities (e.g. JobScheduler)
• Domain Event – something that happened in the domain (e.g. ImageResizedEvent)

🧩 Application Layer
• Application Service / Use Case – orchestrates domain behavior (e.g. ProcessImageUseCase)
• Command – request to perform an action (e.g. ResizeImageCommand)
• Query – read model request
• DTO – data transfer object (input/output, not a domain model)

📦 Infrastructure Layer
• Repository – access to domain objects (e.g. Kafka, FS, DB)
• Adapter – a plug-in that conforms to interface (e.g. KafkaProducerAdapter)
• Gateway – external service abstraction (e.g. S3Client)

### MVC <-> DDD

| Common Term | DDD Equivalent (ish)                           |
| ----------- | ---------------------------------------------- |
| Model       | Often includes Entities + Value Objects        |
| Repository  | Same in DDD                                    |
| Service     | Could be Domain Service or Application Service |
| Controller  | UI/Application Layer Entry Point               |

### 🔍 Tips to Recognize DDD-style Structure

Look for:
• Entity + Value Object split
• Interfaces for repositories
• UseCase or CommandHandler classes/files
• Language like Domain, Aggregate, Event
• Folders like domain, application, infrastructure, interfaces
