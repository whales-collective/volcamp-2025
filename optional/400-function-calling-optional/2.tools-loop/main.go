package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
	"github.com/openai/openai-go/v2/shared/constant"
)

func main() {
	// Step 1: Initialize context for request management
	ctx := context.Background()

	// Step 2: Configure connection to the model runner
	// Read configuration from environment variables
	chatURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_TOOLS")

	// Step 3: Create OpenAI client configured for local model runner
	client := openai.NewClient(
		option.WithBaseURL(chatURL),
		option.WithAPIKey(""), // No API key needed for local deployment
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
	// Step 7: Create user message with multiple requests
	// This will trigger multiple sequential tool calls from the AI
	userQuestion := openai.UserMessage(`
		Je voudrais parler de truffade.
		Peux-tu aussi dire bonjour à Bob pour moi s'il te plaît ?
		Je voudrais aussi parler d'aligot.
		Je voudrais parler de Pizza.
	`)

	// Step 8: Initialize loop control variables
	stopped := false   // Controls the conversation loop
	finishReason := "" // Tracks why AI stopped responding
	//results := []string{}        // Stores tool execution results
	lastAssistantMessage := "" // Final AI message

	// MEMORY:
	// Step 9: Initialize conversation history with user's question
	messages := []openai.ChatCompletionMessageParamUnion{
		userQuestion,
	}

	// Step 10: Configure chat completion parameters
	// ParallelToolCalls set to false for sequential execution
	params := openai.ChatCompletionNewParams{
		ParallelToolCalls: openai.Bool(false), // Execute tools one by one
		Tools:             tools,              // Available tools for AI
		Model:             model,
		Temperature:       openai.Opt(0.0), // Deterministic responses
	}

	fmt.Println(strings.Repeat("=", 80))

	// Step 11: Main conversation loop - continues until AI says "stop"
	for !stopped {

		// Step 12: Update parameters with current conversation history
		params.Messages = messages

		// [COMPLETION] request
		// Step 13: Send request to AI model and get response
		completion, err := client.Chat.Completions.New(ctx, params)
		if err != nil {
			panic(err)
		}
		// Step 14: Extract finish reason to determine next action
		finishReason = completion.Choices[0].FinishReason

		// Step 15: Handle AI response based on finish reason
		switch finishReason {
		case "tool_calls":
			// Step 16: AI wants to use tools - extract tool calls
			detectedToolCalls := completion.Choices[0].Message.ToolCalls

			if len(detectedToolCalls) > 0 {

				// Step 17: Convert tool calls to proper message format
				// WHY: When AI decides to use tools, it returns toolCalls in its response.
				// We must convert these into ChatCompletionMessageToolCallUnionParam format
				// to add them to conversation history. Without this conversion, AI would
				// lose context of what it requested.
				toolCallParams := make([]openai.ChatCompletionMessageToolCallUnionParam, len(detectedToolCalls))

				for i, toolCall := range detectedToolCalls {
					toolCallParams[i] = openai.ChatCompletionMessageToolCallUnionParam{
						OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
							ID:   toolCall.ID,
							Type: constant.Function("function"),
							Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
								Name:      toolCall.Function.Name,
								Arguments: toolCall.Function.Arguments,
							},
						},
					}
				}

				// Step 18: Create assistant message with tool calls using proper union type
				// WHY: We need to create an "assistant" message containing the tool calls
				// for conversation history. This is like saying: "AI said: 'I want to call
				// these functions with these parameters'". This message will be added to
				// history before executing tools, so AI remembers what it requested.
				assistantMessage := openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						ToolCalls: toolCallParams,
					},
				}

				// Step 19: Add the assistant message with tool calls to the conversation history
				messages = append(messages, assistantMessage)

				// TOOL CALLS:
				// Step 20: Process each detected tool call sequentially
				for _, toolCall := range detectedToolCalls {
					functionName := toolCall.Function.Name
					functionArgs := toolCall.Function.Arguments

					// Step 21: Execute the requested function
					fmt.Printf("▶️ Executing function: %s with args: %s\n", functionName, functionArgs)

					resultContent, err := ExecTool(functionName, functionArgs)

					// Step 22: Handle function execution errors
					if err != nil {
						resultContent = fmt.Sprintf(`{"error": "Function execution failed: %s"}`, err)
					}

					// Step 23: Store result for potential later use
					//results = append(results, resultContent)

					// Step 24: Add tool execution result to conversation history
					// WHY: After executing each tool, we must tell the AI what the result was.
					// This is like a conversation:
					// - AI: "I want to call sayHello with name='Jean-Luc'"
					// - System: "Result: 'Hello Jean-Luc'"
					// AI needs these results to: 1) Know tool executed successfully,
					// 2) Use results for final response, 3) Decide if more tools needed.
					// Without this step, AI would have no idea what happened after requesting
					// tool execution and couldn't generate the requested final report.
					messages = append(
						messages,
						openai.ToolMessage(
							resultContent,
							toolCall.ID,
						),
					)
					fmt.Println("✅ ResultContent", resultContent)
					fmt.Println()
				}

			} else {
				// Step 25: Handle unexpected case with no tool calls
				fmt.Println("😢 No tool calls found in response")
			}

		case "stop":
			// Step 26: AI has finished - no more tools needed
			fmt.Println("🟥 Stopping due to 'stop' finish reason.")
			stopped = true
			lastAssistantMessage = completion.Choices[0].Message.Content

			// Step 27: Add final assistant message to conversation history
			messages = append(messages, openai.AssistantMessage(lastAssistantMessage))
			fmt.Print(strings.Repeat("=", 5), "[Last Assistant Message]", strings.Repeat("=", 51), "\n")
			fmt.Println(lastAssistantMessage)
			fmt.Println(strings.Repeat("=", 80))

		default:
			// Step 28: Handle unexpected finish reasons
			fmt.Printf("🔴 Unexpected response: %s\n", finishReason)
			stopped = true

		}

	}

}

// ExecTool routes function calls to appropriate implementations
// Returns the function result as a string and any execution errors
func ExecTool(functionName, functionArgs string) (string, error) {

	// Route function calls to appropriate handlers
	switch functionName {
	case "parler_de":
		// Parse JSON arguments and execute castSpell function
		args, err := JsonStringToMap(functionArgs)
		return speakAbout(args), err

	case "dire_bonjour":
		// Parse JSON arguments and execute bardicInspiration function
		args, err := JsonStringToMap(functionArgs)
		return sayHello(args), err

	default:
		// Handle unknown function calls
		fmt.Println("Unknown function call:", functionName)
		return "", errors.New("unknown function call")
	}

}

// JsonStringToMap converts a JSON string to a Go map
// Used to parse function arguments from AI tool calls
func JsonStringToMap(jsonString string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// castSpell implements the spell casting functionality
// Extracts the target from arguments and returns a spell casting message
func speakAbout(arguments map[string]interface{}) string {
	// Type assertion to safely extract target parameter
	if topic, ok := arguments["sujet"].(string); ok {

		switch topic {
		case "aligot":
			return "🤖 Il faut parler avec André"
		case "truffade":
			return "🤖 Il faut parler avec Édouard"
		default:
			return "🤔 Pour tout autre sujet il faut parler à Vercingetorix"
		}

	} else {
		return ""
	}
}

// bardicInspiration implements the bardic inspiration functionality
// Extracts the ally from arguments and returns an inspiring message
func sayHello(arguments map[string]interface{}) string {
	// Type assertion to safely extract ally parameter
	if name, ok := arguments["nom"].(string); ok {
		return "👋 Bonjour ! Salutation envoyée à " + name
	} else {
		return ""
	}
}
