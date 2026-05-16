package clavesApp

import (
	"DIMISA/src/claves/clavesDomain"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type SearchMedClave struct {
	Repo clavesDomain.ClaveInterface
}

func (uc *SearchMedClave) Execute(s string) ([]*claveEntity.ClaveEntity, error) {
	return uc.Repo.SearchMedClave(s)
}
