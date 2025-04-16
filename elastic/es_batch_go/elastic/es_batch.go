package elastic

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	jsoniter "github.com/json-iterator/go"
)

type EsBatchInsert struct {
	host     string
	userName string
	password string
	client   *elasticsearch.Client
}
type Document struct {
	ID     string         `json:"id"`
	Index  string         `json:"index"`
	Source map[string]any `json:"source"`
}

func NewEsBatchInsert(host, userName, password string) (*EsBatchInsert, error) {
	cert, _ := os.ReadFile("/Users/liupeng/Documents/study/elasticsearch-8.17.0/config/certs/http_ca.crt")
	cfg := elasticsearch.Config{
		Username: userName,
		// Password:  "xpE4DQGWE9bCkoj7WXYE",
		Password: password,
		// Addresses: []string{"https://localhost:9200"},
		Addresses: []string{host},
		CACert:    cert,
	}
	client, err := elasticsearch.NewClient(cfg)
	return &EsBatchInsert{
		host:     host,
		userName: userName,
		password: password,
		client:   client,
	}, err
}

func (b *EsBatchInsert) BatchInSert(docs []Document) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	buf := bytes.Buffer{}
	for _, doc := range docs {
		meta := map[string]any{
			"index": map[string]any{
				"_index": doc.Index,
				"_id":    doc.ID,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
		if err := json.NewEncoder(&buf).Encode(doc.Source); err != nil {
			return err
		}
	}
	// resp, err := client.Bulk(&buf, client.Bulk.WithContext(context.Background()))
	req := esapi.BulkRequest{
		Body:   &buf,
		Pretty: true, // 格式化响应
	}
	resp, err := req.Do(context.Background(), b.client)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Println("error status code: ", resp.StatusCode)
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}
	return nil
}
