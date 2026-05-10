package clavesApp

import (
	"DIMISA/src/claves/clavesDomain"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type SearchMedInInventory struct {
	Repo clavesDomain.ClaveInterface
}

func (uc *SearchMedInInventory) Execute(s string, id int32) ([]*claveEntity.ClaveEntity, error) {
	return uc.Repo.SearchMedInInventory(s, id)
}
