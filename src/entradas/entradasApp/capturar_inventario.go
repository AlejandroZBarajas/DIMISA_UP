package entradasApp

import (
	"DIMISA/src/entradas/entradasDomain"
	"DIMISA/src/entradas/entradasDomain/entradaEntity"
)

type CapturarInventarioUseCase struct {
	Repo entradasDomain.EntradaInterface
}

func (uc *CapturarInventarioUseCase) Execute(inventario *entradaEntity.InventarioRequest) error {
	return uc.Repo.CapturarInventario(inventario)
}
