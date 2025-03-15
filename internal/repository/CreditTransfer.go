package repository

import (
	"log"

	db "github.com/alanwade2001/go-sepa-db"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/repository/entity"
)

type CreditTransfer struct {
	persist *db.Persist
}

func NewCreditTransfer(persist *db.Persist) *CreditTransfer {
	ct := &CreditTransfer{
		persist: persist,
	}

	return ct
}

func (s *CreditTransfer) Perist(entity *entity.Transaction) (*entity.Transaction, error) {
	log.Printf("entity: [%v]", entity)
	tx := s.persist.DB.Save(entity)
	err := tx.Error
	return entity, err
}
