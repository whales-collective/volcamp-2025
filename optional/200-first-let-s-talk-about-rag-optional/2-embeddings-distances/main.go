package main

import (
	"context"
	"embeddings-demo-next/rag"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var chunks = []string{
	`Les écureuils grimpent dans les arbres`,
	`Les truites nagent dans la rivière`,
	`Les grenouilles nagent dans l'étang`,
	`Les lapins courent dans le champ`,
}

func main() {
	ctx := context.Background()

	baseURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	embeddingsModel := os.Getenv("EMBEDDING_MODEL")

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
	// Create and save the embeddings from the chunks
	// -------------------------------------------------
	fmt.Println("⏳ Creating embeddings from chunks...")

	for _, chunk := range chunks {
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
	userQuestion := "Quels sont les animaux qui nagent ?"
	//userQuestion := "Quels animaux peut on trouver en forêt ?"
	//userQuestion := "Quels animaux peut on trouver dans les champs ?"

	fmt.Println("⏳ Searching for similarities...")

	// -------------------------------------------------
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
	// Create a vector record from the user embedding
	// -------------------------------------------------
	embeddingFromUserQuestion := rag.VectorRecord{
		Embedding: embeddingsResponse.Data[0].Embedding,
	}

	similarities, _ := store.SearchTopNSimilarities(embeddingFromUserQuestion, 0.6, 2)
	// if the limit is to near from 1, the risk is to lose the best match

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
	}
	fmt.Println("✋", "Similarities found, total of records", len(similarities))
	fmt.Println()

}
