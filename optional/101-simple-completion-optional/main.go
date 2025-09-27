package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
		//openai.UserMessage("What is your name?"),
		openai.UserMessage("Quel est ton nom?"),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(0.5),
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println(completion.Choices[0].Message.Content)

}
