package inventariosApp

import (
	inventarioDomain "DIMISA/src/inventarioODT/inventariosDomain"
	inventarioEntity "DIMISA/src/inventarioODT/inventariosDomain/inventariosEntity"
)

type GetInventarioByCendisID struct {
	Repo inventarioDomain.InventarioInterface
}

func (uc *GetInventarioByCendisID) Execute(id_cendis int32) (inventarioEntity.InventarioEntity, error) {
	return uc.Repo.GetInventarioByCendisID(id_cendis)
}
