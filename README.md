# openai-go

A Go client library for the [OpenAI API](https://platform.openai.com/docs/api-reference).

This is a fork of [openai/openai-go](https://github.com/openai/openai-go) with additional features and improvements.

## Installation

```bash
go get github.com/your-org/openai-go
```

## Requirements

- Go 1.21 or later
- An OpenAI API key

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/your-org/openai-go"
	"github.com/your-org/openai-go/option"
)

func main() {
	client := openai.NewClient(
		option.WithAPIKey("your-api-key"), // defaults to OPENAI_API_KEY env var
	)

	chatCompletion, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("What is the meaning of life?"),
			}),
			Model: openai.F(openai.ChatModelGPT4o),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(chatCompletion.Choices[0].Message.Content)
}
```

## Features

- **Chat Completions** — Generate text with GPT-4o, GPT-4, GPT-3.5-turbo, and more
- **Streaming** — Stream responses token-by-token
- **Embeddings** — Generate vector embeddings for text
- **Images** — Generate and edit images with DALL·E
- **Audio** — Transcribe and translate audio with Whisper
- **Files** — Upload and manage files
- **Fine-tuning** — Fine-tune models on custom data
- **Assistants** — Build AI assistants with threads and runs
- **Structured Outputs** — Parse responses into typed Go structs

## Streaming

```go
stream := client.Chat.Completions.NewStreaming(
	context.Background(),
	openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Count to 10"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	},
)

for stream.Next() {
	chunk := stream.Current()
	if len(chunk.Choices) > 0 {
		fmt.Print(chunk.Choices[0].Delta.Content)
	}
}

if err := stream.Err(); err != nil {
	log.Fatal(err)
}
```

## Configuration

```go
client := openai.NewClient(
	option.WithAPIKey("your-api-key"),
	option.WithBaseURL("https://api.openai.com/v1"),
	option.WithOrganization("your-org-id"),
	option.WithMaxRetries(3),
)
```

## Environment Variables

| Variable | Description |
|---|---|
| `OPENAI_API_KEY` | Your OpenAI API key |
| `OPENAI_ORG_ID` | Your OpenAI organization ID (optional) |
| `OPENAI_BASE_URL` | Custom API base URL (optional) |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and contribution guidelines.

## License

This project is licensed under the Apache License 2.0 — see the [LICENSE](LICENSE) file for details.
