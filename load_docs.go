package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rcdmk/go-rag-tutorial/internal/loaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func loadDocs(ctx context.Context, spaces []string, store pgvector.Store) error {
	fmt.Println("loading data from", spaces)

	docs, err := getDocs(ctx, spaces)
	if err != nil {
		return err
	}

	fmt.Println("no. of documents to be loaded", len(docs))

	batchSize := 100
	lastBatch := len(docs) % batchSize

	for i := 0; i < len(docs); i = i + batchSize {
		if i+batchSize > len(docs) {
			batchSize = lastBatch
		}
		_, err = store.AddDocuments(ctx, docs[i:i+batchSize])
		if err != nil {
			return err
		}
		fmt.Println("  loaded", i+batchSize, "docs")
	}

	fmt.Println("data successfully loaded into vector store")

	return nil
}

func getDocs(ctx context.Context, spaces []string) (docs []schema.Document, err error) {
	loader := loaders.NewConfluenceLoader("https://hootsuite.atlassian.net/wiki/", ConfluenceAPIKey, ConfluenceUsername, spaces)

	data, err := loader.Load(ctx)
	if err != nil {
		log.Fatalf("Failed to load documents: %v", err)
	}

	textSplitter := textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(500), textsplitter.WithChunkOverlap(0))

	docs, err = textsplitter.SplitDocuments(textSplitter, data)
	if err != nil {
		return docs, err
	}

	return docs, nil
}
