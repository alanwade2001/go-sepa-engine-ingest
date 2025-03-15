package receiver

import (
	"encoding/json"
	"log"

	"github.com/alanwade2001/go-sepa-engine-data/model"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/receiver/message"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/service"
)

type Initiation struct {
	Service *service.Ingestor
}

func NewInitiation(Service *service.Ingestor) *Initiation {
	i := &Initiation{
		Service: Service,
	}

	return i
}

func (i *Initiation) Process(body []byte) error {
	text := string(body)
	msg := &message.Initiation{}
	if err := json.Unmarshal([]byte(text), msg); err != nil {
		log.Println(err)
	} else {
		log.Printf("msg:[%v]", msg.String())
		model := &model.PaymentGroup{}

		if err = message.Map(model, msg); err != nil {
			log.Println(err)
		} else {
			log.Printf("model:[%v]", model.String())
			if err = i.Service.Ingest(model); err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}
