# RAGFlow Go Client

A comprehensive Go client library for interacting with the RAGFlow HTTP API. RAGFlow is a RAG (Retrieval-Augmented Generation) engine based on deep document understanding.

## Features

- **OpenAI-compatible API**: Chat completions with streaming support
- **Dataset Management**: Create, update, delete, and list datasets
- **Document Management**: Upload, download, and manage documents within datasets
- **Chunk Management**: Retrieve and update document chunks
- **Assistant Management**: Create and manage chat assistants
- **Session Management**: Handle conversation sessions
- **Agent Management**: Work with RAGFlow agents and DSL
- **Error Handling**: Comprehensive error types and handling
- **Streaming Support**: Real-time streaming for chat completions

## Installation

```bash
go get github.com/staklabs/ragflow-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    ragflow "github.com/staklabs/ragflow-go"
)

func main() {
    // Initialize client
    client := ragflow.NewClient(
        os.Getenv("RAGFLOW_API_KEY"),
        ragflow.WithBaseURL("http://your-ragflow-instance.com"),
    )

    ctx := context.Background()

    // Create a dataset
    dataset, err := client.CreateDataset(ctx, ragflow.CreateDatasetRequest{
        Name:        "My Dataset",
        Description: "A dataset for my documents",
        Language:    "English",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Upload a document
    doc, err := client.UploadDocument(ctx, dataset.ID, "/path/to/document.pdf")
    if err != nil {
        log.Fatal(err)
    }

    // Create an assistant
    assistant, err := client.CreateAssistant(ctx, ragflow.CreateAssistantRequest{
        Name:       "My Assistant",
        DatasetIDs: []string{dataset.ID},
        LLMModel:   "deepseek-chat",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Chat with the assistant
    response, err := client.CreateChatCompletion(ctx, ragflow.ChatCompletionRequest{
        Model: assistant.ID,
        Messages: []ragflow.ChatMessage{
            {Role: "user", Content: "What's in my document?"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Choices[0].Message.Content)
}
```

## Configuration

### Client Options

```go
client := ragflow.NewClient(apiKey,
    ragflow.WithBaseURL("http://localhost:9380"),  // Custom base URL
    ragflow.WithTimeout(60*time.Second),          // Custom timeout
    ragflow.WithHTTPClient(&http.Client{}),       // Custom HTTP client
)
```

### Environment Variables

- `RAGFLOW_API_KEY`: Your RAGFlow API key
- `RAGFLOW_BASE_URL`: Base URL of your RAGFlow instance (optional, defaults to `http://127.0.0.1`)

## API Reference

### Datasets

```go
// Create dataset
dataset, err := client.CreateDataset(ctx, ragflow.CreateDatasetRequest{
    Name:        "Dataset Name",
    Description: "Dataset description",
    Language:    "English",
})

// List datasets
datasets, err := client.ListDatasets(ctx, &ragflow.ListDatasetsOptions{
    Page:     1,
    PageSize: 10,
})

// Get dataset
dataset, err := client.GetDataset(ctx, datasetID)

// Update dataset
dataset, err := client.UpdateDataset(ctx, datasetID, ragflow.UpdateDatasetRequest{
    Description: "Updated description",
})

// Delete dataset
err := client.DeleteDataset(ctx, datasetID)
```

### Documents

```go
// Upload from file
doc, err := client.UploadDocument(ctx, datasetID, "/path/to/file.pdf")

// Upload from bytes
doc, err := client.UploadDocumentFromBytes(ctx, datasetID, "file.txt", []byte("content"))

// List documents
docs, err := client.ListDocuments(ctx, datasetID, &ragflow.ListDocumentsOptions{
    Keywords: "search term",
})

// Download document
data, err := client.DownloadDocument(ctx, datasetID, documentID)

// Delete document
err := client.DeleteDocument(ctx, datasetID, documentID)
```

### Assistants

```go
// Create assistant
assistant, err := client.CreateAssistant(ctx, ragflow.CreateAssistantRequest{
    Name:        "Assistant Name",
    DatasetIDs:  []string{datasetID},
    LLMModel:    "deepseek-chat",
    Temperature: 0.7,
    MaxTokens:   1000,
})

// List assistants
assistants, err := client.ListAssistants(ctx, nil)

// Update assistant
assistant, err := client.UpdateAssistant(ctx, assistantID, ragflow.UpdateAssistantRequest{
    Temperature: 0.8,
})
```

### Chat Completions

```go
// Standard chat completion
response, err := client.CreateChatCompletion(ctx, ragflow.ChatCompletionRequest{
    Model: assistantID,
    Messages: []ragflow.ChatMessage{
        {Role: "user", Content: "Hello!"},
    },
    Temperature: 0.7,
    MaxTokens:  1000,
})

// Streaming chat completion
respChan, errChan := client.CreateChatCompletionStream(ctx, ragflow.ChatCompletionRequest{
    Model: assistantID,
    Messages: []ragflow.ChatMessage{
        {Role: "user", Content: "Tell me a story"},
    },
})

for {
    select {
    case resp, ok := <-respChan:
        if !ok {
            // Stream finished
            return
        }
        fmt.Print(resp.Choices[0].Delta.Content)
    case err := <-errChan:
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

### Sessions

```go
// Create session
session, err := client.CreateSession(ctx, assistantID, ragflow.CreateSessionRequest{
    Name: "Session Name",
})

// List sessions
sessions, err := client.ListSessions(ctx, assistantID, nil)

// Update session
session, err := client.UpdateSession(ctx, assistantID, sessionID, ragflow.UpdateSessionRequest{
    Name: "New Name",
})
```

### Agents

```go
// Create agent
agent, err := client.CreateAgent(ctx, ragflow.CreateAgentRequest{
    Name: "Agent Name",
    DSL:  map[string]interface{}{"key": "value"},
})

// Run agent
response, err := client.RunAgent(ctx, agentID, "Hello", sessionID)

// Run agent with streaming
respChan, errChan := client.RunAgentStream(ctx, agentID, "Tell me a story", sessionID)
```

## Error Handling

The client provides structured error handling:

```go
if err != nil {
    if ragflow.IsErrorCode(err, ragflow.ErrorCodeUnauthorized) {
        log.Println("Invalid API key")
    } else if apiErr, ok := err.(*ragflow.APIError); ok {
        log.Printf("API error: %d - %s", apiErr.Code, apiErr.Message)
    } else {
        log.Printf("Other error: %v", err)
    }
}
```

## Examples

See the `/examples` directory for more comprehensive examples:

- `main.go`: Complete workflow example
- More examples coming soon

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.