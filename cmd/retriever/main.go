package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	retrieverAgent "github.com/LouYuanbo1/go-eino-agent/agents/retriever"
	"github.com/LouYuanbo1/go-eino-agent/config"
	"github.com/LouYuanbo1/go-eino-agent/tools/retriever"
	elasticsearchRetriever "github.com/LouYuanbo1/go-eino-agent/tools/retriever/elasticsearch"
	embeddingOllama "github.com/cloudwego/eino-ext/components/embedding/ollama"
	modelOllama "github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	ctx := context.Background()
	// 初始化配置
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Printf("Error initializing config: %v", err)
		return
	}
	retrieverModel, err := modelOllama.NewChatModel(ctx, &modelOllama.ChatModelConfig{
		Model:   "qwen3:1.7b",
		BaseURL: "http://localhost:11434",
	})
	if err != nil {
		fmt.Printf("Error creating chat model: %v", err)
		return
	}
	elasticsearchClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
		Addresses: []string{"http://localhost:9200"},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 30 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			// 跳过TLS验证（仅在开发环境中使用）
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		fmt.Printf("Error creating elasticsearch client: %v", err)
		return
	}
	embeddingModel, err := embeddingOllama.NewEmbedder(ctx, &embeddingOllama.EmbeddingConfig{
		Model:   "nomic-embed-text",
		BaseURL: "http://localhost:11434",
	})
	if err != nil {
		fmt.Printf("Error creating embedder: %v", err)
		return
	}
	typedRetriever := elasticsearchRetriever.NewElasticsearchRetriever(elasticsearchClient, embeddingModel, &elasticsearchRetriever.ElasticsearchRetrieverConfig{
		RetrieverConfig: retriever.RetrieverConfig{
			K:               5,
			IndexName:       "boss_jobs",
			VectorFieldName: "embedding",
		},
		NumCandidates: 100,
	})
	retrieverAgent := retrieverAgent.NewDefaultRetrieverAgent(ctx, retrieverModel, typedRetriever)
	retrieverAgent.OutputMessage(ctx, "帮我寻找北京适合应届毕业生的Go语言高薪工作岗位", true)
}
