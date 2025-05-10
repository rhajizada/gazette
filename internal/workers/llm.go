package workers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cdipaolo/goml/cluster"
	"github.com/ollama/ollama/api"
	"github.com/pgvector/pgvector-go"
	"github.com/rhajizada/gazette/internal/config"
	"golang.org/x/net/html"
)

const (
	ClusterCount   = 64
	ClusterMaxIter = 512
)

type Centroid struct {
	Vector      *pgvector.Vector
	MemberCount int32
}

func GetOllamaClient(cfg *config.OllamaConfig) (*api.Client, error) {
	baseURL, err := url.Parse(cfg.BaseUrl)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(baseURL, http.DefaultClient)
	return client, nil
}

func InitModels(c *api.Client, cfg *config.OllamaConfig) error {
	model := cfg.EmbeddingsModel
	ctx := context.Background()
	listReponse, err := c.List(ctx)
	if err != nil {
		return err
	}
	found := false
	for _, v := range listReponse.Models {
		if model == v.Model {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("model '%s' not found", model)
	}
	return nil
}

func ExtractTextFromHTML(input string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(input))
	var result strings.Builder

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return result.String()
		case html.TextToken:
			result.WriteString(tokenizer.Token().Data)
		}
	}
}

func ClusterizeEmbeddings(points [][]float64) ([]Centroid, error) {
	// If there are no points, nothing to cluster
	if len(points) == 0 {
		return nil, nil
	}

	k := ClusterCount
	if len(points) < k {
		k = len(points)
	}

	model := cluster.NewKMeans(k, ClusterMaxIter, points)
	if err := model.Learn(); err != nil {
		return nil, err
	}

	guesses := model.Guesses()
	cents := model.Centroids

	// count how many points landed in each cluster
	counts := make([]int, len(cents))
	for _, g := range guesses {
		if g >= 0 && g < len(counts) {
			counts[g]++
		}
	}

	// build your Centroid structs
	centroids := make([]Centroid, len(cents))
	for i, cent := range cents {
		vec := float64sToVector(cent)
		centroids[i] = Centroid{
			Vector:      &vec,
			MemberCount: int32(counts[i]),
		}
	}

	return centroids, nil
}
