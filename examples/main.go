package main

import (
	"context"
	"log"
	"os"

	ragflow "github.com/kevinroleke/ragflow-go"
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

	client := ragflow.NewClient(apiKey, ragflow.WithBaseURL(baseURL), ragflow.WithUserPass("kevin@zerogon.consulting", "http://34.23.156.236/login"))

	ctx := context.Background()
	/*
	res, err := client.UploadDocumentFromBytes(ctx, "a931096e73e811f0aaf70242ac120006", "test.txt", []byte(""))
	if err != nil {
		log.Fatalln(err)
	}
	*/
	assistant, err := client.CreateAssistant(ctx, ragflow.CreateAssistantRequest{
		Name:        "Example 122231001004010101055412134",
		Description: "An assistant for testing the RAGFlow Go client",
		DatasetIDs:  []string{"a931096e73e811f0aaf70242ac120006"},
		LLMModel: "gpt-4o",
	})
	if err != nil {
		log.Fatalf("Error creating assistant: %v", err)
	}
	log.Printf("Created assistant: %s (ID: %s)\n", assistant.Name, assistant.ID)

	log.Println("Creating a session...")
	session, err := client.CreateSession(ctx, assistant.ID, ragflow.CreateSessionRequest{
		Name: "Example 234235410",
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}
	log.Printf("Created session: %s (ID: %s)\n", session.Name, session.ID)

	// Test streaming chat completion
	log.Println("Testing streaming chat completion...")
	respChan, errChan := client.CreateChatCompletionStream(ctx, ragflow.ChatCompletionRequest{
		Model: "09016520740011f080720242ac120006",//assistant.ID,
		Messages: []ragflow.ChatMessage{
			{
				Role:    "user",
				Content: "Hello! Can you tell me about yourself and what you can help with?",
			},
		},
		ConversationID: "090a4aa0740011f0bc920242ac120006",//session.ID,
	})

	var fullResponse string
	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				log.Printf("Complete streaming response: %s\n", fullResponse)
				goto done
			}
			if len(resp.Choices) > 0 {
				fullResponse += resp.Choices[0].Delta.Content
				log.Printf("Stream chunk: %s", resp.Choices[0].Delta.Content)
			}
		case err := <-errChan:
			if err != nil {
				log.Fatalf("Error in streaming chat completion: %v", err)
			}
		}
	}

done:
	log.Println("Streaming test completed successfully!")
	/*suc, err := client.SetAPIKey(ctx, ragflow.SetAPIKeyRequest{
		ApiKey: "123",
		FactoryName: "OpenAI",
	})
	if !suc || err != nil {
		log.Fatalf("Err: %v\n", err)
	}*/

	/*datasets, err := client.ListDatasets(ctx, &ragflow.ListDatasetsOptions{})
	if err != nil {
		log.Fatalln(err)
	}*/

	/*
	docs, err := client.ListDocuments(ctx, "a931096e73e811f0aaf70242ac120006", &ragflow.ListDocumentsOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(docs)
*/

	/*
	suc, err := client.AddLLM(ctx, ragflow.AddLLMRequest{
		ApiBase: "123",
		ModelName: "4",
		MaxTokens: 4,
		FactoryName: "Ollama",
		ModelType: "123",
	})
	if !suc || err != nil {
		log.Fatalf("Err: %v\n", err)
	}

	fmt.Println("Getting factories...")
	factories, err := client.GetFactories(ctx)
	if err != nil {
		log.Fatalf("Error getting factories: %v", err)
	}
	fmt.Printf("Found %d factories:\n", len(factories))
	for _, factory := range factories {
		fmt.Printf("- %s (Status: %s, Types: %v)\n", factory.Name, factory.Status, factory.ModelTypes)
	}
	fmt.Println()

	llm, err := client.GetMyLLMs(ctx)
	if err != nil {
		log.Fatalf("Error get llms: %v", err)
	}
	fmt.Println(llm)
*/

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
	*/
}
