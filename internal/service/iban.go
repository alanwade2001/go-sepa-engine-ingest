package service

import "github.com/jbub/banking/iban"

type Iban struct {
}

func NewIban() *Iban {
	iban := &Iban{}

	return iban
}

func (i Iban) CheckIban(s string) error {
	return iban.Validate(s)
}
