package agents

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
)

type GhostAgent struct {
	name           string
	messages       []openai.ChatCompletionMessageParamUnion
	responseFormat openai.ChatCompletionNewParamsResponseFormatUnion
}

// NewFakeAgent creates a new fake agent instance
func NewGhostAgent(name string) mu.Agent {
	return &GhostAgent{
		name:     name,
		messages: []openai.ChatCompletionMessageParamUnion{},
	}
}



// GetDescription implements mu.Agent.
func (g *GhostAgent) GetDescription() string {
	panic("unimplemented")
}

// GetMetaData implements mu.Agent.
func (g *GhostAgent) GetMetaData() any {
	panic("unimplemented")
}

// SetDescription implements mu.Agent.
func (g *GhostAgent) SetDescription(description string) {
	panic("unimplemented")
}

// SetMetaData implements mu.Agent.
func (g *GhostAgent) SetMetaData(metaData any) {
	panic("unimplemented")
}

// AddMessage implements mu.Agent.
func (g *GhostAgent) AddMessage(message openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// AddMessages implements mu.Agent.
func (g *GhostAgent) AddMessages(messages []openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// GetFirstNMessages implements mu.Agent.
func (g *GhostAgent) GetFirstNMessages(n int) []openai.ChatCompletionMessageParamUnion {
	panic("unimplemented")
}

// GetLastMessage implements mu.Agent.
func (g *GhostAgent) GetLastMessage() (openai.ChatCompletionMessageParamUnion, bool) {
	panic("unimplemented")
}

// GetLastNMessages implements mu.Agent.
func (g *GhostAgent) GetLastNMessages(n int) []openai.ChatCompletionMessageParamUnion {
	panic("unimplemented")
}

// PrependMessage implements mu.Agent.
func (g *GhostAgent) PrependMessage(message openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// PrependMessages implements mu.Agent.
func (g *GhostAgent) PrependMessages(messages []openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// RemoveFirstMessage implements mu.Agent.
func (g *GhostAgent) RemoveFirstMessage() {
	panic("unimplemented")
}

// RemoveLastMessage implements mu.Agent.
func (g *GhostAgent) RemoveLastMessage() {
	panic("unimplemented")
}

// RemoveLastNMessages implements mu.Agent.
func (g *GhostAgent) RemoveLastNMessages(n int) {
	panic("unimplemented")
}

// ResetMessages implements mu.Agent.
func (g *GhostAgent) ResetMessages() {
	panic("unimplemented")
}

// GetName implements mu.Agent.
func (g *GhostAgent) GetName() string {
	return g.name
}

// SetName implements mu.Agent.
func (g *GhostAgent) SetName(name string) {
	g.name = name
}

// DetectToolCalls implements mu.Agent.
func (g *GhostAgent) DetectToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolCallBack func(functionName string, arguments string) (string, error)) (string, []string, string, error) {
	panic("unimplemented")
}

// DetectToolCallsStream implements mu.Agent.
func (g *GhostAgent) DetectToolCallsStream(messages []openai.ChatCompletionMessageParamUnion, toolCallback func(functionName string, arguments string) (string, error), streamCallback func(content string) error) (string, []string, string, error) {
	panic("unimplemented")
}

// GenerateEmbeddingVector implements mu.Agent.
func (g *GhostAgent) GenerateEmbeddingVector(content string) ([]float64, error) {
	panic("unimplemented")
}

// GetMessages implements mu.Agent.
func (g *GhostAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	return g.messages
}

// GetResponseFormat implements mu.Agent.
func (g *GhostAgent) GetResponseFormat() openai.ChatCompletionNewParamsResponseFormatUnion {
	return g.responseFormat
}

// Run implements mu.Agent.
func (g *GhostAgent) Run(Messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	panic("unimplemented")
}

// RunStream simulates streaming completion
func (g *GhostAgent) RunStream(Messages []openai.ChatCompletionMessageParamUnion, callBack func(content string) error) (string, error) {
	g.messages = append(g.messages, Messages...)

	// Extract user message content for simulation
	var userMessage string
	for _, msg := range Messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}

	response := g.simulateResponse(userMessage)

	// Simulate streaming by sending chunks
	words := strings.Fields(response)
	fullResponse := ""

	for _, word := range words {
		chunk := word + " "
		fullResponse += chunk

		// Simulate streaming delay
		time.Sleep(50 * time.Millisecond)

		if err := callBack(chunk); err != nil {
			return fullResponse, err
		}
	}

	return fullResponse, nil
}

// RunStreamWithReasoning implements mu.Agent.
func (g *GhostAgent) RunStreamWithReasoning(Messages []openai.ChatCompletionMessageParamUnion, contentCallback func(content string) error, reasoningCallback func(reasoning string) error) (string, string, error) {
	panic("unimplemented")
}

// RunWithReasoning implements mu.Agent.
func (g *GhostAgent) RunWithReasoning(Messages []openai.ChatCompletionMessageParamUnion) (string, string, error) {
	panic("unimplemented")
}

// SetMessages implements mu.Agent.
func (g *GhostAgent) SetMessages(messages []openai.ChatCompletionMessageParamUnion) {
	g.messages = messages
}

// SetResponseFormat implements mu.Agent.
func (g *GhostAgent) SetResponseFormat(format openai.ChatCompletionNewParamsResponseFormatUnion) {
	g.responseFormat = format
}

// GetModel implements mu.Agent.
func (g *GhostAgent) GetModel() string {
	return "ghost-model"
}

// SetModel implements mu.Agent.
func (g *GhostAgent) SetModel(model string) {
	// No-op for ghost agent
}


// simulateResponse generates a fake AI response based on the input
func (g *GhostAgent) simulateResponse(userMessage string) string {
	responses := map[string]string{
		"hello":     fmt.Sprintf("Hello! I'm %s, your fake AI assistant. How can I help you today?", g.name),
		"weather":   "I'm a fake agent, so I can't check real weather, but let's pretend it's sunny and 72°F!",
		"code":      "Here's some fake code: `func main() { fmt.Println(\"Hello from fake agent!\") }`",
		"time":      "The current time is... well, I'm fake, so let's say it's always coffee time! ☕",
		"calculate": "I calculated that 2+2 = 4 (even fake agents know basic math!)",
		"search":    "I found exactly what you were looking for! (Just kidding, I'm a fake agent)",
	}

	userLower := strings.ToLower(userMessage)
	for keyword, response := range responses {
		if strings.Contains(userLower, keyword) {
			return response
		}
	}

	return fmt.Sprintf("I'm %s, a fake AI agent. You said: \"%s\". I don't have real AI capabilities, but I'm pretending to understand and respond!", g.name, userMessage)
}
