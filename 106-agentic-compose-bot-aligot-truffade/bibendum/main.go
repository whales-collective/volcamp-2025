package main

import (
	"bufio"
	"context"
	"encoding/json"
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

	// ✋ NOTE: load the system instructions from a file
	systemInstructions, err := os.ReadFile(os.Getenv("SYSTEM_INSTRUCTIONS_PATH"))
	if err != nil {
		log.Fatal("😡:", err)
	}
	aligotKnowledgeBase, err := os.ReadFile(os.Getenv("ALIGOT_KNOWLEDGE_BASE_PATH"))
	if err != nil {
		log.Fatal("😡:", err)
	}
	truffadeKnowledgeBase, err := os.ReadFile(os.Getenv("TRUFFADE_KNOWLEDGE_BASE_PATH"))
	if err != nil {
		log.Fatal("😡:", err)
	}

	// NOTE: initialize the messages slice with a system message to set the behavior of the assistant
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(string(systemInstructions)),
		openai.SystemMessage("BASE DE CONNAISSANCE:\n"+string(aligotKnowledgeBase)),
		openai.SystemMessage("BASE DE CONNAISSANCE:\n"+string(truffadeKnowledgeBase)),
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖🧠 [%s](%s) ask me something - /bye to exit> ", agentName, model)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		if strings.HasPrefix(userMessage, "/memory") {
			DisplayConversationalMemory(messages)
			continue
		}

		// NOTE: append the user message to the messages slice
		messages = append(messages, openai.UserMessage(userMessage))

		param := openai.ChatCompletionNewParams{
			Messages:    messages,
			Model:       model,
			Temperature: openai.Opt(temperature),
			TopP:        openai.Opt(topP),	
		}

		stream := client.Chat.Completions.NewStreaming(ctx, param)

		fmt.Println()

		// IMPORTANT:
		answer := ""
		for stream.Next() {
			chunk := stream.Current()
			// Stream each chunk as it arrives
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				content := chunk.Choices[0].Delta.Content
				// NOTE: accumulate the content of the assistant's response
				answer += content
				fmt.Print(content)
			}
		}

		if err := stream.Err(); err != nil {
			log.Fatalln("😡 stream completion error:", err)
		}

		// NOTE: Append the assistant's response to the messages slice
		messages = append(messages, openai.AssistantMessage(answer))

		fmt.Println("\n\n", strings.Repeat("-", 80))

	}

}

// MessageToMap converts an OpenAI chat message to a map with string keys and values
func MessageToMap(message openai.ChatCompletionMessageParamUnion) (map[string]string, error) {
	jsonData, err := message.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	stringMap := make(map[string]string)
	for key, value := range result {
		if str, ok := value.(string); ok {
			stringMap[key] = str
		}
	}

	return stringMap, nil
}

func DisplayConversationalMemory(messages []openai.ChatCompletionMessageParamUnion) {
	// remove the /debug part from the input
	fmt.Println()
	fmt.Println("📝 Messages history / Conversational memory:")
	for i, message := range messages {
		printableMessage, err := MessageToMap(message)
		if err != nil {
			fmt.Printf("😡 Error converting message to map: %v\n", err)
			continue
		}
		fmt.Print("-", i, " ")
		fmt.Print(printableMessage["role"], ": ")
		fmt.Println(printableMessage["content"])
	}
	fmt.Println("📝 End of the messages")
	fmt.Println()
}
