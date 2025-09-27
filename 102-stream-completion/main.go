package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	// Docker Model Runner Chat base URL
	baseURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("COOK_MODEL")

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(os.Getenv("SYSTEM_INSTRUCTION")),
		openai.UserMessage(os.Getenv("USER_PROMPT")),
	}

	// IMPORTANT: Adjust temperature and top_p for desired creativity and coherence
	temperature, _ := strconv.ParseFloat(os.Getenv("TEMPERATURE"), 64)
	topP , _ := strconv.ParseFloat(os.Getenv("TOP_P"), 64)

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(temperature),
		TopP: 	 openai.Opt(topP),
	}
	// NOTE:: Starting a streaming chat completion
	stream := client.Chat.Completions.NewStreaming(ctx, param)

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		log.Fatalln("😡:", err)
	}
}
