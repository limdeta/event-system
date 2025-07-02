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
