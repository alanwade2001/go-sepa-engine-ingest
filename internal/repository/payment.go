package repository

import (
	"log"

	db "github.com/alanwade2001/go-sepa-db"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/repository/entity"
)

type Payment struct {
	persist *db.Persist
}

func NewPayment(persist *db.Persist) *Payment {
	payment := &Payment{
		persist: persist,
	}

	return payment
}

func (s *Payment) Perist(entity *entity.Payment) (*entity.Payment, error) {
	log.Printf("entity: [%v]", entity)

	tx := s.persist.DB.Save(entity)
	err := tx.Error

	return entity, err
}
