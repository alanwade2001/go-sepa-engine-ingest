package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log/slog"

	"net/http"

	"strconv"

	"time"

	"github.com/alanwade2001/go-sepa-engine-data/model"
	utils "github.com/alanwade2001/go-sepa-utils"
)

type Store struct {
	Address string
	client  http.Client
}

func NewStore() *Store {

	address := utils.Getenv("DOCS_ADDRESS", "https://0.0.0.0:8443")

	s := &Store{
		Address: address,
		client: http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	return s
}

func (s *Store) GetDocument(docId uint) (doc *model.Document, err error) {
	doc = &model.Document{}

	id := strconv.FormatUint(uint64(docId), 10)
	url := s.Address + "/documents/" + id

	slog.Info("get document", "url", url)

	response, err := s.client.Get(url)

	if err != nil {
		slog.Error("failed to post", "error", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New("failed to get document")
		slog.Error("failed to get", "error", err)
		return nil, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(doc)
	if err != nil {
		slog.Error("failed to decode", "error", err)
		return nil, err
	}

	slog.Info("Got Document:", "ID", doc.ID)

	return doc, nil
}
