package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nEnter your question: ")
		scanner.Scan()

		question := strings.TrimSpace(scanner.Text())
		if len(question) == 0 {
			continue
		}

		command := strings.ToLower(question)

		switch command {
		case "exit", "quit":
			return
		}

		fmt.Print("Answer: ")

		err = askAssistant(ctx, llm, store, question)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
	}
}
