package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func getVectorStore(ctx context.Context, llm *ollama.LLM, deleteCollection bool) (pgvector.Store, error) {
	host := os.Getenv("PG_HOST")
	if host == "" {
		log.Fatal("missing PG_HOST")
	}

	user := os.Getenv("PG_USER")
	if user == "" {
		log.Fatal("missing PG_USER")
	}

	password := os.Getenv("PG_PASSWORD")
	if password == "" {
		log.Fatal("missing PG_PASSWORD")
	}

	dbName := os.Getenv("PG_DB")
	if dbName == "" {
		log.Fatal("missing PG_DB")
	}

	connURLFormat := "postgres://%s:%s@%s:5432/%s?sslmode=disable"

	pgConnURL := fmt.Sprintf(connURLFormat, user, url.QueryEscape(password), host, dbName)

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}
	store, err := pgvector.New(
		ctx,
		pgvector.WithPreDeleteCollection(deleteCollection),
		pgvector.WithConnectionURL(pgConnURL),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		return store, err
	}

	return store, nil
}
