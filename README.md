# go-rag-search

A RAG (Retrieval Augmented Generation) exploratory project on how to load data from an external source to augment an LLM result.
This is a PoC to load data from an online source like Confluence in order to help answer questions about internal documentation.

## Prerequisites

1. [Go 1.20+](https://go.dev/dl/)
2. [Docker](https://www.docker.com/)

### Local prerequisites

If you wish to make full use of the machine resources, the best option is to run the app locally, using a locally running Ollama server.
In that case, the prerequisites would be:

1. [Go 1.20+](https://go.dev/dl/)
2. [Ollama 0.5+](https://ollama.com/download) with llama 3.2
3. [Postgres 15+](https://www.postgresql.org/download/)

## How to run it

To run the complete suite, navigate to the repository root and run:

```sh
docker-compose up
```

That will spin up a Postgres container tweaked for pgvector store, an Ollama container with llama3.2, and will also run the application once the database is ready.

To run it locally, outside of the container, it is possible to use an already existing Postgres instance or run the database container and connect to it:

```sh
docker-compose up pgvector -d # only once

expose LLAMA_HOST=localhost
expose PG_HOST=localhost
expose PG_USER=postgres
expose PG_PASSWORD=postgres
expose PG_DB=postgres

go run . -load
```

## Loading data

The docker image for the app always loads the data from the remote source, but if you are running it locally, you just need to run with the `-local` flag for the first time. Next runs can run a lot faster without that, as the data would be already loaded into the database.
