package model

import (
	"encoding/xml"
	"fmt"

	"github.com/alanwade2001/go-sepa-engine-data/repository/entity"
	"github.com/alanwade2001/go-sepa-iso/gen"
	"gorm.io/gorm"
)

type PaymentGroup struct {
	ID      uint
	InitnID uint
	MsgID   string
	CtrlSum float64
	NbOfTxs string
	State   string
	DocID   uint
	GrpHdr  *gen.GroupHeader32
}

func (i PaymentGroup) String() string {
	return fmt.Sprintf("ID=[%d], InitnID=[%d], MsgId=[%s], CtrlSum=[%f], NbOfTxs=[%s], State=[%s], DocID=[%d]", i.ID, i.InitnID, i.MsgID, i.CtrlSum, i.NbOfTxs, i.State, i.DocID)
}

func (p PaymentGroup) ToEntity() (ent *entity.PaymentGroup, err error) {

	if bytes, err := xml.Marshal(p.GrpHdr); err != nil {
		return nil, err
	} else {

		ent = &entity.PaymentGroup{
			Model: &gorm.Model{
				ID: p.ID,
			},
			InitnID: p.InitnID,
			CtrlSum: p.CtrlSum,
			NbOfTxs: p.NbOfTxs,
			MsgID:   p.MsgID,
			State:   p.State,
			DocID:   p.DocID,
			GrpHdr:  string(bytes),
		}

		return ent, nil
	}

}
