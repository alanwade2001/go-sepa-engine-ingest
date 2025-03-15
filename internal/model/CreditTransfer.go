package model

import (
	"encoding/xml"

	"github.com/alanwade2001/go-sepa-engine-data/repository/entity"
	"github.com/alanwade2001/go-sepa-iso/gen"
	"gorm.io/gorm"
)

type CreditTransfer struct {
	ID          uint
	EndToEndID  string
	Amt         float64
	CdtrAcc     *Account
	CdtTrfTxInf *gen.CreditTransferTransactionInformation10
}

func NewCreditTransfer(txInf *gen.CreditTransferTransactionInformation10) *CreditTransfer {

	tx := &CreditTransfer{
		EndToEndID:  txInf.PmtId.EndToEndId,
		Amt:         txInf.Amt.InstdAmt.Value,
		CdtrAcc:     &Account{Nm: txInf.Cdtr.Nm, Iban: txInf.CdtrAcct.Id.IBAN, Bic: txInf.CdtrAgt.FinInstnId.BIC},
		CdtTrfTxInf: txInf,
	}

	return tx
}

func NewCreditTransfers(txInves []*gen.CreditTransferTransactionInformation10) []*CreditTransfer {
	txs := []*CreditTransfer{}

	for _, txInf := range txInves {
		tx := NewCreditTransfer(txInf)
		txs = append(txs, tx)
	}

	return txs
}

func (t CreditTransfer) ToEntity() (ent *entity.Transaction, err error) {

	if bytes, err := xml.Marshal(t.CdtTrfTxInf); err != nil {
		return nil, err
	} else {

		ent = &entity.Transaction{
			Model: &gorm.Model{
				ID: t.ID,
			},
			EndToEndID:  t.EndToEndID,
			Amt:         t.Amt,
			Nm:          t.CdtrAcc.Nm,
			Iban:        t.CdtrAcc.Iban,
			Bic:         t.CdtrAcc.Bic,
			CdtTrfTxInf: string(bytes),
		}

		return ent, nil
	}

}
