package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
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

	log.Printf("get document: [%s]", url)

	response, err := s.client.Get(url)

	if err != nil {
		log.Printf("failed to post, %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New("failed to get document")
		log.Printf("failed to get, %v", err)
		return nil, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(doc)
	if err != nil {
		log.Printf("failed to decode, %v", err)
		return nil, err
	}

	log.Printf("Document ID: %d", doc.ID)

	return doc, nil
}
