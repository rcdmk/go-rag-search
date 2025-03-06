package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
)

var ConfluenceAPIKey = os.Getenv("CONFLUENCE_API_KEY")
var ConfluenceUsername = os.Getenv("CONFLUENCE_USERNAME")

var confluenceSpaces = []string{"CI"}

func main() {
	loadFromSource := flag.Bool("load", false, "loads source from external link to vector store")
	flag.Parse()

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Minute)

	llm, err := getLLM()
	if err != nil {
		log.Fatal(err)
	}

	store, err := getVectorStore(ctx, llm, *loadFromSource)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if *loadFromSource {
		err = loadDocs(ctx, confluenceSpaces, store)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = askAssistant(ctx, llm, store)
	if err != nil {
		log.Fatal(err)
	}
}
