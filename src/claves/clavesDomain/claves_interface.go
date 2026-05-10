package clavesDomain

import (
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
)

type ClaveInterface interface {
	SearchMedClave(s string) ([]*claveEntity.ClaveEntity, error)
	SearchMatClave(s string) ([]*claveEntity.ClaveEntity, error)
	SearchMedInInventory(s string, id int32) ([]*claveEntity.ClaveEntity, error)
	SearchMatInInventory(s string, id int32) ([]*claveEntity.ClaveEntity, error)
}
