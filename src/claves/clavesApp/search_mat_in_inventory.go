package clavesApp

import (
	"DIMISA/src/claves/clavesDomain"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type SearchMatInInventory struct {
	Repo clavesDomain.ClaveInterface
}

func (uc *SearchMatInInventory) Execute(s string, id int32) ([]*claveEntity.ClaveEntity, error) {
	return uc.Repo.SearchMatInInventory(s, id)
}
