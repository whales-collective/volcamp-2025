package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/rag"
	"github.com/micro-agent/micro-agent-go/agent/ui"
	"github.com/openai/openai-go/v2" // imported as openai
	"github.com/openai/openai-go/v2/option"
)

var store rag.MemoryVectorStore
var embeddingsModel string

func main() {
	ctx := context.Background()

	// -------------------------------------------------
	// Create MCP server
	// -------------------------------------------------
	s := server.NewMCPServer(
		"mcp-aligot-server",
		"0.0.0",
	)

	baseURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1/")
	embeddingsModel = helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/granite-embedding-multilingual:latest")

	// -------------------------------------------------
	// Create an OpenAI client
	// -------------------------------------------------
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)
	// -------------------------------------------------
	// Create a vector store (in memory)
	// -------------------------------------------------
	store = rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	// -------------------------------------------------
	// [RAG]  Initialize or loqd the data into the store
	// -------------------------------------------------
	errEmbedding := GenerateEmbeddings(ctx, &client, "aligot_agent", helpers.GetEnvOrDefault("ALIGOT_AGENT_KNOWLEDGE_BASE_PATH", ""))
	if errEmbedding != nil {
		fmt.Println("🔶 Error generating embeddings for aligot agent:", errEmbedding)
	}

	// =================================================
	// TOOLS:
	// =================================================
	aligotTool := mcp.NewTool("search_information_about_aligot",
		mcp.WithDescription(`Search for information about Aligot in the knowledge base.`),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to search for similarities."),
		),
	)
	s.AddTool(aligotTool, aligotHandler(&client))

	aligotToolFrench := mcp.NewTool("trouver_des_informations_sur_aligot",
		mcp.WithDescription(`Rechercher des informations sur l'aligot dans la base de connaissances.`),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to search for similarities."),
		),
	)
	s.AddTool(aligotToolFrench, aligotHandler(&client))

	// -------------------------------------------------
	// Start the HTTP server
	// -------------------------------------------------
	httpPort := helpers.GetEnvOrDefault("MCP_HTTP_PORT", "6060")
	fmt.Println("🌍 MCP HTTP Port:", httpPort)

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)

	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)

	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
}

func aligotHandler(client *openai.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	contentArg, exists := args["content"]
	if !exists || contentArg == nil {
		return nil, fmt.Errorf("missing required parameter 'content'")
	}
	content, ok := contentArg.(string)
	if !ok {
		return nil, fmt.Errorf("parameter 'content' must be a string")
	}

	fmt.Println("🔍 Searching similarities for content:", content)

	threshold := helpers.StringToFloat(helpers.GetEnvOrDefault("LIMIT", "0.6"))
	topN := helpers.StringToInt(helpers.GetEnvOrDefault("MAX_RESULTS", "2"))

	fmt.Println("🔍 Using threshold:", threshold, "and topN:", topN)

	similarities, err := SearchSimilarities(ctx, client, "aligot_agent", content, threshold, topN)
	if err != nil {
		return nil, fmt.Errorf("error searching similarities: %v", err)
	}

	documentsContent := "Similarities found:\n"
	for _, similarity := range similarities {
		documentsContent += fmt.Sprintf("Similarity: %.4f - %s\n", similarity.CosineSimilarity, similarity.Prompt)
	}

	fmt.Println("✅ Found", len(similarities), "similarities")

	return mcp.NewToolResultText(documentsContent), nil
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"status":           "healthy",
		"records":          len(store.Records),
		"embeddings_model": embeddingsModel,
	}
	json.NewEncoder(w).Encode(response)
}

// GenerateEmbeddings reads a context file, splits it into chunks, generates embeddings,
// and stores them in the vector store for the specified agent
func GenerateEmbeddings(ctx context.Context, client *openai.Client, name string, contextInstructionsContentPath string) error {

	// Load the vector store from a file if it exists
	jsonStoreFilePath := helpers.GetEnvOrDefault("VECTOR_STORES_PATH", "./data") + "/" + strings.ToLower(name) + "_vector_store.json"
	fmt.Println("🔶 Loading vector store from:", jsonStoreFilePath)

	// ---------------------------------------------------------
	// [VECTOR STORE] Loading or creating the vector store
	// ---------------------------------------------------------
	err := store.Load(jsonStoreFilePath)
	if err != nil {
		// ---------------------------------------------------------
		// BEGIN: If the file does not exist, create a new vector store
		// ---------------------------------------------------------
		if os.IsNotExist(err) {
			fmt.Println("🔶 No existing vector store found, starting fresh:", err)

			ui.Println(ui.Green, strings.Repeat("─", 80))
			ui.Println(ui.Green, "🚧 Generating embeddings for agent:", name)
			ui.Println(ui.Green, strings.Repeat("─", 80))

			// EMBEDDING AGENT: Create an embedding agent to generate embeddings
			embeddingAgent, err := mu.NewAgent(ctx, "vector-agent",
				mu.WithClient(*client),
				mu.WithEmbeddingParams(
					openai.EmbeddingNewParams{
						Model: embeddingsModel,
					},
				),
			)
			if err != nil {
				fmt.Println("🔶 Error creating embedding agent", err)
				return err
			}

			fmt.Println("✅ Embedding agent created successfully")

			if contextInstructionsContentPath == "" {
				fmt.Println("🔶 No context path provided, using default instructions.")
				return fmt.Errorf("no context path provided")
			}

			// Read the content of the file at contextInstructionsContentPath
			contextInstructionsContent, err := helpers.ReadTextFile(contextInstructionsContentPath)
			if err != nil {
				fmt.Println("🔶 Error reading the file, using default instructions:", err)
				return err
			}

			// CHUNKS: Split the content into chunks for embedding
			chunks := rag.SplitMarkdownBySections(contextInstructionsContent)

			for idx, chunk := range chunks {
				fmt.Println("🔶 Chunk", idx, ":", chunk)
				embeddingVector, err := embeddingAgent.GenerateEmbeddingVector(chunk)
				if err != nil {
					return err
				}
				_, errSave := store.Save(rag.VectorRecord{
					Prompt:    chunk,
					Embedding: embeddingVector,
				})

				if errSave != nil {
					fmt.Println("🔴 When saving the vector", errSave)
					return errSave
				}
				fmt.Println("✅ Chunk", idx, "saved with embedding:", len(embeddingVector))
			}
			fmt.Println("📝 Total records in the vector store:", len(store.Records))

			// [RAG] Save the vector store to a file
			err = store.Persist(jsonStoreFilePath)
			if err != nil {
				fmt.Println("🔶 Error saving vector store:", err)
				return err
			}
			fmt.Println("✅ Vector store saved to", jsonStoreFilePath)
			fmt.Println("💾 Vector store initialized with", len(store.Records), "records.")

			ui.Println(ui.Green, strings.Repeat("─", 80))
			fmt.Println()

			return nil
			// ---------------------------------------------------------
			// END: If the file does not exist, create a new vector store
			// ---------------------------------------------------------
		} else {
			fmt.Println("🔶 Error loading vector store:", err)
			return err
		}

	} else {
		fmt.Println("✅ Vector store loaded successfully with", len(store.Records), "records")
		return nil // If the store is loaded successfully, no need to regenerate embeddings

	}

}

// SearchSimilarities searches for similar content in the agent's vector store
// based on the input question and returns the top N similar records
func SearchSimilarities(ctx context.Context, client *openai.Client, agentName string, input string, threshold float64, topN int) ([]rag.VectorRecord, error) {

	embeddingAgent, err := mu.NewAgent(ctx, "vector-agent",
		mu.WithClient(*client),
		mu.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest"),
			},
		),
	)
	if err != nil {
		fmt.Println("🔶 Error creating embedding agent", err)
		return nil, err
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("🧠 Creating embedding for input:", input)
	questionEmbeddingVector, err := embeddingAgent.GenerateEmbeddingVector(input)
	if err != nil {
		return nil, err
	}

	questionRecord := rag.VectorRecord{Embedding: questionEmbeddingVector}

	similarities, err := store.SearchTopNSimilarities(questionRecord, threshold, topN)
	if err != nil {
		return nil, err
	}

	fmt.Println("📝 Similarities found:", len(similarities))

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
	}

	fmt.Println(strings.Repeat("-", 80))

	return similarities, nil
}
