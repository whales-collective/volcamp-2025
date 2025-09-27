package main

import (
	"bibendum/agents"
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/msg"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/tools"
	"github.com/micro-agent/micro-agent-go/agent/ui"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()
	baseURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	mcpHost := helpers.GetEnvOrDefault("MCP_HOST", "http://localhost:9011/mcp")

	mcpClient, err := tools.NewStreamableHttpMCPClient(ctx, mcpHost)
	if err != nil {
		panic(fmt.Errorf("failed to create MCP client: %v", err))
	}

	ui.Println(ui.Purple, "MCP Client initialized successfully")

	// ---------------------------------------------------------
	// TOOLS CATALOG: get the list of tools from the [MCP] client
	// ---------------------------------------------------------
	toolsIndex := mcpClient.OpenAITools()

	DisplayToolsIndex(toolsIndex)

	// for _, tool := range toolsIndex {
	// 	ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	// }

	// ---------------------------------------------------------
	// AGENT: This is the Bibendum agent
	// ---------------------------------------------------------
	bibendumAgent := agents.GetCookAgent(ctx, client, toolsIndex)

	for {
		promptText := "🤖 (/bye to exit) [" + bibendumAgent.GetName() + "]>"
		// PROMPT:
		content, _ := ui.SimplePrompt(promptText, "Type your command here...")

		// USER MESSAGE: content.Input

		// ---------------------------------------------------------
		// Bye [COMMAND]
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/bye") {
			fmt.Println("👋 Goodbye! Thanks for the chat!")
			break
		}

		// DEBUG:
		if strings.HasPrefix(content.Input, "/memory") {
			msg.DisplayHistory(bibendumAgent)
			continue
		}

		// ---------------------------------------------------------
		// AGENT:: + [RAG]
		// ---------------------------------------------------------
		ui.Println(ui.Purple, "<", bibendumAgent.GetName(), "speaking...>")

		// thinkingCtrl := ui.NewThinkingController()
		// thinkingCtrl.Start(ui.Cyan, "Tools detection.....")

		// Create executeFunction with MCP client option
		// Tool execution callback
		executeFn := ExecuteFunction(mcpClient)

		bibendumAgentMessages := []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(content.Input),
		}

		// TOOLS DETECTION:

		// Stream callback for real-time content display
		streamCallback := func(content string) error {
			fmt.Print(content)
			return nil
		}

		// fmt.Println("🚀 Starting streaming tool completion...")
		// fmt.Println(strings.Repeat("=", 50))

		finishReason, toolCallsResults, assistantMessage, err := bibendumAgent.DetectToolCallsStream(bibendumAgentMessages, executeFn, streamCallback)
		if err != nil {
			panic(err)
		}

		// finishReason, toolCallsResults, assistantMessage, err := bibendumAgent.DetectToolCalls(bibendumAgentMessages, executeFn)
		// if err != nil {
		// 	panic(err)
		// }
		// thinkingCtrl.Stop()

		fmt.Printf("Finish Reason: %s\n", finishReason)

		if len(toolCallsResults) > 0 {
			// IMPORTANT: This is the answer from the [MCP] server
			DisplayMCPToolCallResult(toolCallsResults)
		}

		// ASSISTANT MESSAGE:
		// This is the final answer from the agent
		DisplayAgentResponse(assistantMessage)

		fmt.Println()
		fmt.Println()

	}

}

func DisplayToolsIndex(toolsIndex []openai.ChatCompletionToolUnionParam) {
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}
	fmt.Println()
}

func DisplayMCPToolCallResult(results []string) {
	fmt.Println(strings.Repeat("-", 3) + "[MCP RESPONSE]" + strings.Repeat("-", 33))
	fmt.Println(results[0])
	fmt.Println(strings.Repeat("-", 50))
}

func DisplayAgentResponse(assistantMessage string) {
	ui.Println(ui.Green, strings.Repeat("-", 3)+"[AGENT RESPONSE]"+strings.Repeat("-", 31))
	fmt.Println(assistantMessage)
	ui.Println(ui.Green, strings.Repeat("-", 50))
	fmt.Println()
}

func ExecuteFunction(mcpClient *tools.MCPClient) func(string, string) (string, error) {

	return func(functionName string, arguments string) (string, error) {

		fmt.Printf("🟢 %s with arguments: %s\n", functionName, arguments)

		// WAITING: for [CONFIRMATION] function is detected, function execution confirmation
		choice := ui.GetChoice(ui.Yellow, "Do you want to execute this function? (y)es (n)o (a)bort", []string{"y", "n", "a"}, "y")

		switch choice {
		case "n":
			return `{"result": "Function not executed"}`, nil
		case "a": // abort
			return `{"result": "Function not executed"}`,
				&mu.ExitToolCallsLoopError{Message: "Tool execution aborted by user"}

		default: // [YES] if the user confirms (yes)
			ctx := context.Background()
			result, err := mcpClient.CallTool(ctx, functionName, arguments)
			if err != nil {
				return "", fmt.Errorf("MCP tool execution failed: %v", err)
			}
			if len(result.Content) > 0 {
				// Take the first content item and return its text
				resultContent := result.Content[0].(mcp.TextContent).Text
				fmt.Println("✅ Tool executed successfully")
				// ✋ could be JSON or not
				return resultContent, nil

			}
			return `{"result": "Tool executed successfully but returned no content"}`, nil
		}

	}
}
