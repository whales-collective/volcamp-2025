package main

import (
	"context"
	"embeddings-chat/rag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var chunks = []string{
	`# Truffade
	La truffade est un plat traditionnel du Cantal à base de pommes de terre et de tomme fraîche. 
	Cette spécialité rustique est préparée en faisant revenir des pommes de terre en lamelles 
	avec de l'ail et du persil, puis en incorporant la tomme qui fond lentement. 
	Sa texture crémeuse et son goût authentique en font un incontournable de la cuisine auvergnate, 
	particulièrement apprécié lors des soirées d'hiver au coin du feu.`,

	`# Aligot
	L'aligot est une purée de pommes de terre mélangée à de la tomme fraîche et de l'ail. 
	Originaire de l'Aubrac, cette préparation demande un savoir-faire particulier pour obtenir 
	la texture filante caractéristique grâce au brassage énergique de la tomme. 
	Accompagnement traditionnel des saucisses de Toulouse, l'aligot est devenu l'emblème 
	de la gastronomie aveyronnaise et auvergnate, symbole de convivialité et de tradition.`,

	`# Pounti
	Le pounti est un pâté rustique aux herbes typique de la Haute-Loire et du Cantal. 
	Cette terrine salée mélange viande de porc hachée, blettes ou épinards, oeufs et lait. 
	Cuit au four dans un moule, il se déguste chaud ou froid, souvent accompagné de salade. 
	Son goût unique aux herbes sauvages et sa texture moelleuse en font un plat convivial 
	parfait pour les pique-niques et les repas familiaux en Auvergne.`,

	`# Cantal
	Le Cantal est un fromage au lait de vache à pâte pressée non cuite, emblème de l'Auvergne. 
	Fabriqué depuis plus de 2000 ans dans les montagnes du Massif Central, il développe 
	une croûte dorée et une pâte ferme aux arômes complexes selon son affinage. 
	Jeune, entre-deux ou vieux, chaque étape offre des saveurs différentes allant du doux 
	au corsé, faisant du Cantal un trésor gastronomique incontournable de la région.`,
}

func main() {
	ctx := context.Background()

	baseURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	embeddingsModel := os.Getenv("EMBEDDING_MODEL")
	chatModel := os.Getenv("COOK_MODEL")
	systemInstructions := os.Getenv("SYSTEM_INSTRUCTIONS")

	fmt.Println("🌍", baseURL)
	fmt.Println("🧠", embeddingsModel)
	fmt.Println("🤖", chatModel)

	temperature, _ := strconv.ParseFloat(os.Getenv("TEMPERATURE"), 64)
	topP, _ := strconv.ParseFloat(os.Getenv("TOP_P"), 64)

	similaritySearchLimit, _ := strconv.ParseFloat(os.Getenv("SIMILARITY_LIMIT"), 64)
	similaritySearchMaxResults, _ := strconv.Atoi(os.Getenv("SIMILARITY_MAX_RESULTS"))

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	// -------------------------------------------------
	// Create a vector store
	// -------------------------------------------------
	store := rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	// -------------------------------------------------
	// STEP 1: Create and save the embeddings from the chunks
	// -------------------------------------------------
	fmt.Println("⏳ Creating the embeddings...")

	for _, chunk := range chunks {
		// EMBEDDING COMPLETION:
		embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(chunk),
			},
			Model: embeddingsModel,
		})

		if err != nil {
			fmt.Println(err)
		} else {
			_, errSave := store.Save(rag.VectorRecord{
				Prompt:    chunk,
				Embedding: embeddingsResponse.Data[0].Embedding,
			})
			if errSave != nil {
				fmt.Println("😡:", errSave)
			}
		}
	}

	fmt.Println("✋", "Embeddings created, total of records", len(store.Records))
	fmt.Println()

	// -------------------------------------------------
	// Search for similarities
	// -------------------------------------------------

	// USER MESSAGE:
	userQuestion := "Explique moi ce qu'est le Pounti ?"

	fmt.Println("⏳ Searching for similarities...")

	// -------------------------------------------------
	// STEP 2: EMBEDDING COMPLETION:
	// Create embedding from the user question
	// -------------------------------------------------
	embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(userQuestion),
		},
		Model: embeddingsModel,
	})
	if err != nil {
		log.Fatal("😡:", err)
	}

	// -------------------------------------------------
	// STEP 3: SIMILARITY SEARCH: use the vector store to find similar chunks
	// -------------------------------------------------
	// Create a vector record from the user embedding
	embeddingFromUserQuestion := rag.VectorRecord{
		Embedding: embeddingsResponse.Data[0].Embedding,
	}

	similarities, _ := store.SearchTopNSimilarities(embeddingFromUserQuestion, similaritySearchLimit, similaritySearchMaxResults)

	documentsContent := "Documents:\n"

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
		documentsContent += similarity.Prompt
	}
	documentsContent += "\n"
	fmt.Println("\n✋", "Similarities found, total of records", len(similarities))
	fmt.Println()

	// -------------------------------------------------
	// STEP 4: Generate CHAT COMPLETION:
	// -------------------------------------------------
	messages := []openai.ChatCompletionMessageParamUnion{
		// SYSTEM MESSAGE:
		openai.SystemMessage(systemInstructions),
		// SIMILARITIES:
		openai.SystemMessage(documentsContent),
		// USER MESSAGE:
		openai.UserMessage(userQuestion),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       chatModel,
		Temperature: openai.Opt(temperature),
		TopP:        openai.Opt(topP),
	}

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

	fmt.Println()
	store.SaveJSONToFile("vectorstore.json")
	fmt.Println("\n✋", "Vector store saved to vectorstore.json")
}
