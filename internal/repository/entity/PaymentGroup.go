package entity

import (
	"fmt"

	"gorm.io/gorm"
)

type PaymentGroup struct {
	Model   *gorm.Model `gorm:"embedded"`
	InitnID uint
	MsgID   string
	CtrlSum float64
	NbOfTxs string
	State   string
	DocID   uint
	GrpHdr  string
}

func (i PaymentGroup) String() string {
	return fmt.Sprintf("ID=[%d], InitnID=[%d], MsgId=[%s], CtrlSum=[%f], NbOfTxs=[%s], State=[%s], DocID=[%d]", i.Model.ID, i.InitnID, i.MsgID, i.CtrlSum, i.NbOfTxs, i.State, i.DocID)
}
