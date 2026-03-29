package entradasInfra

import (
	"DIMISA/src/entradas/entradasApp"
	"DIMISA/src/entradas/entradasDomain/entradaEntity"
	"encoding/json"
	"net/http"
)

type EntradasController struct {
	CapturarEntradaUC    *entradasApp.CapturarEntradaUseCase
	CapturarInventarioUC *entradasApp.CapturarInventarioUseCase
}

func NewEntradaController(capturar *entradasApp.CapturarEntradaUseCase, inventario *entradasApp.CapturarInventarioUseCase) *EntradasController {
	return &EntradasController{
		CapturarEntradaUC:    capturar,
		CapturarInventarioUC: inventario,
	}
}

func (c *EntradasController) CapturarEntrada(w http.ResponseWriter, r *http.Request) {
	var entrada entradaEntity.EntradaRequest

	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := c.CapturarEntradaUC.Execute(&entrada); err != nil {
		http.Error(w, "Error al capturar entrada", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Entrada capturada correctamente"})
}

func (c *EntradasController) CapturarInventario(w http.ResponseWriter, r *http.Request) {
	var inventario entradaEntity.InventarioRequest

	if err := json.NewDecoder(r.Body).Decode(&inventario); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if inventario.Id_cendis == 0 || len(inventario.Detalles) == 0 {
		http.Error(w, "id_cendis y detalles son requeridos", http.StatusBadRequest)
		return
	}

	if err := c.CapturarInventarioUC.Execute(&inventario); err != nil {
		http.Error(w, "Error al cargar inventario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Inventario cargado correctamente"})
}
