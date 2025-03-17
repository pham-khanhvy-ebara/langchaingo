package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/azureaisearch"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create OpenAI client for embeddings
	llm, err := openai.New(
		openai.WithEmbeddingModel(os.Getenv("AZURE_OPENAI_MODEL_NAME")),
		openai.WithAPIType(openai.APITypeAzure),
		openai.WithToken(os.Getenv("AZURE_OPENAI_API_KEY")),
		openai.WithBaseURL(os.Getenv("AZURE_OPENAI_BASE_URL")),
		openai.WithAPIVersion(os.Getenv("AZURE_OPENAI_API_VERSION")),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI client: %v", err)
	}

	// Create embeddings client
	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatalf("Failed to create embeddings client: %v", err)
	}

	// Create a new Azure AI Search client
	searchClient, err := azureaisearch.New(
		azureaisearch.WithEmbedder(e),
		azureaisearch.WithAPIKey(os.Getenv("AZURE_AI_SEARCH_KEY")),
		azureaisearch.WithVectorField("text_vector"),
	)
	if err != nil {
		log.Fatalf("Failed to create search client: %v", err)
	}

	// Create a context
	ctx := context.Background()
	
	// Perform vector search with options
	searchText := "Whistleblower Policy"
	results, err := searchClient.SimilaritySearch(
		ctx,
		searchText,
		3, // Get top 3 results
		vectorstores.WithNameSpace(os.Getenv("AZURE_AI_SEARCH_INDEX")),
		// vectorstores.WithScoreThreshold(scoreThreshold),
		// azureaisearch.WithFilters(filter),
	)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Print results
	fmt.Println("Search results for:", searchText)
	for i, doc := range results {
		fmt.Printf("%d. Content: %s\n   ChunkID: %s\n   ParentID: %s\n   Title: %s\n\n",
			i+1,
			doc.PageContent,
			doc.Metadata["chunk_id"],
			doc.Metadata["parent_id"],
			doc.Metadata["title"],
		)
	}
}
