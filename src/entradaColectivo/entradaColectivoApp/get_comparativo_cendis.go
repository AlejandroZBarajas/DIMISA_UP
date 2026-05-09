// getComparativoCendis.go
package entradaColectivoApp

import (
	entradaColectivoDomain "DIMISA/src/entradaColectivo/entradaColectivoDomain"
	entity "DIMISA/src/entradaColectivo/entradaColectivoDomain/entradaColectivoEntity"
)

type GetComparativoCendis struct {
	Repo entradaColectivoDomain.EntradaColectivoInterface
}

func (uc *GetComparativoCendis) Execute(mes int, anio int) ([]entity.EntradaColectivoReporte, error) {
	return uc.Repo.GetComparativoCendis(mes, anio)
}
