package azureaisearch

import (
	"context"
	// "encoding/json"
	"errors"
	// "fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// Store is a wrapper to use azure AI search rest API.
type Store struct {
	azureAISearchEndpoint string
	azureAISearchAPIKey   string
	embedder              embeddings.Embedder
	client                *http.Client
	vectorField           string
}

var (
	// ErrNumberOfVectorDoesNotMatch when providing documents,
	// the number of vectors generated should be equal to the number of docs.
	ErrNumberOfVectorDoesNotMatch = errors.New(
		"number of vectors from embedder does not match number of documents",
	)
	// ErrAssertingMetadata SearchScore is stored as float64.
	ErrAssertingSearchScore = errors.New(
		"couldn't assert @search.score to float64",
	)
	// ErrAssertingMetadata Metadata is stored as string.
	ErrAssertingMetadata = errors.New(
		"couldn't assert metadata to string",
	)
	// ErrAssertingContent Content is stored as string.
	ErrAssertingContent = errors.New(
		"couldn't assert content to string",
	)
)

// New creates a vectorstore for azure AI search
// and returns the `Store` object needed by the other accessors.
func New(opts ...Option) (Store, error) {
	s := Store{
		client:      http.DefaultClient,
		vectorField: "contentVector", // default vector field
	}

	if err := applyClientOptions(&s, opts...); err != nil {
		return s, err
	}

	return s, nil
}

var _ vectorstores.VectorStore = &Store{}

// AddDocuments adds the text and metadata from the documents to the Chroma collection associated with 'Store'.
// and returns the ids of the added documents.
func (s *Store) AddDocuments(
	ctx context.Context,
	docs []schema.Document,
	options ...vectorstores.Option,
) ([]string, error) {
	opts := s.getOptions(options...)
	ids := []string{}

	texts := []string{}

	for _, doc := range docs {
		texts = append(texts, doc.PageContent)
	}

	vectors, err := s.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return ids, err
	}

	if len(vectors) != len(docs) {
		return ids, ErrNumberOfVectorDoesNotMatch
	}
	for i, doc := range docs {
		id := uuid.NewString()
		if err = s.UploadDocument(ctx, id, opts.NameSpace, doc.PageContent, vectors[i], doc.Metadata); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// SimilaritySearch creates a vector embedding from the query using the embedder
// and queries to find the most similar documents.
func (s *Store) SimilaritySearch(
	ctx context.Context,
	query string,
	numDocuments int,
	options ...vectorstores.Option,
) ([]schema.Document, error) {
	opts := s.getOptions(options...)

	queryVector, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	payload := SearchDocumentsRequestInput{
		VectorQueries: []SearchDocumentsRequestInputVector{{
			Kind:   "vector",
			Fields: s.vectorField,
			Vector: queryVector,
			K:      numDocuments,
		}},
	}

	if filter, ok := opts.Filters.(string); ok {
		payload.Filter = filter
	}

	searchResults := SearchDocumentsRequestOuput{}
	if err := s.SearchDocuments(ctx, opts.NameSpace, payload, &searchResults); err != nil {
		return nil, err
	}

	documents := make([]schema.Document, 0, len(searchResults.Value))
	for _, result := range searchResults.Value {
		doc, err := assertResultValues(result)
		if err != nil {
			return nil, err
		}
		documents = append(documents, *doc)
	}

	return documents, nil
}

func assertResultValues(searchResult map[string]interface{}) (*schema.Document, error) {
	var score float32
	if scoreFloat64, ok := searchResult["@search.score"].(float64); ok {
		score = float32(scoreFloat64)
	} else {
		return nil, ErrAssertingSearchScore
	}

	metadata := map[string]interface{}{}
	for key, value := range searchResult {
		if key == "@search.score" || key == "chunk" || key == "text_vector" {
			continue
		}
		if strValue, ok := value.(string); ok {
			metadata[key] = strValue
		} else {
			return nil, ErrAssertingMetadata
		}
	}

	var pageContent string
	var ok bool
	if pageContent, ok = searchResult["chunk"].(string); !ok {
		return nil, ErrAssertingContent
	}

	return &schema.Document{
		PageContent: pageContent,
		Metadata:    metadata,
		Score:       score,
	}, nil
}
