package agents

import (
	"context"
	"fmt"
	"sync"

	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/mu"

	"github.com/openai/openai-go/v2"
)

var (
	cookAgentInstance mu.Agent
	cookAgentOnce     sync.Once
)

// GetCookAgent returns the singleton instance of the sorcerer agent
func GetCookAgent(ctx context.Context, client openai.Client) mu.Agent {
	cookAgentOnce.Do(func() {
		cookAgentInstance = createCookAgent(ctx, client)
	})
	return cookAgentInstance
}

// Huey, Dewey, and Louie
func createCookAgent(ctx context.Context, client openai.Client) mu.Agent {

	name := helpers.GetEnvOrDefault("COOK_NAME", "Dewey")
	model := helpers.GetEnvOrDefault("COOK_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("COOK_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("COOK_MODEL_TOP_P", "0.9"))

	// [RAG]  Initialize the vector store for the agent
	errEmbedding := GenerateEmbeddings(ctx, &client, name, helpers.GetEnvOrDefault("KNOWLEDGE_BASE_PATH", ""))
	if errEmbedding != nil {
		fmt.Println("🔶 Error generating embeddings for cook agent:", errEmbedding)
	}

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	var systemInstructions openai.ChatCompletionMessageParamUnion

	systemInstructionsContentPath := helpers.GetEnvOrDefault("SYSTEM_INSTRUCTIONS_PATH", "")
	if systemInstructionsContentPath == "" {
		fmt.Println("🔶 No SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
	}
	// Read the content of the file at systemInstructionsContentPath
	systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		systemInstructions = openai.SystemMessage("You are a useful assistant.")
	} else {
		systemInstructions = openai.SystemMessage(systemInstructionsContent)
	}

	chatAgent, err := mu.NewAgent(ctx, name,
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       model,
			Temperature: openai.Opt(temperature),
			TopP:        openai.Opt(topP),
			Messages: []openai.ChatCompletionMessageParamUnion{
				systemInstructions,
			},
		}),
	)
	// IMPORTANT: Fake agent if error
	if err != nil {
		fmt.Println("🔶 Error creating cook agent, creating ghost agent instead:", err)
		return NewGhostAgent("[Ghost] " + name)
	}

	return chatAgent

}
