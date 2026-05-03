package entradaColectivoInfra

import (
	"DIMISA/src/entradaColectivo/entradaColectivoApp"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type EntradaColectivoController struct {
	GetReporteMensualUC    *entradaColectivoApp.GetReporteMensual
	GetDeficitCronicoUC    *entradaColectivoApp.GetDeficitCronico
	GetComparativoCendisUC *entradaColectivoApp.GetComparativoCendis
}

func NewEntradaColectivoController(
	getReporteMensual *entradaColectivoApp.GetReporteMensual,
	getDeficitCronico *entradaColectivoApp.GetDeficitCronico,
	getComparativoCendis *entradaColectivoApp.GetComparativoCendis,
) *EntradaColectivoController {
	return &EntradaColectivoController{
		GetReporteMensualUC:    getReporteMensual,
		GetDeficitCronicoUC:    getDeficitCronico,
		GetComparativoCendisUC: getComparativoCendis,
	}
}

// POST body: { "id_cendis": 1, "mes": 4, "anio": 2026 }
func (c *EntradaColectivoController) GetReporteMensual(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdCendis int32 `json:"id_cendis"`
		Mes      int   `json:"mes"`
		Anio     int   `json:"anio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"body inválido"}`, http.StatusBadRequest)
		return
	}
	if body.IdCendis == 0 || body.Mes == 0 || body.Anio == 0 {
		http.Error(w, `{"error":"id_cendis, mes y anio son requeridos"}`, http.StatusBadRequest)
		return
	}

	reporte, err := c.GetReporteMensualUC.Execute(body.IdCendis, body.Mes, body.Anio)
	if err != nil {
		http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reporte)
}

// POST body: { "id_cendis": 1, "anio": 2026 }
func (c *EntradaColectivoController) GetDeficitCronico(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdCendis int32 `json:"id_cendis"`
		Anio     int   `json:"anio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"body inválido"}`, http.StatusBadRequest)
		return
	}
	if body.IdCendis == 0 || body.Anio == 0 {
		http.Error(w, `{"error":"id_cendis y anio son requeridos"}`, http.StatusBadRequest)
		return
	}

	detalles, err := c.GetDeficitCronicoUC.Execute(body.IdCendis, body.Anio)
	if err != nil {
		http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles)
}

// POST body: { "mes": 4, "anio": 2026 }
func (c *EntradaColectivoController) GetComparativoCendis(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Mes  int `json:"mes"`
		Anio int `json:"anio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"body inválido"}`, http.StatusBadRequest)
		return
	}
	if body.Mes == 0 || body.Anio == 0 {
		http.Error(w, `{"error":"mes y anio son requeridos"}`, http.StatusBadRequest)
		return
	}

	reportes, err := c.GetComparativoCendisUC.Execute(body.Mes, body.Anio)
	if err != nil {
		http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reportes)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func mesAnioActual() (int, int) {
	now := time.Now()
	return int(now.Month()), now.Year()
}

func parseIntParam(r *http.Request, key string) (int, error) {
	return strconv.Atoi(r.URL.Query().Get(key))
}
