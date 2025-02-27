package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func main() {
	llamaHost := os.Getenv("LLAMA_HOST")
	if llamaHost == "" {
		llamaHost = "localhost"
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	llm, err := ollama.New(ollama.WithModel("llama3.2"), ollama.WithServerURL("http://"+llamaHost+":11434"))
	if err != nil {
		log.Fatal(err)
	}

	store, err := getVectorStore(ctx, llm)
	if err != nil {
		log.Fatal(err)
	}

	err = loadDocs(ctx, "https://support.atlassian.com/jira-software-cloud/docs/what-is-the-jira-family-of-products/", store)
	if err != nil {
		log.Fatal(err)
	}

	question := "Human: Who can see projects in Jira?\nAssistant:"

	fmt.Println()
	fmt.Println(question)

	numOfResults := 3

	_, err = chains.Run(
		ctx,
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, numOfResults),
		),
		question,
		chains.WithTemperature(0.8),
		chains.WithMaxTokens(2048),
		chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func getVectorStore(ctx context.Context, llm *ollama.LLM) (pgvector.Store, error) {

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
		pgvector.WithPreDeleteCollection(true),
		pgvector.WithConnectionURL(pgConnURL),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		return store, err
	}

	fmt.Println("vector store ready")

	return store, nil
}

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
