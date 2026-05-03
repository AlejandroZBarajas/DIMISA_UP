package entradaColectivoDomain

import entity "DIMISA/src/entradaColectivo/entradaColectivoDomain/entradaColectivoEntity"

type EntradaColectivoInterface interface {
	GetReporteMensual(idCendis int32, mes int, anio int) (entity.EntradaColectivoReporte, error)
	GetDeficitCronico(idCendis int32, anio int) ([]entity.EntradaColectivoDetalle, error)
	GetComparativoCendis(mes int, anio int) ([]entity.EntradaColectivoReporte, error)
}
