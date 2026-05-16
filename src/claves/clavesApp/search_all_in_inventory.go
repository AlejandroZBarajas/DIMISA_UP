package clavesApp

import (
	"DIMISA/src/claves/clavesDomain"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type SearchAllInInventory struct {
	Repo clavesDomain.ClaveInterface
}

func (uc *SearchAllInInventory) Execute(s string, id int32) ([]*claveEntity.ClaveEntity, error) {
	return uc.Repo.SearchAllInInventory(s, id)
}
