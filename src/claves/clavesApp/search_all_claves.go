package clavesApp

import (
	"DIMISA/src/claves/clavesDomain"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type SearchAllClaves struct {
	Repo clavesDomain.ClaveInterface
}

func (uc *SearchAllClaves) Execute(s string) ([]*claveEntity.ClaveEntity, error) {
	return uc.Repo.SearchAllClaves(s)
}
