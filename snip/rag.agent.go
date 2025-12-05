package snip

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/firebase/genkit/go/plugins/localvec"
	"github.com/openai/openai-go/option"
)

type RagAgent struct {
	ctx context.Context

	Name    string
	ModelID string

	storeName string
	storePath string

	genKitInstance    *genkit.Genkit
	embedder          ai.Embedder
	docStore          *localvec.DocStore
	documentRetriever ai.Retriever

	embeddingDimension int

	//engineURL string
}

func (agent *RagAgent) GetName() string {
	return agent.Name
}

func NewRagAgent(ctx context.Context, ragAgentConfig RagAgentConfig, storeConfig StoreConfig, opts ...RagAgentOption) (*RagAgent, error) {
	oaiPlugin := &oai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(ragAgentConfig.EngineURL),
		},
	}
	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	if !IsModelAvailable(ctx, ragAgentConfig.EngineURL, ragAgentConfig.ModelID) {
		return nil, fmt.Errorf("model %s is not available at %s", ragAgentConfig.ModelID, ragAgentConfig.EngineURL)
	} else {
		log.Printf("✅ Model %s is available at %s", ragAgentConfig.ModelID, ragAgentConfig.EngineURL)
	}

	// NOTE: Embedder definition/creation
	// you don't need to specify the provider again, it's already set in the plugin 🤔
	// == you don't need to prefix the model name with the provider
	embedder := oaiPlugin.DefineEmbedder(ragAgentConfig.ModelID, nil)

	// get embedder to calculate embedding dimension
	// calculate embedding dimension
	embeddingDimension, err := calculateEmbeddingDimensionForModel(ctx, genKitInstance, embedder)
	if err != nil {
		//log.Printf("Warning: could not calculate embedding dimension: %v", err)
		return nil, fmt.Errorf("error calculating embedding dimension: %w", err)
	}

	if err := localvec.Init(); err != nil {
		return nil, fmt.Errorf("error initializing localvec: %w", err)
	}
	docStore, documentRetriever, err := localvec.DefineRetriever(
		genKitInstance,
		storeConfig.StoreName,
		localvec.Config{
			Embedder: embedder,
			Dir:      storeConfig.StorePath,
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error defining retriever: %w", err)
	}

	ragAgent := &RagAgent{
		ctx:                ctx,
		Name:               ragAgentConfig.Name,
		ModelID:            ragAgentConfig.ModelID,
		storeName:          storeConfig.StoreName,
		storePath:          storeConfig.StorePath,
		genKitInstance:     genKitInstance,
		embedder:           embedder,
		docStore:           docStore,
		documentRetriever:  documentRetriever,
		embeddingDimension: embeddingDimension,
	}

	// Temporary: for debugging
	fmt.Println("🔎 Retriever:", documentRetriever)
	fmt.Println("📚 DocStore:", len(docStore.Data))
	if len(docStore.Data) == 0 {
		fmt.Println("🚧 The document store is empty.")
	} else {
		fmt.Printf("🚧 The document store has %d documents.\n", len(docStore.Data))
	}
	// End Temporary: for debugging

	for _, opt := range opts {
		opt(ragAgent)
	}

	return ragAgent, nil

}

// isStoreInitialized checks if the document store is initialized
func (agent *RagAgent) IsStoreInitialized() bool {
	return agent.docStore != nil
}

// getNumberOfDocuments returns the number of documents in the store
func (agent *RagAgent) GetNumberOfDocuments() int {
	if agent.docStore == nil {
		return 0
	}
	return len(agent.docStore.Data)
}

func (agent *RagAgent) Kind() AgentKind {
	return Rag
}

func (agent *RagAgent) GetInfo() (RagAgentInfo, error) {
	return RagAgentInfo{
		Name:               agent.Name,
		ModelID:            agent.ModelID,
		EmbeddingDimension: agent.embeddingDimension,
		StoreName:          agent.storeName,
		StorePath:          agent.storePath,
		NumberOfDocuments:  agent.GetNumberOfDocuments(),
	}, nil
}

type TextChunk struct {
	Content  string
	Metadata map[string]any
}

func (agent *RagAgent) AddTextChunksToStore(chunks []TextChunk) (int, error) {
	docs := []*ai.Document{}

	for idx, chunk := range chunks {
		// Temporary: for debugging
		fmt.Println("💾", idx, chunk)
		// End Temporary: for debugging
		if chunk.Metadata != nil {
			docs = append(docs, ai.DocumentFromText(chunk.Content, chunk.Metadata))
		} else {
			docs = append(docs, ai.DocumentFromText(chunk.Content, nil))
		}
	}
	fmt.Println("🗂️ Indexing documents...", docs)
	err := localvec.Index(agent.ctx, docs, agent.docStore)
	if err != nil {
		return 0, fmt.Errorf("error indexing documents: %w", err)
	}
	fmt.Println("✅ Document indexing completed.")
	return len(docs), nil
}

func (agent *RagAgent) SearchSimilarities(query string) ([]string, error) {
	// === SIMILARITY SEARCH ===
	// Create a query document from the user question
	queryDoc := ai.DocumentFromText(query, nil)
	// Create a retriever request with custom options
	request := &ai.RetrieverRequest{
		Query: queryDoc,
	}
	// Retrieve documents relevant to a query
	retrieveResponse, err := agent.documentRetriever.Retrieve(agent.ctx, request)
	if err != nil {
		retrieveResponse = &ai.RetrieverResponse{Documents: []*ai.Document{}}
		return nil, fmt.Errorf("error retrieving documents: %w", err)
	}
	//fmt.Println("📝 Retrieved documents:", retrieveResponse.Documents)

	// Process the retrieved documents
	similarDocuments := []string{}
	for _, doc := range retrieveResponse.Documents {
		//fmt.Println(doc.Metadata, doc.Content[0].Text)
		similarDocuments = append(similarDocuments, doc.Content[0].Text)
	}
	return similarDocuments, nil
}

// calculateEmbeddingDimensionForModel calculates the embedding dimension for a given model by generating a sample embedding.
func calculateEmbeddingDimensionForModel(ctx context.Context, genKitInstance *genkit.Genkit, embedder ai.Embedder) (int, error) {
	res, err := genkit.Embed(
		ctx,
		genKitInstance,
		ai.WithEmbedder(embedder),
		ai.WithTextDocs("Hello World"),
	)
	if err != nil {
		return 0, fmt.Errorf("error when calculating embedding dimension: %w", err)
	}

	return len(res.Embeddings[0].Embedding), nil
}
