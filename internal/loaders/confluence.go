package loaders

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tmc/langchaingo/schema"
)

type ConfluenceSearchResponse struct {
	Results []ConfluencePage
	Links   struct {
		Next string
	} `json:"_links"`
}

type ConfluencePage struct {
	ID    string
	Title string
	Body  struct {
		Storage struct {
			Value string
		}
	}
	Links struct {
		Self  string
		WebUI string
	} `json:"_links"`
}

type ConfluenceLoader struct {
	BaseURL                  string
	APIKey                   string
	Username                 string
	SpaceKeys                []string
	IncludeRestrictedContent bool
	IncludeArchivedContent   bool
	IncludeAttachments       bool
	Limit                    int
	MaxPages                 int

	client *http.Client
}

func NewConfluenceLoader(baseURL, apiKey, username string, spaceKeys []string) *ConfluenceLoader {
	return &ConfluenceLoader{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		Username:  username,
		SpaceKeys: spaceKeys,
		Limit:     50,
		MaxPages:  10000,
		client:    &http.Client{},
	}
}

func (loader *ConfluenceLoader) Load(ctx context.Context) ([]schema.Document, error) {
	var documents []schema.Document
	pages, err := loader.getPages(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("found %v pages\n", len(pages))
	for _, page := range pages {
		documents = append(documents, schema.Document{
			PageContent: page.Body.Storage.Value,
			Metadata: map[string]any{
				"title":  page.Title,
				"id":     page.ID,
				"source": fmt.Sprintf("%s%s", loader.BaseURL, page.Links.WebUI),
			},
		})
	}
	return documents, nil
}

func (loader *ConfluenceLoader) getPages(ctx context.Context) ([]ConfluencePage, error) {
	var pages []ConfluencePage

	for _, spaceKey := range loader.SpaceKeys {
		next := fmt.Sprintf("/rest/api/content?spaceKey=%s&limit=%d&expand=body.storage", spaceKey, loader.Limit)

		for next != "" {
			url := loader.BaseURL + next

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return nil, err
			}

			req.SetBasicAuth(loader.Username, loader.APIKey)

			resp, err := loader.client.Do(req)
			if err != nil {
				if resp.Body != nil {
					resp.Body.Close()
				}
				return nil, err
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				return nil, err
			}

			var result ConfluenceSearchResponse
			if err := json.Unmarshal(body, &result); err != nil {
				resp.Body.Close()
				return nil, err
			}

			pages = append(pages, result.Results...)

			next = result.Links.Next
			resp.Body.Close()
		}
	}

	return pages, nil
}
