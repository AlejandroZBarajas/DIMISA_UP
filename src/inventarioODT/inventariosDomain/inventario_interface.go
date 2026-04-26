package inventarioDomain

import (
	inventarioEntity "DIMISA/src/inventarioODT/inventariosDomain/inventariosEntity"
)

type InventarioInterface interface {
	GetInventarioByCendisID(id_cendis int32) (inventarioEntity.InventarioEntity, error)
	GetAllInventarios() ([]inventarioEntity.InventarioEntity, error)
}
