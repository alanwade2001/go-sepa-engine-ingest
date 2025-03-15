package message

import (
	"fmt"

	"github.com/alanwade2001/go-sepa-engine-ingest/internal/model"
)

type Initiation struct {
	ID      uint    `json:"id"`
	MsgID   string  `json:"msgId"`
	CtrlSum float64 `json:"ctrlSum"`
	NbOfTxs string  `json:"NbOfTxs"`
	State   string  `json:"state"`
	DocID   uint    `json:"docId"`
}

func (i Initiation) String() string {
	return fmt.Sprintf("ID=[%d], MsgId=[%s], CtrlSum=[%f], NbOfTxs=[%s], State=[%s], DocID=[%d]", i.ID, i.MsgID, i.CtrlSum, i.NbOfTxs, i.State, i.DocID)
}

func Map(mdl *model.PaymentGroup, msg *Initiation) error {

	mdl.CtrlSum = msg.CtrlSum
	mdl.DocID = msg.DocID
	mdl.InitnID = msg.ID
	mdl.MsgID = msg.MsgID
	mdl.NbOfTxs = msg.NbOfTxs
	mdl.State = msg.State

	return nil
}
