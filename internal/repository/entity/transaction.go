package entity

import (
	"fmt"

	"gorm.io/gorm"
)

type Transaction struct {
	Model       *gorm.Model `gorm:"embedded"`
	EndToEndID  string
	Amt         float64
	Nm          string
	Iban        string
	Bic         string
	CdtTrfTxInf string
}

func (t Transaction) String() string {
	return fmt.Sprintf("ID=[%d], EndToEndID=[%s], Amt=[%f], Nm=[%s], Iban=[%s], bic=[%s]", t.Model.ID, t.EndToEndID, t.Amt, t.Nm, t.Iban, t.Bic)
}
