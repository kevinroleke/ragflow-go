package ragflow

import (
	"encoding/json"
	"strconv"
	"time"
)

// UnixTime handles Unix timestamp unmarshaling from JSON
type UnixTime struct {
	time.Time
}

func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		return nil
	}

	// Try to parse as Unix timestamp (number)
	if data[0] != '"' {
		timestamp, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		ut.Time = time.Unix(timestamp, 0)
		return nil
	}

	// Try to parse as RFC3339 string
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}
	ut.Time = parsedTime
	return nil
}

func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ut.Time.Format(time.RFC3339))
}

type ChatCompletionRequest struct {
	Model            string                 `json:"model"`
	Messages         []ChatMessage          `json:"messages"`
	Stream           bool                   `json:"stream"`
	ConversationID   string                 `json:"conversation_id,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int64                        `json:"created"`
	Model             string                       `json:"model"`
	SystemFingerprint string                       `json:"system_fingerprint"`
	Choices           []ChatCompletionChoice       `json:"choices"`
	Usage             ChatCompletionUsage          `json:"usage"`
	Reference         ChatCompletionReference      `json:"reference"`
}

type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	Delta        ChatMessage `json:"delta,omitempty"`
	FinishReason string      `json:"finish_reason"`
}

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionReference struct {
	Chunks []ReferenceChunk `json:"chunks"`
}

type ReferenceChunk struct {
	ChunkID      string                 `json:"chunk_id"`
	ContentLTKS  string                 `json:"content_ltks"`
	ContentWTKS  string                 `json:"content_with_weight"`
	DocumentID   string                 `json:"document_id"`
	DocumentName string                 `json:"document_name"`
	Dataset      []string               `json:"dataset"`
	Similarity   float64                `json:"similarity"`
	Vector       []float64              `json:"vector"`
	Positions    [][]int                `json:"positions"`
	Image        string                 `json:"img_id"`
	Term         map[string]interface{} `json:"term_similarity"`
}

type Dataset struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Language            string                 `json:"language"`
	Permission          string                 `json:"permission"`
	DocumentCount       int                    `json:"document_count"`
	ChunkCount          int                    `json:"chunk_count"`
	ParseMethod         string                 `json:"parse_method"`
	ParserConfig        map[string]interface{} `json:"parser_config"`
	CreateTime          UnixTime               `json:"create_time"`
	UpdateTime          UnixTime               `json:"update_time"`
	CreatedBy           string                 `json:"created_by"`
	Avatar              string                 `json:"avatar"`
	EmbeddingModel      string                 `json:"embedding_model"`
	TenantID            string                 `json:"tenant_id"`
	VectorSimilarity    float64                `json:"vector_similarity_weight"`
	Parser              string                 `json:"parser"`
	ChunkTokenCount     int                    `json:"chunk_token_count"`
	ChunkTokenNumber    int                    `json:"chunk_token_number"`
}

type CreateDatasetRequest struct {
	Name                string                 `json:"name"`
	Description         string                 `json:"description,omitempty"`
	Language            string                 `json:"language,omitempty"`
	Permission          string                 `json:"permission,omitempty"`
	ParseMethod         string                 `json:"parse_method,omitempty"`
	ParserConfig        map[string]interface{} `json:"parser_config,omitempty"`
	Avatar              string                 `json:"avatar,omitempty"`
	EmbeddingModel      string                 `json:"embedding_model,omitempty"`
	VectorSimilarity    float64                `json:"vector_similarity_weight,omitempty"`
	Parser              string                 `json:"parser,omitempty"`
	ChunkTokenNumber    int                    `json:"chunk_token_number,omitempty"`
}

type UpdateDatasetRequest struct {
	Name                string                 `json:"name,omitempty"`
	Description         string                 `json:"description,omitempty"`
	Language            string                 `json:"language,omitempty"`
	Permission          string                 `json:"permission,omitempty"`
	ParseMethod         string                 `json:"parse_method,omitempty"`
	ParserConfig        map[string]interface{} `json:"parser_config,omitempty"`
	Avatar              string                 `json:"avatar,omitempty"`
	EmbeddingModel      string                 `json:"embedding_model,omitempty"`
	VectorSimilarity    float64                `json:"vector_similarity_weight,omitempty"`
	Parser              string                 `json:"parser,omitempty"`
	ChunkTokenNumber    int                    `json:"chunk_token_number,omitempty"`
}

type Document struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Size        int64                  `json:"size"`
	Token       int                    `json:"token"`
	ChunkNumber int                    `json:"chunk_number"`
	Progress    float64                `json:"progress"`
	ProgressMsg string                 `json:"progress_msg"`
	Process     string                 `json:"process"`
	Source      string                 `json:"source"`
	CreateTime  UnixTime               `json:"create_time"`
	UpdateTime  UnixTime               `json:"update_time"`
	CreatedBy   string                 `json:"created_by"`
	Run         string                 `json:"run"`
	Parser      map[string]interface{} `json:"parser"`
	Location    string                 `json:"location"`
}

type Chunk struct {
	ID           string                 `json:"id"`
	Content      string                 `json:"content"`
	DocumentID   string                 `json:"document_id"`
	DocumentName string                 `json:"document_name"`
	DatasetIDs   []string               `json:"dataset_ids"`
	Important    bool                   `json:"important"`
	CreateTime   UnixTime               `json:"create_time"`
	UpdateTime   UnixTime               `json:"update_time"`
	Positions    [][]int                `json:"positions"`
	Available    bool                   `json:"available"`
	TermWeights  map[string]interface{} `json:"term_weights"`
}

type UpdateChunkRequest struct {
	Content   string `json:"content,omitempty"`
	Important bool   `json:"important,omitempty"`
}

type Variable struct {
	Key      string `json:"key"`
	Optional bool   `json:"optional"`
}

type Prompt struct {
	EmptyResponse             string     `json:"empty_response"`
	KeywordsSimilarityWeight  float64    `json:"keywords_similarity_weight"`
	Opener                    string     `json:"opener"`
	Prompt                    string     `json:"prompt"`
	RefineMultiturn           bool       `json:"refine_multiturn"`
	RerankModel               string     `json:"rerank_model"`
	ShowQuote                 bool       `json:"show_quote"`
	SimilarityThreshold       float64    `json:"similarity_threshold"`
	TopN                      int        `json:"top_n"`
	TTS                       bool       `json:"tts"`
	Variables                 []Variable `json:"variables"`
}

type LLMModelSettings struct {
	ModelName string `json:"model_name"`
}

type Assistant struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Avatar          string                 `json:"avatar"`
	Language        string                 `json:"language"`
	Prompt          Prompt 				   `json:"prompt"`
	LLMSetting      map[string]interface{} `json:"llm_setting"`
	LLMModel        string                 `json:"llm_model"`
	DatasetIDs      []string               `json:"dataset_ids"`
	TopK            int                    `json:"top_k"`
	SimilarityThreshold float64            `json:"similarity_threshold"`
	VectorSimilarityWeight float64         `json:"vector_similarity_weight"`
	TopP            float64                `json:"top_p"`
	Temperature     float64                `json:"temperature"`
	MaxTokens       int                    `json:"max_tokens"`
	PresencePenalty float64                `json:"presence_penalty"`
	FrequencyPenalty float64               `json:"frequency_penalty"`
	CreateTime      UnixTime               `json:"create_time"`
	UpdateTime      UnixTime               `json:"update_time"`
	CreatedBy       string                 `json:"created_by"`
	TenantID        string                 `json:"tenant_id"`
	ReRank          bool                   `json:"rerank"`
	EmptyResponse   string                 `json:"empty_response"`
	MaxReference    int                    `json:"max_reference"`
	ReRankModel     string                 `json:"rerank_model"`
}

type CreateAssistantRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description,omitempty"`
	DatasetIDs      []string               `json:"dataset_ids,omitempty"`
	Avatar          string                 `json:"avatar,omitempty"`
	Prompt          Prompt `json:"prompt,omitempty"`
	LLMModel		LLMModelSettings `json:"llm"`
}

type UpdateAssistantRequest struct {
	Name            string                 `json:"name,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Avatar          string                 `json:"avatar,omitempty"`
	Language        string                 `json:"language,omitempty"`
	Prompt          string                 `json:"prompt,omitempty"`
	LLMSetting      map[string]interface{} `json:"llm_setting,omitempty"`
	LLMModel        string                 `json:"llm_model,omitempty"`
	DatasetIDs      []string               `json:"dataset_ids,omitempty"`
	TopK            int                    `json:"top_k,omitempty"`
	SimilarityThreshold float64            `json:"similarity_threshold,omitempty"`
	VectorSimilarityWeight float64         `json:"vector_similarity_weight,omitempty"`
	TopP            float64                `json:"top_p,omitempty"`
	Temperature     float64                `json:"temperature,omitempty"`
	MaxTokens       int                    `json:"max_tokens,omitempty"`
	PresencePenalty float64                `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64               `json:"frequency_penalty,omitempty"`
	ReRank          bool                   `json:"rerank,omitempty"`
	EmptyResponse   string                 `json:"empty_response,omitempty"`
	MaxReference    int                    `json:"max_reference,omitempty"`
	ReRankModel     string                 `json:"rerank_model,omitempty"`
}

type Session struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Messages   []ChatMessage `json:"messages"`
	CreateTime UnixTime `json:"create_time"`
	UpdateTime UnixTime `json:"update_time"`
	CreatedBy  string    `json:"created_by"`
}

type CreateSessionRequest struct {
	Name string `json:"name"`
}

type UpdateSessionRequest struct {
	Name string `json:"name"`
}

type Agent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Avatar      string                 `json:"avatar"`
	Language    string                 `json:"language"`
	DSL         map[string]interface{} `json:"dsl"`
	CreateTime  UnixTime               `json:"create_time"`
	UpdateTime  UnixTime               `json:"update_time"`
	CreatedBy   string                 `json:"created_by"`
	TenantID    string                 `json:"tenant_id"`
}

type CreateAgentRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Avatar      string                 `json:"avatar,omitempty"`
	Language    string                 `json:"language,omitempty"`
	DSL         map[string]interface{} `json:"dsl,omitempty"`
}

type UpdateAgentRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Avatar      string                 `json:"avatar,omitempty"`
	Language    string                 `json:"language,omitempty"`
	DSL         map[string]interface{} `json:"dsl,omitempty"`
}

type ListResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total int `json:"total"`
		Items []T `json:"items"`
	} `json:"data"`
}

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ArrayResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
}
