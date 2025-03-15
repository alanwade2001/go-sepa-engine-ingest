package service

import (
	"encoding/xml"
	"log"

	"github.com/alanwade2001/go-sepa-engine-data/repository"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/model"
	"github.com/alanwade2001/go-sepa-iso/schema"
)

type Ingestor struct {
	ghRepos  *repository.PaymentGroup
	pmtRepos *repository.Payment
	txRepos  *repository.CreditTransfer
	store    *Store
	//iban     *utils.Iban
	delivery *Delivery
}

func NewIngestor(ghRepos *repository.PaymentGroup, pmtRepos *repository.Payment, txRepos *repository.CreditTransfer, store *Store /*iban *utils.Iban,*/, delivery *Delivery) *Ingestor {
	initiation := &Ingestor{
		ghRepos:  ghRepos,
		pmtRepos: pmtRepos,
		txRepos:  txRepos,
		store:    store,
		//iban:     iban,
		delivery: delivery,
	}

	return initiation
}

func (s *Ingestor) Ingest(mdl *model.PaymentGroup) (err error) {

	log.Printf("ingest: [%v]", mdl)

	var doc *model.Document
	// get the full document
	if doc, err = s.store.GetDocument(mdl.DocID); err != nil {
		return err
	}

	p1Doc := &schema.P1Document{}

	// TODO handle parsing error better
	if err = xml.Unmarshal([]byte(doc.Content), p1Doc); err != nil {
		return err
	}

	mdl.GrpHdr = p1Doc.CstmrCdtTrfInitn.GrpHdr
	if err = s.IngestPaymentGroup(mdl); err != nil {
		return err
	}

	pmtInves := model.NewPayments(p1Doc.CstmrCdtTrfInitn.PmtInf)

	if err = s.IngestPayments(pmtInves); err != nil {
		return err
	}

	if err = s.delivery.PaymentGroupIngested(mdl); err != nil {
		return err
	}

	return nil
}

func (s *Ingestor) IngestPaymentGroup(mdl *model.PaymentGroup) error {

	if pgEnt, err := mdl.ToEntity(); err != nil {
		return err
	} else if pgEnt, err = s.ghRepos.Perist(pgEnt); err != nil {
		return err
	} else {
		mdl.ID = pgEnt.Model.ID
		log.Printf("payment group entity: [%s]", pgEnt.String())
	}

	return nil
}

func (s *Ingestor) IngestPayments(payments []*model.Payment) error {
	for _, payment := range payments {
		if err := s.IngestPayment(payment); err != nil {
			return err
		}
	}

	return nil
}

func (s *Ingestor) IngestPayment(pmtInf *model.Payment) error {

	if pEnt, err := pmtInf.ToEntity(); err != nil {
		return err
	} else if pEnt, err = s.pmtRepos.Perist(pEnt); err != nil {
		return err
	} else {
		pmtInf.ID = pEnt.Model.ID
		log.Printf("payment entity: [%s]", pEnt.String())
		txInves := model.NewCreditTransfers(pmtInf.PmtInf.CdtTrfTxInf)
		if err = s.IngestTransactions(txInves); err != nil {
			return err
		}
	}

	return nil
}

func (s *Ingestor) IngestTransactions(transactions []*model.CreditTransfer) error {
	for _, transaction := range transactions {
		if err := s.IngestTransaction(transaction); err != nil {
			return err
		}
	}

	return nil
}

func (s *Ingestor) IngestTransaction(cdtTrfTxInf *model.CreditTransfer) error {

	if tEnt, err := cdtTrfTxInf.ToEntity(); err != nil {
		return err
	} else if tEnt, err = s.txRepos.Perist(tEnt); err != nil {
		return err
	} else {
		cdtTrfTxInf.ID = tEnt.Model.ID
		log.Printf("transaction entity: [%s]", tEnt.String())
	}

	return nil
}
