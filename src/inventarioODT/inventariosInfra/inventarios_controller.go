package inventariosInfra

import (
	"DIMISA/src/inventarioODT/inventariosApp"
	"database/sql"
	"encoding/json"
	"net/http"
)

type InventariosController struct {
	GetInventarioByCendisIDUC *inventariosApp.GetInventarioByCendisID
	GetAllInventariosUC       *inventariosApp.GetAllInventarios
}

func NewInventariosController(
	getInventarioByCendisID *inventariosApp.GetInventarioByCendisID,
	getAllInventarios *inventariosApp.GetAllInventarios,
) *InventariosController {
	return &InventariosController{
		GetInventarioByCendisIDUC: getInventarioByCendisID,
		GetAllInventariosUC:       getAllInventarios,
	}
}

func (c *InventariosController) GetInventarioByCendisID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdCendis int32 `json:"id_cendis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"body inválido"}`, http.StatusBadRequest)
		return
	}

	if body.IdCendis == 0 {
		http.Error(w, `{"error":"id_cendis es requerido"}`, http.StatusBadRequest)
		return
	}

	inventario, err := c.GetInventarioByCendisIDUC.Execute(body.IdCendis)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"inventario no encontrado"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventario)
}

func (c *InventariosController) GetAllInventarios(w http.ResponseWriter, r *http.Request) {
	inventarios, err := c.GetAllInventariosUC.Execute()
	if err != nil {
		http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventarios)
}
