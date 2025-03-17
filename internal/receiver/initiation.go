package receiver

import (
	"encoding/json"
	"log/slog"

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
		slog.Error("unmarshall", "error", err, "body", body)
	} else {
		slog.Info("received msg", "msg", msg)
		model := &model.PaymentGroup{}

		if err = message.Map(model, msg); err != nil {
			slog.Error("failed to map", "error", err)
		} else {
			slog.Info("mapped", "model", model.String())
			if err = i.Service.Ingest(model); err != nil {
				slog.Error("failed to ingest", "error", err)
			}
		}
	}

	return nil
}
