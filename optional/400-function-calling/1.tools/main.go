package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
)

func main() {
	ctx := context.Background()

	// Docker Model Runner base URL
	chatURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_TOOLS")

	client := openai.NewClient(
		option.WithBaseURL(chatURL),
		option.WithAPIKey(""),
	)

	// TOOL:
	speakAboutSomethingTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "parler_de",
		Description: openai.String("Parler d'un sujet spécifique"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"sujet": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"sujet"},
		},
	})

	// TOOL:
	sayHelloTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "dire_bonjour",
		Description: openai.String("Dire bonjour à quelqu'un"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"nom": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"nom"},
		},
	})

	// TOOLS: used by the parameters request
	tools := []openai.ChatCompletionToolUnionParam{
		speakAboutSomethingTool,
		sayHelloTool,
	}

	// USER MESSAGE:
	userQuestion := openai.UserMessage(`
		Je voudrais parler de truffade.
		Peux-tu aussi dire bonjour à Bob pour moi s'il te plaît ?
		Je voudrais aussi parler d'aligot.
		Je voudrais parler de Pizza.
	`)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			userQuestion,
		},
		// IMPORTANT: try this:
		//ParallelToolCalls: openai.Bool(false), // default value
		ParallelToolCalls: openai.Bool(true), // Sequential tool calls
		Tools:             tools,
		Model:             model,
		Temperature:       openai.Opt(0.0),
	}

	// Make [COMPLETION] request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	// TOOL CALLS: Extract tool calls from the response
	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Println("😡 No function call")
		fmt.Println()
		return
	}

	fmt.Println(strings.Repeat("=", 80))

	// Display the tool calls
	for _, toolCall := range toolCalls {
		fmt.Println("🛠️", toolCall.Function.Name, toolCall.Function.Arguments)

		// Handle parler_de tool
		if toolCall.Function.Name == "parler_de" {
			// Parse the arguments to extract the "sujet"
			arguments := toolCall.Function.Arguments

			type Parameters struct {
				Sujet string `json:"sujet"`
			}
			var parameters Parameters

			err := json.Unmarshal([]byte(arguments), &parameters)
			if err != nil {
				fmt.Println("😡 Error unmarshaling arguments:", err)
				continue
			}

			switch parameters.Sujet {
			case "aligot":
				fmt.Println("🤖 Il faut parler avec André")
			case "truffade":
				fmt.Println("🤖 Il faut parler avec Édouard")
			default:
				fmt.Println("🤔 Pour tout autre sujet il faut parler à Vercingetorix")
			}
		}

		// Handle dire_bonjour tool
		if toolCall.Function.Name == "dire_bonjour" {
			// Parse the arguments to extract the "nom"
			arguments := toolCall.Function.Arguments
			fmt.Printf("👋 Bonjour ! Salutation envoyée à %s 🤗\n", arguments)
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Println(strings.Repeat("=", 80))
}
