// getDeficitCronico.go
package entradaColectivoApp

import (
	entradaColectivoDomain "DIMISA/src/entradaColectivo/entradaColectivoDomain"
	entity "DIMISA/src/entradaColectivo/entradaColectivoDomain/entradaColectivoEntity"
)

type GetDeficitCronico struct {
	Repo entradaColectivoDomain.EntradaColectivoInterface
}

func (uc *GetDeficitCronico) Execute(idCendis int32, anio int) ([]entity.EntradaColectivoDetalle, error) {
	return uc.Repo.GetDeficitCronico(idCendis, anio)
}
