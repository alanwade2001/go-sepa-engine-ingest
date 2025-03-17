package service

import (
	"encoding/xml"
	"log/slog"

	"github.com/alanwade2001/go-sepa-engine-data/model"
	"github.com/alanwade2001/go-sepa-engine-data/repository"
	"github.com/alanwade2001/go-sepa-engine-data/repository/entity"
	"github.com/alanwade2001/go-sepa-iso/schema"
)

type Ingestor struct {
	ghRepos  *repository.PaymentGroup
	pmtRepos *repository.Payment
	txRepos  *repository.Transaction
	store    *Store
	iban     *Iban
	delivery *Delivery
}

func NewIngestor(ghRepos *repository.PaymentGroup, pmtRepos *repository.Payment, txRepos *repository.Transaction, store *Store, delivery *Delivery, iban *Iban) *Ingestor {
	initiation := &Ingestor{
		ghRepos:  ghRepos,
		pmtRepos: pmtRepos,
		txRepos:  txRepos,
		store:    store,
		iban:     iban,
		delivery: delivery,
	}

	return initiation
}

func (s *Ingestor) Ingest(mdl *model.PaymentGroup) (err error) {

	slog.Info("ingest", "model", mdl)

	var doc *model.Document
	// get the full document
	if doc, err = s.store.GetDocument(mdl.DocID); err != nil {
		return err
	}

	p1Doc := &schema.Pain001Document{}

	// TODO handle parsing error better
	if err = xml.Unmarshal([]byte(doc.Content), p1Doc); err != nil {
		return err
	}

	var ePg *entity.PaymentGroup

	mdl.GrpHdr = p1Doc.CstmrCdtTrfInitn.GrpHdr
	if ePg, err = s.IngestPaymentGroup(mdl); err != nil {
		return err
	}

	pmtInves := model.NewPayments(p1Doc.CstmrCdtTrfInitn.PmtInf)

	if _, err = s.IngestPayments(ePg, pmtInves); err != nil {
		return err
	}

	if err = s.delivery.PaymentGroupIngested(mdl); err != nil {
		return err
	}

	return nil
}

func (s *Ingestor) IngestPaymentGroup(mdl *model.PaymentGroup) (*entity.PaymentGroup, error) {

	if pgEnt, err := mdl.ToEntity(); err != nil {
		return nil, err
	} else if pgEnt, err = s.ghRepos.Perist(pgEnt); err != nil {
		return nil, err
	} else {
		mdl.ID = pgEnt.Model.ID
		slog.Info("payment group persisted", "entity", pgEnt.String())
		return pgEnt, nil
	}

}

func (s *Ingestor) IngestPayments(ePg *entity.PaymentGroup, payments []*model.Payment) ([]*entity.Payment, error) {
	entities := make([]*entity.Payment, len(payments))

	for _, payment := range payments {
		if entity, err := s.IngestPayment(ePg, payment); err != nil {
			return nil, err
		} else {
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (s *Ingestor) IngestPayment(ePg *entity.PaymentGroup, pmtInf *model.Payment) (*entity.Payment, error) {

	if pEnt, err := pmtInf.ToEntity(); err != nil {
		return nil, err
	} else {
		pEnt.PaymentGroupID = ePg.Model.ID
		if pEnt, err = s.pmtRepos.Perist(pEnt); err != nil {
			return nil, err
		} else {
			pmtInf.ID = pEnt.Model.ID
			slog.Info("persisted payment", "entity", pEnt.String())
			txInves := model.NewCreditTransfers(pmtInf.PmtInf.CdtTrfTxInf)
			if _, err = s.IngestTransactions(pEnt, txInves); err != nil {
				return nil, err
			}
		}

		return pEnt, nil
	}

}

func (s *Ingestor) IngestTransactions(eP *entity.Payment, transactions []*model.CreditTransfer) ([]*entity.Transaction, error) {
	entities := make([]*entity.Transaction, len(transactions))
	for _, transaction := range transactions {
		if entity, err := s.IngestTransaction(eP, transaction); err != nil {
			return nil, err
		} else {
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (s *Ingestor) IngestTransaction(eP *entity.Payment, cdtTrfTxInf *model.CreditTransfer) (*entity.Transaction, error) {

	if tEnt, err := cdtTrfTxInf.ToEntity(); err != nil {
		return nil, err
	} else {
		tEnt.PaymentID = eP.Model.ID
		if tEnt, err = s.txRepos.Perist(tEnt); err != nil {
			return nil, err
		} else {
			cdtTrfTxInf.ID = tEnt.Model.ID
			slog.Info("persisted transaction", "entity", tEnt.String())
			return tEnt, nil
		}
	}

}
