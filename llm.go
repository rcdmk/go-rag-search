package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func getLLM() (*ollama.LLM, error) {
	llamaHost := os.Getenv("LLAMA_HOST")
	if llamaHost == "" {
		llamaHost = "localhost"
	}

	return ollama.New(ollama.WithModel("llama3.2"), ollama.WithServerURL("http://"+llamaHost+":11434"))
}

func askAssistant(ctx context.Context, llm *ollama.LLM, store pgvector.Store) error {
	questionTemplate := "Human: %s\nAssistant:"
	question := fmt.Sprintf(questionTemplate, "Who can see projects in Jira?")

	fmt.Println()
	fmt.Println(question)

	numOfResults := 5

	_, err := chains.Run(
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
	fmt.Println()

	return err
}
