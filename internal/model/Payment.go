package model

import (
	"encoding/xml"

	"github.com/alanwade2001/go-sepa-engine-data/repository/entity"
	"github.com/alanwade2001/go-sepa-iso/gen"
	"gorm.io/gorm"
)

type Payment struct {
	ID       uint
	PmtInfId string
	CtrlSum  float64
	NbOfTxs  string
	DbtrAcc  *Account
	PmtInf   *gen.PaymentInstructionInformation3
}

func NewPayment(pmtInf *gen.PaymentInstructionInformation3) *Payment {

	pmt := &Payment{
		PmtInfId: pmtInf.PmtInfId,
		CtrlSum:  pmtInf.CtrlSum,
		NbOfTxs:  pmtInf.NbOfTxs,
		DbtrAcc:  &Account{Nm: pmtInf.Dbtr.Nm, Iban: pmtInf.DbtrAcct.Id.IBAN, Bic: pmtInf.DbtrAgt.FinInstnId.BIC},
		PmtInf:   pmtInf,
	}

	return pmt
}

func NewPayments(pmtInves []*gen.PaymentInstructionInformation3) []*Payment {
	pmts := []*Payment{}

	for _, pmtInf := range pmtInves {
		pmt := NewPayment(pmtInf)
		pmts = append(pmts, pmt)
	}

	return pmts
}

func (p Payment) ToEntity() (ent *entity.Payment, err error) {

	if bytes, err := xml.Marshal(p.PmtInf); err != nil {
		return nil, err
	} else {

		ent = &entity.Payment{
			Model: &gorm.Model{
				ID: p.ID,
			},
			PmtInfID: p.PmtInfId,
			CtrlSum:  p.CtrlSum,
			NbOfTxs:  p.NbOfTxs,
			Nm:       p.DbtrAcc.Nm,
			Iban:     p.DbtrAcc.Iban,
			Bic:      p.DbtrAcc.Bic,
			PmtInf:   string(bytes),
		}

		return ent, nil
	}

}
