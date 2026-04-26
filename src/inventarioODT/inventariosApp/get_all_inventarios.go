package inventariosApp

import (
	inventarioDomain "DIMISA/src/inventarioODT/inventariosDomain"
	inventarioEntity "DIMISA/src/inventarioODT/inventariosDomain/inventariosEntity"
)

type GetAllInventarios struct {
	Repo inventarioDomain.InventarioInterface
}

func (uc *GetAllInventarios) Execute() ([]inventarioEntity.InventarioEntity, error) {
	return uc.Repo.GetAllInventarios()
}
