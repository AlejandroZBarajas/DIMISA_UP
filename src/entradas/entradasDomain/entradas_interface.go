package entradasDomain

import (
	"DIMISA/src/entradas/entradasDomain/entradaEntity"
)

type EntradaInterface interface {
	CapturarEntrada(entrada *entradaEntity.EntradaRequest) error
	CapturarInventario(inventario *entradaEntity.InventarioRequest) error
}
