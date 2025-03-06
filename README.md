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

### Docker Compose

Make sure you update the `CONFLUENCE_API_KEY` and `CONFLUENCE_USERNAME` environament variables in the `.env` file before running `docker-compose`.

To run the complete suite, navigate to the repository root and run:

```sh
docker-compose up
```

That will spin up a Postgres container tweaked for pgvector store, an Ollama container with llama3.2, and will also run the application once the database is ready.

### Running locally

To run it locally, outside of the container, it is possible to use an already existing Postgres instance or run the database container and connect to it:

```sh
docker-compose up pgvector -d # use this if you don't have a pgvector instance

export LLAMA_HOST=localhost
export PG_HOST=localhost
export PG_USER=postgres
export PG_PASSWORD=postgres
export PG_DB=postgres

export CONFLUENCE_API_KEY=your-super-secret-key
export CONFLUENCE_USERNAME=john.doe@acme.com

go run . -load
```

## Loading data

The docker image for the app always loads the data from the remote source, but if you are running it locally, you just need to run with the `-local` flag for the first time. Next runs can run a lot faster without that, as the data would be already loaded into the database.

## Environment Variables

Make sure to set the following environment variables before running the application:

- `LLAMA_HOST`: The host address for the Ollama server.
- `PG_HOST`: The host address for the Postgres database.
- `PG_USER`: The username for the Postgres database.
- `PG_PASSWORD`: The password for the Postgres database.
- `PG_DB`: The database name for the Postgres database.
- `CONFLUENCE_API_KEY`: Your Confluence/Atlassian API key/token.
- `CONFLUENCE_USERNAME`: Your Confluence/Atlassian username/email address

## Troubleshooting

If you encounter any issues, ensure that all prerequisites are installed and properly configured. Check the logs for any error messages and verify that all services are running correctly.

Most common issues are related to not enough memory allocated by the docker daemon (eg. by default Docker desktop allocates half of the available memory). Make sure you have at least 3.5 GB available for the Ollama container.

Inserting documents into the containeraised pgvector can be painfully slow, so you may have to adjust the timeout setting in the `main.go` file depending on the amount of data you load.
