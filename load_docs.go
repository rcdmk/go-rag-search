package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func loadDocs(ctx context.Context, source string, store pgvector.Store) error {
	fmt.Println("loading data from", source)

	docs, err := getDocs(ctx, source)
	if err != nil {
		return err
	}

	fmt.Println("no. of documents to be loaded", len(docs))

	_, err = store.AddDocuments(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Println("data successfully loaded into vector store")

	return nil
}

func getDocs(ctx context.Context, source string) ([]schema.Document, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	docs, err := documentloaders.NewHTML(resp.Body).LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())
	if err != nil {
		return nil, err
	}

	return docs, nil
}
