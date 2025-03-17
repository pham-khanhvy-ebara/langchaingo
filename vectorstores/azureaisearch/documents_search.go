package azureaisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// QueryType pseudo enum for SearchDocumentsRequestInput queryType property.
type QueryType string

const (
	QueryTypeSimple   QueryType = "simple"
	QueryTypeFull     QueryType = "full"
	QueryTypeSemantic QueryType = "semantic"
)

// QueryCaptions pseudo enum for SearchDocumentsRequestInput queryCaptions property.
type QueryCaptions string

const (
	QueryTypeExtractive QueryCaptions = "extractive"
	QueryTypeNone       QueryCaptions = "none"
)

// SpellerType pseudo enum for SearchDocumentsRequestInput spellerType property.
type SpellerType string

const (
	SpellerTypeLexicon SpellerType = "lexicon"
	SpellerTypeNone    SpellerType = "none"
)

// HybridSearchParams represents hybrid search parameters
type HybridSearchParams struct {
	MaxTextRecallSize int    `json:"maxTextRecallSize,omitempty"`
	CountAndFacetMode string `json:"countAndFacetMode,omitempty"`
}

// SearchDocumentsRequestInput is the input struct to format a payload in order to search for a document.
type SearchDocumentsRequestInput struct {
	Count                         bool                                `json:"count,omitempty"`
	Captions                      QueryCaptions                       `json:"captions,omitempty"`
	Facets                        []string                            `json:"facets,omitempty"`
	Filter                        string                              `json:"filter,omitempty"`
	Highlight                     string                              `json:"highlight,omitempty"`
	HighlightPostTag              string                              `json:"highlightPostTag,omitempty"`
	HighlightPreTag               string                              `json:"highlightPreTag,omitempty"`
	MinimumCoverage               int16                               `json:"minimumCoverage,omitempty"`
	Orderby                       string                              `json:"orderby,omitempty"`
	QueryType                     QueryType                           `json:"queryType,omitempty"`
	QueryLanguage                 string                              `json:"queryLanguage,omitempty"`
	Speller                       SpellerType                         `json:"speller,omitempty"`
	SemanticConfiguration         string                              `json:"semanticConfiguration,omitempty"`
	SemanticErrorHandling         string                              `json:"semanticErrorHandling,omitempty"`
	SemanticMaxWaitInMilliseconds int16                               `json:"semanticMaxWaitInMilliseconds,omitempty"`
	SemanticQuery                 string                              `json:"semanticQuery,omitempty"`
	SemanticFields                string                              `json:"semanticFields,omitempty"`
	Answers                       string                              `json:"answers,omitempty"`
	QueryRewrites                 string                              `json:"queryRewrites,omitempty"`
	ScoringParameters             []string                            `json:"scoringParameters,omitempty"`
	ScoringProfile                string                              `json:"scoringProfile,omitempty"`
	Search                        string                              `json:"search,omitempty"`
	SearchFields                  string                              `json:"searchFields,omitempty"`
	SearchMode                    string                              `json:"searchMode,omitempty"`
	SessionID                     string                              `json:"sessionId,omitempty"`
	ScoringStatistics             string                              `json:"scoringStatistics,omitempty"`
	Select                        string                              `json:"select,omitempty"`
	Skip                          int                                 `json:"skip,omitempty"`
	Top                           int                                 `json:"top,omitempty"`
	VectorQueries                 []SearchDocumentsRequestInputVector `json:"vectorQueries,omitempty"`
	VectorFilterMode              string                              `json:"vectorFilterMode,omitempty"`
	HybridSearch                  HybridSearchParams                  `json:"hybridSearch,omitempty"`
}

// VectorQueryThreshold represents a threshold in vector queries
type ThresholdParams struct {
	Kind string `json:"kind,omitempty"`
}

// SearchDocumentsRequestInputVector is the input struct for vector search.
type SearchDocumentsRequestInputVector struct {
	Kind           string          `json:"kind,omitempty"`
	Vector         []float32       `json:"vector,omitempty"`
	Fields         string          `json:"fields,omitempty"`
	K              int             `json:"k,omitempty"`
	Exhaustive     bool            `json:"exhaustive,omitempty"`
	Oversampling   int             `json:"oversampling,omitempty"`
	Weight         float32         `json:"weight,omitempty"`
	Threshold      ThresholdParams `json:"threshold,omitempty"`
	FilterOverride string          `json:"filterOverride,omitempty"`
}

// FacetResult represents a single facet result in search facets
type FacetResult struct {
	Count           int                    `json:"count,omitempty"`
	Sum             int                    `json:"sum,omitempty"`
	SearchFacets    string                 `json:"@search.facets,omitempty"`
	AdditionalProp1 map[string]interface{} `json:"additionalProp1,omitempty"`
}

// SearchAnswer represents an answer in search answers
type SearchAnswer struct {
	Score           float64                `json:"score,omitempty"`
	Key             string                 `json:"key,omitempty"`
	Text            string                 `json:"text,omitempty"`
	Highlights      string                 `json:"highlights,omitempty"`
	AdditionalProp1 map[string]interface{} `json:"additionalProp1,omitempty"`
}

// QueryRewrite represents a query rewrite in search debug
type QueryRewrite struct {
	InputQuery string   `json:"inputQuery,omitempty"`
	Rewrites   []string `json:"rewrites,omitempty"`
}

// SearchDebug represents debug information in search results
type SearchDebug struct {
	QueryRewrites struct {
		Text    QueryRewrite   `json:"text,omitempty"`
		Vectors []QueryRewrite `json:"vectors,omitempty"`
	} `json:"queryRewrites,omitempty"`
}

// SearchCaption represents a caption in search results
type SearchCaption struct {
	Text            string                 `json:"text,omitempty"`
	Highlights      string                 `json:"highlights,omitempty"`
	AdditionalProp1 map[string]interface{} `json:"additionalProp1,omitempty"`
}

// SemanticFieldState represents the state of a semantic field
type SemanticFieldState struct {
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
}

// SemanticRerankerInput represents reranker input in semantic debug info
type SemanticRerankerInput struct {
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Keywords string `json:"keywords,omitempty"`
}

// VectorScore represents a vector score in document debug info
type VectorScore struct {
	SearchScore      float64 `json:"searchScore,omitempty"`
	VectorSimilarity float64 `json:"vectorSimilarity,omitempty"`
}

// DocumentDebugInfo represents debug information for a document
type DocumentDebugInfo struct {
	Semantic struct {
		TitleField    SemanticFieldState    `json:"titleField,omitempty"`
		ContentFields []SemanticFieldState  `json:"contentFields,omitempty"`
		KeywordFields []SemanticFieldState  `json:"keywordFields,omitempty"`
		RerankerInput SemanticRerankerInput `json:"rerankerInput,omitempty"`
	} `json:"semantic,omitempty"`
	Vectors struct {
		Subscores struct {
			Text struct {
				SearchScore float64 `json:"searchScore,omitempty"`
			} `json:"text,omitempty"`
			Vectors       []map[string]VectorScore `json:"vectors,omitempty"`
			DocumentBoost float64                  `json:"documentBoost,omitempty"`
		} `json:"subscores,omitempty"`
	} `json:"vectors,omitempty"`
}

// SearchResultValue represents a single result in the search results
type SearchResultValue struct {
	SearchScore             float64                `json:"@search.score,omitempty"`
	SearchRerankerScore     float64                `json:"@search.rerankerScore,omitempty"`
	SearchHighlights        map[string][]string    `json:"@search.highlights,omitempty"`
	SearchCaptions          []SearchCaption        `json:"@search.captions,omitempty"`
	SearchDocumentDebugInfo []DocumentDebugInfo    `json:"@search.documentDebugInfo,omitempty"`
	AdditionalProp1         map[string]interface{} `json:"additionalProp1,omitempty"`
}

// SearchDocumentsRequestOuput is the output struct for search.
type SearchDocumentsRequestOuput struct {
	OdataCount                            int                         `json:"@odata.count,omitempty"`
	SearchCoverage                        int                         `json:"@search.coverage,omitempty"`
	SearchFacets                          map[string][]FacetResult    `json:"@search.facets,omitempty"`
	SearchAnswers                         []SearchAnswer              `json:"@search.answers,omitempty"`
	SearchDebug                           SearchDebug                 `json:"@search.debug,omitempty"`
	SearchNextPageParameters              SearchDocumentsRequestInput `json:"@search.nextPageParameters,omitempty"`
	Value                                 []map[string]interface{}    `json:"value,omitempty"`
	OdataNextLink                         string                      `json:"@odata.nextLink,omitempty"`
	SearchSemanticPartialResponseReason   string                      `json:"@search.semanticPartialResponseReason,omitempty"`
	SearchSemanticPartialResponseType     string                      `json:"@search.semanticPartialResponseType,omitempty"`
	SearchSemanticQueryRewritesResultType string                      `json:"@search.semanticQueryRewritesResultType,omitempty"`
}

// SearchDocuments send a request to azure AI search Rest API for searching documents.
func (s *Store) SearchDocuments(
	ctx context.Context,
	indexName string,
	payload SearchDocumentsRequestInput,
	output *SearchDocumentsRequestOuput,
) error {
	URL := fmt.Sprintf("%s/indexes('%s')/docs/search.post.search?api-version=2025-03-01-preview", s.azureAISearchEndpoint, indexName)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("err marshalling document for azure ai search: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("err setting request for azure ai search document: %w", err)
	}

	req.Header.Add("content-Type", "application/json")
	if s.azureAISearchAPIKey != "" {
		req.Header.Add("api-key", s.azureAISearchAPIKey)
	}
	return s.httpDefaultSend(req, "search documents on azure ai search", output)
}
