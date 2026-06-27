package extractor

import (
	"QueryGuard/internal/cache"
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"

	"google.golang.org/genai"
)

type Extactor struct{
	Cache cache.Cache
	Baseline []float64
	Client *genai.Client
	Ctx context.Context
}

func  NewExtractor(baselinePath string)(*Extactor,error){
	ctx := context.Background()
	client,err :=genai.NewClient(ctx,&genai.ClientConfig{APIKey:  os.Getenv("GEMINI_API_KEY"),
        Backend: genai.BackendGeminiAPI,
    })
	if(err !=nil){
		return nil,err
	}

	data,err :=os.ReadFile(baselinePath)
	if(err!=nil){
		return nil,err
	}
	var stored struct {
        Embedding []float64 `json:"embedding"`
    }
    if err := json.Unmarshal(data, &stored); err != nil {
        return nil, err
    }
	c :=cache.NewCache()

	return &Extactor{
		Cache: c,
		Baseline: stored.Embedding,
		Client: client,
		Ctx: ctx,
	},nil
}


func (e *Extactor) createEmbedding(query string ) ([]float64, error) {

    content := []*genai.Content{
        genai.NewContentFromText(query, genai.RoleUser),
    }
    result, err := e.Client.Models.EmbedContent(e.Ctx,
        "gemini-embedding-2",
        content,
        nil,
    )
    if err != nil {
        return nil, err
    }
    values := result.Embeddings[0].Values
    vec := make([]float64, len(values))
    for i, v := range values {
        vec[i] = float64(v)
    }
    return vec, nil
}
func (e *Extactor)GetEmbedding(query string)[]float64{

	emb,result :=e.Cache.Get(query)
	if(result){
		return emb
	}
	emb,_ = e.createEmbedding(query)
	e.Cache.Check(query, emb)
	return emb
}	


func CosineSimilarity(a,b []float64)(float64,error){
	if len(a) != len(b) {
		return 0, errors.New("vectors must have the same length")
	}
	if len(a) == 0 {
		return 0, errors.New("vectors cannot be empty")
	}
	var dot ,normA,normB float64
	for i:= range len(a){
		dot +=(a[i]*b[i])
		normA +=a[i]*a[i]
		normB +=b[i]*b[i]
	}
	if normA == 0 || normB == 0 {
		return 0, errors.New("cannot calculate similarity for zero-magnitude vectors")
	}
	return dot/(math.Sqrt(normA) * math.Sqrt(normB)),nil
}


func (ext *Extactor) Extract (query string) float64{
	queryEmbedding := ext.GetEmbedding(query)
	similarityScore,_ := CosineSimilarity(queryEmbedding,ext.Baseline)
	return similarityScore
}