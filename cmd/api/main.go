package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"google.golang.org/genai"
)

func embeddingCreate(q string) []float64{
	ctx :=context.Background()
	client,err:=genai.NewClient(ctx, &genai.ClientConfig{APIKey: "",Backend: genai.BackendGeminiAPI})
	if err != nil {
        log.Fatal(err)
    }
	content :=[]*genai.Content{
        genai.NewContentFromText(q, genai.RoleUser),
    }
	result, err := client.Models.EmbedContent(ctx,
        "gemini-embedding-2",
        content,
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
	
    values := result.Embeddings[0].Values
	vec :=make([]float64, len(values))
	for i,v := range values{
		vec[i] =float64(v)
	}
	return vec
}

func vectorAverage(v [][]float64) []float64 {
    avg := make([]float64, len(v[0]))  
    for _, r := range v {
        for i, val := range r {
            avg[i] += val
        }
    }
    size := float64(len(v))
    for i, r := range avg {
        avg[i] = r / size
    }
    return avg
}

func main(){
	query :=[]string{
	"How do goroutines work in Go?",
    "What is a channel in Go?",
    "How to handle errors in Go?",
    "What is an interface in Go?",
    "How does garbage collection work in Go?",
    "What is retrieval augmented generation?",
    "How does vector similarity search work?",
    "What is an embedding in machine learning?",
    "How does cosine similarity work?",
    "What is semantic search?",
    "What is CUSUM algorithm?",
    "How does EWMA detect anomalies?",
    "What is statistical process control?",
    "How to detect distribution shift?",
    "What is query drift in LLM systems?",
    "How does HTTP middleware work?",
    "What is observability in distributed systems?",
    "How does Prometheus collect metrics?",
    "What is a sidecar pattern?",
    "How does Docker networking work?",
	}

	var allVectors [][]float64
    for _, q := range query{
        vec := embeddingCreate(q)
        allVectors = append(allVectors, vec)
    }
	baseline := vectorAverage(allVectors)
    fmt.Println("Baseline length:", len(baseline))
    fmt.Println("First 5:", baseline[:5])

	data, _ := json.MarshalIndent(map[string][]float64{
        "embedding": baseline,
    }, "", "  ")
    
    os.WriteFile("baseline_embedding.json", data, 0644)
    fmt.Println("Baseline saved!")


}