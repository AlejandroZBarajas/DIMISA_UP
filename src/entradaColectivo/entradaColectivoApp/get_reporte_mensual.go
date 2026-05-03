// getReporteMensual.go
package entradaColectivoApp

import (
	entradaColectivoDomain "DIMISA/src/entradaColectivo/entradaColectivoDomain"
	entity "DIMISA/src/entradaColectivo/entradaColectivoDomain/entradaColectivoEntity"
)

type GetReporteMensual struct {
	Repo entradaColectivoDomain.EntradaColectivoInterface
}

func (uc *GetReporteMensual) Execute(idCendis int32, mes int, anio int) (entity.EntradaColectivoReporte, error) {
	return uc.Repo.GetReporteMensual(idCendis, mes, anio)
}
