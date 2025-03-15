package main

import (
	db "github.com/alanwade2001/go-sepa-db"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/receiver"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/repository"
	"github.com/alanwade2001/go-sepa-engine-ingest/internal/service"
	inf "github.com/alanwade2001/go-sepa-infra"
	q "github.com/alanwade2001/go-sepa-q"
)

type App struct {
	infra    *inf.Infra
	receiver *receiver.Initiation
	Postgres *db.Persist
	PGRepos  *repository.PaymentGroup
	PmtRepos *repository.Payment
	//Iban     *utils.Iban
	TxRepos  *repository.CreditTransfer
	Store    *service.Store
	Listener *q.Listener
}

func NewApp() *App {

	infra := inf.NewInfra()

	Stomp := q.NewStomp()
	Postgres := db.NewPersist()
	PGRepos := repository.NewPaymentGroup(Postgres)
	PmtRepos := repository.NewPayment(Postgres)
	TxRepos := repository.NewCreditTransfer(Postgres)
	Store := service.NewStore()
	//Iban := utils.NewIban()
	Delivery := service.NewDelivery(Stomp)
	Service := service.NewIngestor(PGRepos, PmtRepos, TxRepos, Store, Delivery)
	Receiver := receiver.NewInitiation(Service)

	Listener := q.Newlistener(Stomp, Receiver)

	app := &App{
		infra:    infra,
		Postgres: Postgres,
		PGRepos:  PGRepos,
		PmtRepos: PmtRepos,
		TxRepos:  TxRepos,
		Store:    Store,
		Listener: Listener,
	}

	return app
}

func (a *App) Run() {
	a.Listener.Listen(service.DEST_ENGINE_INITIATION_INGEST)
}

func main() {
	app := NewApp()
	app.Run()
}
