# go-rag-search

A RAG (Retrieval Augmented Generation) exploratory project on how to load data from an external source to augment an LLM result.
This is a PoC to load data from an online source like Confluence in order to help answer questions about internal documentation.

## Prerequisites

1. [Go 1.20+](https://go.dev/dl/)
2. [Ollama 0.5+](https://ollama.com/download) with llama 3.2
3. [Docker](https://www.docker.com/)

## How to run it

To run, once Ollama is running, navigate to the repository root and run:

```sh
docker-compose up
```

That will spin up a Postgres container tweaked for pgvector store, and will also run the application once the database is ready.

To run it locally, outside of the container, it is possible to use an already existing Postgres instance or run the database container and connect to it:

```sh
docker-compose up pgvector -d # only once

expose LLAMA_HOST=localhost
expose PG_HOST=localhost
expose PG_USER=postgres
expose PG_PASSWORD=postgres
expose PG_DB=postgres

go run .
```
