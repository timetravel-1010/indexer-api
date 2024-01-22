package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Email struct {
	MessageID               string   `json:"messageId"`
	Date                    string   `json:"date"`
	From                    string   `json:"from"`
	To                      []string `json:"to"`
	CC                      []string `json:"cc"`
	BCC                     []string `json:"bcc"`
	Subject                 string   `json:"subject"`
	MimeVersion             string   `json:"mimeVersion"`
	ContentType             string   `json:"contentType"`
	ContentTransferEncoding string   `json:"contentTransferEncoding"`
	XFrom                   string   `json:"xFrom"`
	XTo                     []string `json:"xTo"`
	Xcc                     []string `json:"xcc"`
	Xbcc                    []string `json:"xbcc"`
	XFolder                 string   `json:"xFolder"`
	XOrigin                 string   `json:"xOrigin"`
	XFileName               string   `json:"xFileName"`
	Body                    string   `json:"body"`
}

type EmailResponse struct {
}

type Response struct {
	Hits []Hit `json:"hits"`
}

type Hit struct {
	Index     string  `json:"_index"`
	Type      string  `json:"_type"`
	Id        string  `json:"_id"`
	Score     float64 `json:"_score"`
	Timestamp string  `json:"@timestamp"`
	Source    Source  `json:"_source"`
}

type Source struct {
	Email     Email  `json:"email"`
	Timestamp string `json:"@timestamp"`
	Path      string `json:"path"`
}

type EmailHandler struct {
}

type Config struct {
	host string
	port string
}

func (eh EmailHandler) GetEmail(w http.ResponseWriter, r *http.Request) {

	query := []byte(`{
	    "search_type": "matchall",
	    "from": 0,
	    "max_results": %s,
	    "_source": []
	}`)
	c := Config{
		host: "localhost",
		port: "4080",
	}
	index := r.URL.Query().Get("index")
	page := r.URL.Query().Get("page")
	url := fmt.Sprintf("http://%s:%s/api/%s/_search", c.host, c.port, index)

	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(fmt.Sprintf(string(query), page)),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	log.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		Hits struct {
			Hits []Hit `json:"hits"`
		} `json:"hits"`
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err, "this")
	}
	//fmt.Println(res.Hits.Hits)
	//fmt.Println(string(body))
	json.NewEncoder(w).Encode(res.Hits.Hits)
}
