# Even System sandbox

This project is a sandbox for learning how to use **Kafka** in a Go application by implementing event sourcing, some Domain-Driven Design (DDD) patterns, and practicing Test-Driven Development (TDD).  
My goal is to create a robust template for future projects. For now, it’s just an experimental playground.

## Getting Started

To run the project locally, use the following commands:

```sh
docker compose up -d       # Start Kafka, Zookeeper, and Kafdrop
go mod tidy                # Install Go dependencies
go run main.go             # Run the application
```

## Project Components

- **Kafdrop**  
  Web UI for monitoring Kafka.  
  Access it at [http://localhost:9000](http://localhost:9000).

- **Zookeeper**  
  Manages Kafka brokers, keeps configuration, and helps detect errors.

## Notes

There isn’t much functionality yet—this is just the first commit and a starting point for further development.


## Project Structure

A brief overview of the directory structure, following a simple Domain-Driven Design (DDD) approach:

---

### `/domain` — Domain Model
- Domain entities, business logic, and contracts (interfaces).
- Contains pure Go code with no dependencies on infrastructure or frameworks.

---

### `/application` — Application Layer
- Use cases and services.
- Coordinates domain logic and calls to the infrastructure layer.
- Acts as a bridge between the domain and the outside world.

---

### `/infrastructure` — Infrastructure Layer
- Implementations of interfaces from the domain layer.
- Integrations with third-party services, databases, Kafka, etc.
- Contains adapters and wrappers for external systems.

---

### `/interface` — Interface Layer
- All interactions with the outside world (REST API, WebSocket, GraphQL, CLI, etc).
- Handles input/output, request parsing, validation, and serialization.

---

```
/domain <---  /application  <---  /interface
   ^                  ^                    ^
   |                  |                    |
   +-----------+----- /infrastructure------+
```

---

> This structure helps to keep business logic isolated, makes the codebase maintainable, and allows easy testing and substitution of infrastructure components.

