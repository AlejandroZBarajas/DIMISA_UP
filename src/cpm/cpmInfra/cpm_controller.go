package cpmInfra

import (
	"DIMISA/src/cpm/cpmApp"
	"encoding/json"
	"net/http"
)

type CpmController struct {
	GetCpmUC *cpmApp.GetCpm
}

func NewCpmController(getCpm *cpmApp.GetCpm) *CpmController {
	return &CpmController{GetCpmUC: getCpm}
}

func (c *CpmController) GetCpm(w http.ResponseWriter, r *http.Request) {
	result, err := c.GetCpmUC.Execute()
	if err != nil {
		http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
