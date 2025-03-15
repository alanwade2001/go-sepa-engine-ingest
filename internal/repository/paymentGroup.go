package repository

import (
	"log"

	db "github.com/alanwade2001/go-sepa-db"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/repository/entity"
)

type PaymentGroup struct {
	persist *db.Persist
}

func NewPaymentGroup(persist *db.Persist) *PaymentGroup {
	initiation := &PaymentGroup{
		persist: persist,
	}

	return initiation
}

func (s *PaymentGroup) Perist(entity *entity.PaymentGroup) (*entity.PaymentGroup, error) {
	log.Printf("entity: [%v]", entity)
	tx := s.persist.DB.Save(entity)
	err := tx.Error
	return entity, err
}

func (s *PaymentGroup) FindByID(id string) (*entity.PaymentGroup, error) {
	paymentGroup := &entity.PaymentGroup{}
	if err := s.persist.DB.First(paymentGroup, id).Error; err != nil {
		return nil, err
	}

	return paymentGroup, nil
}
