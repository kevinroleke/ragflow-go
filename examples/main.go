package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ragflow "github.com/staklabs/ragflow-go"
)

func main() {
	apiKey := os.Getenv("RAGFLOW_API_KEY")
	if apiKey == "" {
		log.Fatal("RAGFLOW_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("RAGFLOW_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1"
	}

	client := ragflow.NewClient(apiKey, ragflow.WithBaseURL(baseURL))

	ctx := context.Background()

	/*
	fmt.Println("Creating a dataset...")
	dataset, err := client.CreateDataset(ctx, ragflow.CreateDatasetRequest{
		Name:        "Example Dataset4",
		Description: "A dataset for testing the RAGFlow Go client",
	})
	if err != nil {
		log.Fatalf("Error creating dataset: %v", err)
	}
	fmt.Printf("Created dataset: %s (ID: %s)\n", dataset.Name, dataset.ID)
	*/
	var dataset ragflow.Dataset
	dataset.ID = "72c10e7c72f011f0aa4a0242ac120006"

	fmt.Println("Uploading a document...")
	doc, err := client.UploadDocumentFromBytes(ctx, dataset.ID, "example.txt", []byte("This is an example document for testing RAGFlow."))
	if err != nil {
		log.Fatalf("Error uploading document: %v", err)
	}
	fmt.Printf("Uploaded document: %s (ID: %s)\n", doc.Name, doc.ID)

	fmt.Println("Creating an assistant...")
	assistant, err := client.CreateAssistant(ctx, ragflow.CreateAssistantRequest{
		Name:        "Example 010000000055412134",
		Description: "An assistant for testing the RAGFlow Go client",
		DatasetIDs:  []string{dataset.ID},
		LLMModel:    "gpt-4o",
		Temperature: 0.7,
		MaxTokens:   1000,
	})
	if err != nil {
		log.Fatalf("Error creating assistant: %v", err)
	}
	fmt.Println(assistant)
	fmt.Printf("Created assistant: %s (ID: %s)\n", assistant.Name, assistant.ID)

	fmt.Println("Creating a session...")
	session, err := client.CreateSession(ctx, assistant.ID, ragflow.CreateSessionRequest{
		Name: "Example 234235410",
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}
	fmt.Printf("Created session: %s (ID: %s)\n", session.Name, session.ID)

	fmt.Println("Sending a chat completion request...")
	response, err := client.CreateChatCompletion(ctx, ragflow.ChatCompletionRequest{
		Model: assistant.ID,
		Messages: []ragflow.ChatMessage{
			{
				Role:    "user",
				Content: "Hello, can you help me understand what's in the uploaded document?",
			},
		},
		ConversationID: session.ID,
	})
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}
	fmt.Printf("Assistant response: %s\n", response.Choices[0].Message.Content)

	fmt.Println("Testing streaming chat completion...")
	respChan, errChan := client.CreateChatCompletionStream(ctx, ragflow.ChatCompletionRequest{
		Model: assistant.ID,
		Messages: []ragflow.ChatMessage{
			{
				Role:    "user",
				Content: "Can you provide a summary of the document content?",
			},
		},
		ConversationID: session.ID,
	})

	var fullResponse string
	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				fmt.Printf("Complete streaming response: %s\n", fullResponse)
				goto cleanup
			}
			if len(resp.Choices) > 0 {
				fullResponse += resp.Choices[0].Delta.Content
			}
		case err := <-errChan:
			if err != nil {
				log.Fatalf("Error in streaming chat completion: %v", err)
			}
		}
	}

cleanup:
	fmt.Println("Cleaning up resources...")

	if err := client.DeleteSession(ctx, assistant.ID, session.ID); err != nil {
		log.Printf("Error deleting session: %v", err)
	}

	if err := client.DeleteAssistant(ctx, assistant.ID); err != nil {
		log.Printf("Error deleting assistant: %v", err)
	}

	if err := client.DeleteDocument(ctx, dataset.ID, doc.ID); err != nil {
		log.Printf("Error deleting document: %v", err)
	}

	if err := client.DeleteDataset(ctx, dataset.ID); err != nil {
		log.Printf("Error deleting dataset: %v", err)
	}

	fmt.Println("Example completed successfully!")
}
