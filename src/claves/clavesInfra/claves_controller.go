package clavesInfra

import (
	"DIMISA/src/claves/clavesApp"
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ClavesController struct {
	SearchMedUC            *clavesApp.SearchMedClave
	searchMatUC            *clavesApp.SearchMatClave
	SearchMedInInventoryUC *clavesApp.SearchMedInInventory
	searchMatInInventoryUC *clavesApp.SearchMatInInventory
	searchAllInInventoryUC *clavesApp.SearchAllInInventory
	searchAllClavesUC      *clavesApp.SearchAllClaves
}

func NewClaveController(
	searchMed *clavesApp.SearchMedClave,
	searchMat *clavesApp.SearchMatClave,
	searchMedInInventory *clavesApp.SearchMedInInventory,
	searchMatInInventory *clavesApp.SearchMatInInventory,
	searchAllInInventory *clavesApp.SearchAllInInventory,
	searchAllClaves *clavesApp.SearchAllClaves,
) *ClavesController {
	return &ClavesController{
		SearchMedUC:            searchMed,
		searchMatUC:            searchMat,
		SearchMedInInventoryUC: searchMedInInventory,
		searchMatInInventoryUC: searchMatInInventory,
		searchAllInInventoryUC: searchAllInInventory,
		searchAllClavesUC:      searchAllClaves,
	}
}

type SearchResponse struct {
	Success bool                       `json:"success"`
	Data    []*claveEntity.ClaveEntity `json:"data,omitempty"`
	Message string                     `json:"message,omitempty"`
	Count   int                        `json:"count"`
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SearchResponse{
		Success: false,
		Message: message,
		Count:   0,
	})
}

// parseSearchQuery extrae y valida el parámetro q
func parseSearchQuery(w http.ResponseWriter, r *http.Request) (string, bool) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		sendError(w, "El parámetro 'q' es requerido", http.StatusBadRequest)
		return "", false
	}
	if len(query) < 2 {
		sendError(w, "La búsqueda debe tener al menos 2 caracteres", http.StatusBadRequest)
		return "", false
	}
	return query, true
}

// parseCendisID extrae y valida el parámetro cendis
func parseCendisID(w http.ResponseWriter, r *http.Request) (int32, bool) {
	cendisStr := r.URL.Query().Get("cendis")
	if cendisStr == "" {
		sendError(w, "El parámetro 'cendis' es requerido", http.StatusBadRequest)
		return 0, false
	}
	id, err := strconv.Atoi(cendisStr)
	if err != nil || id <= 0 {
		sendError(w, "El parámetro 'cendis' debe ser un número válido", http.StatusBadRequest)
		return 0, false
	}
	return int32(id), true
}

// buildResponse arma el SearchResponse estándar
func buildResponse(results []*claveEntity.ClaveEntity, emptyMsg string) SearchResponse {
	response := SearchResponse{
		Success: true,
		Data:    results,
		Count:   len(results),
	}
	if len(results) == 0 {
		response.Message = emptyMsg
	}
	return response
}

// ─── MEDICAMENTOS ────────────────────────────────────────────────────────────

func (c *ClavesController) SearchMedClave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := parseSearchQuery(w, r)
	if !ok {
		return
	}

	results, err := c.SearchMedUC.Execute(query)
	if err != nil {
		sendError(w, "Error al buscar medicamentos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildResponse(results, "No se encontraron medicamentos"))
}

func (c *ClavesController) SearchMedInInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := parseSearchQuery(w, r)
	if !ok {
		return
	}

	cendisID, ok := parseCendisID(w, r)
	if !ok {
		return
	}

	results, err := c.SearchMedInInventoryUC.Execute(query, cendisID)
	if err != nil {
		sendError(w, "Error al buscar medicamentos en inventario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildResponse(results, "No se encontraron medicamentos en el inventario"))
}

// ─── MATERIAL DE CURACION ────────────────────────────────────────────────────

func (c *ClavesController) SearchMatClave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := parseSearchQuery(w, r)
	if !ok {
		return
	}

	results, err := c.searchMatUC.Execute(query)
	if err != nil {
		sendError(w, "Error al buscar material de curación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildResponse(results, "No se encontró material de curación"))
}

func (c *ClavesController) SearchMatInInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := parseSearchQuery(w, r)
	if !ok {
		return
	}

	cendisID, ok := parseCendisID(w, r)
	if !ok {
		return
	}

	results, err := c.searchMatInInventoryUC.Execute(query, cendisID)
	if err != nil {
		sendError(w, "Error al buscar material de curación en inventario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildResponse(results, "No se encontró material de curación en el inventario"))
}

func (c *ClavesController) SearchAllInInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := parseSearchQuery(w, r)
	if !ok {
		return
	}

	cendisID, ok := parseCendisID(w, r)
	if !ok {
		return
	}

	results, err := c.searchAllInInventoryUC.Execute(query, cendisID)
	if err != nil {
		sendError(w, "Error al buscar en inventario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildResponse(results, "No se encontraron resultados en el inventario"))
}

func (c *ClavesController) SearchAllClaves(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := parseSearchQuery(w, r)
	if !ok {
		return
	}

	results, err := c.searchAllClavesUC.Execute(query)
	if err != nil {
		sendError(w, "Error al buscar claves: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildResponse(results, "No se encontraron resultados"))
}
