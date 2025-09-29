package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	// IMPORTANT: Adjust temperature and top_p for desired creativity and coherence
	temperature, _ := strconv.ParseFloat(os.Getenv("TEMPERATURE"), 64)
	topP, _ := strconv.ParseFloat(os.Getenv("TOP_P"), 64)
	agentName := os.Getenv("AGENT_NAME")
	systemInstructions := os.Getenv("SYSTEM_INSTRUCTIONS")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s](%s) ask me something - /bye to exit> ", agentName, model)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemInstructions),
			openai.UserMessage(userMessage),
		}

		param := openai.ChatCompletionNewParams{
			Messages:    messages,
			Model:       model,
			Temperature: openai.Opt(temperature),
			TopP:        openai.Opt(topP),
		}

		stream := client.Chat.Completions.NewStreaming(ctx, param)

		fmt.Println()

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
		fmt.Println("\n\n", strings.Repeat("-", 80))

	}

}
