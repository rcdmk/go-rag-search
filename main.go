package main

import (
	"context"
	"flag"
	"log"
	"time"
)

func main() {
	loadFromSource := flag.Bool("load", false, "loads source from external link to vector store")
	flag.Parse()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	llm, err := getLLM()
	if err != nil {
		log.Fatal(err)
	}

	store, err := getVectorStore(ctx, llm)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if loadFromSource != nil && *loadFromSource {
		err = loadDocs(ctx, "https://support.atlassian.com/jira-software-cloud/docs/what-is-the-jira-family-of-products/", store)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = askAssistant(ctx, llm, store)
	if err != nil {
		log.Fatal(err)
	}
}
