package tasks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/rhajizada/gazette/internal/config"
	"golang.org/x/net/html"
)

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
