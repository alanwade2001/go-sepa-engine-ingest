package service

import (
	"encoding/json"

	"github.com/alanwade2001/go-sepa-engine-data/model"
	q "github.com/alanwade2001/go-sepa-q"
)

var DEST_ENGINE_INITIATION_INGEST string = "topic:portal.initiation.approved"
var DEST_ENGINE_PAYMENT_GROUP_INGESTED string = "queue:portal.payment_group.ingested"

type Delivery struct {
	Stomp *q.Stomp
}

func NewDelivery(Stomp *q.Stomp) *Delivery {
	i := &Delivery{
		Stomp: Stomp,
	}

	return i
}

func (d *Delivery) PaymentGroupIngested(pg *model.PaymentGroup) error {
	if bytes, err := json.Marshal(pg); err != nil {
		return err
	} else {
		if err = d.Stomp.SendMessage(DEST_ENGINE_PAYMENT_GROUP_INGESTED, bytes); err != nil {
			return err
		}
	}

	return nil
}
