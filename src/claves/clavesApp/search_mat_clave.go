package clavesApp

import (
	"DIMISA/src/claves/clavesDomain"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type SearchMatClave struct {
	Repo clavesDomain.ClaveInterface
}

func (uc *SearchMatClave) Execute(s string) ([]*claveEntity.ClaveEntity, error) {
	return uc.Repo.SearchMatClave(s)
}
