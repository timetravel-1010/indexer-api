package zinc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/timetravel-1010/indexer-api/internal/models"
)

// Config
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

// ZincResponse
type ZincResponse struct {
	Hits struct {
		Hits []models.Hit `json:"hits"`
	} `json:"hits"`
}

const (
	DEFAULT_PAGE_SIZE = 10
)

// DoZincRequest
func DoZincRequest(r *http.Request, query string, c Config) (*ZincResponse, error) {

	index := r.URL.Query().Get("index")
	url := fmt.Sprintf("http://%s:%s/api/%s/_search", c.Host, c.Port, index)
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(query),
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println(resp.StatusCode)

	zr := &ZincResponse{}
	err = json.NewDecoder(resp.Body).Decode(&zr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return zr, nil
}
