// src/claves/clavesInfra/claveRepository.go
package clavesInfra

import (
	claveEntity "DIMISA/src/claves/clavesDomain/entity"
	"DIMISA/src/core/config"
	"DIMISA/src/core/utils"
	"database/sql"
	"fmt"
)

type ClaveRepository struct {
	DB *sql.DB
}

// ─── MEDICAMENTOS ────────────────────────────────────────────────────────────

func (r *ClaveRepository) SearchMedClave(s string) ([]*claveEntity.ClaveEntity, error) {
	prefixWhere, prefixArgs := utils.BuildPrefixWhere(config.MedPrefixes, config.MedExactKeys, nil)

	query := `
		SELECT 
			id_medicamento, 
			clave_med, 
			descripcion
		FROM medicamentos
		WHERE 
			(clave_med LIKE ? OR descripcion LIKE ?)
			AND ` + prefixWhere + `
		ORDER BY 
			CASE 
				WHEN clave_med = ? THEN 1
				WHEN clave_med LIKE ? THEN 2
				ELSE 3
			END,
			clave_med
		LIMIT 50
	`

	searchTerm := "%" + s + "%"
	args := []interface{}{searchTerm, searchTerm}
	args = append(args, prefixArgs...)
	args = append(args, s, s+"%")

	return r.scanClaves(query, args, false)
}

func (r *ClaveRepository) SearchMedInInventory(s string, id int32) ([]*claveEntity.ClaveEntity, error) {
	prefixWhere, prefixArgs := utils.BuildPrefixWhere(config.MedPrefixes, config.MedExactKeys, nil)

	query := `
		SELECT 
			m.id_medicamento, 
			m.clave_med, 
			m.descripcion,
			d.cantidad
		FROM medicamentos m
		INNER JOIN inventarios i 
			ON i.id_cendis = ?
		INNER JOIN inventario_detalle d 
			ON d.id_inventario = i.id_inventario 
			AND d.id_medicamento = m.id_medicamento
		WHERE 
			(m.clave_med LIKE ? OR m.descripcion LIKE ?)
			AND ` + prefixWhere + `
		ORDER BY 
			CASE 
				WHEN m.clave_med = ? THEN 1
				WHEN m.clave_med LIKE ? THEN 2
				ELSE 3
			END,
			m.clave_med
		LIMIT 50
	`

	searchTerm := "%" + s + "%"
	args := []interface{}{id, searchTerm, searchTerm}
	args = append(args, prefixArgs...)
	args = append(args, s, s+"%")

	return r.scanClaves(query, args, true)
}

// ─── MATERIAL DE CURACION ────────────────────────────────────────────────────

func (r *ClaveRepository) SearchMatClave(s string) ([]*claveEntity.ClaveEntity, error) {
	prefixWhere, prefixArgs := utils.BuildPrefixWhere(config.MatPrefixes, nil, config.MatExclusions)

	query := `
		SELECT 
			id_medicamento, 
			clave_med, 
			descripcion
		FROM medicamentos
		WHERE 
			(clave_med LIKE ? OR descripcion LIKE ?)
			AND ` + prefixWhere + `
		ORDER BY 
			CASE 
				WHEN clave_med = ? THEN 1
				WHEN clave_med LIKE ? THEN 2
				ELSE 3
			END,
			clave_med
		LIMIT 50
	`

	searchTerm := "%" + s + "%"
	args := []interface{}{searchTerm, searchTerm}
	args = append(args, prefixArgs...)
	args = append(args, s, s+"%")

	return r.scanClaves(query, args, false)
}

func (r *ClaveRepository) SearchMatInInventory(s string, id int32) ([]*claveEntity.ClaveEntity, error) {
	prefixWhere, prefixArgs := utils.BuildPrefixWhere(config.MatPrefixes, nil, config.MatExclusions)

	query := `
		SELECT 
			m.id_medicamento, 
			m.clave_med, 
			m.descripcion,
			d.cantidad
		FROM medicamentos m
		INNER JOIN inventarios i 
			ON i.id_cendis = ?
		INNER JOIN inventario_detalle d 
			ON d.id_inventario = i.id_inventario 
			AND d.id_medicamento = m.id_medicamento
		WHERE 
			(m.clave_med LIKE ? OR m.descripcion LIKE ?)
			AND ` + prefixWhere + `
		ORDER BY 
			CASE 
				WHEN m.clave_med = ? THEN 1
				WHEN m.clave_med LIKE ? THEN 2
				ELSE 3
			END,
			m.clave_med
		LIMIT 50
	`

	searchTerm := "%" + s + "%"
	args := []interface{}{id, searchTerm, searchTerm}
	args = append(args, prefixArgs...)
	args = append(args, s, s+"%")

	return r.scanClaves(query, args, true)
}

// ─── HELPER ──────────────────────────────────────────────────────────────────

// scanClaves centraliza el scan de rows para evitar repetición.
// withCantidad indica si el query incluye la columna cantidad (inventory queries).
func (r *ClaveRepository) scanClaves(query string, args []interface{}, withCantidad bool) ([]*claveEntity.ClaveEntity, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error en query: %w", err)
	}
	defer rows.Close()

	claves := []*claveEntity.ClaveEntity{}

	for rows.Next() {
		var c claveEntity.ClaveEntity
		var scanErr error

		if withCantidad {
			scanErr = rows.Scan(&c.Id_medicamento, &c.Clave_med, &c.Descripcion, &c.Cantidad_actual)
		} else {
			scanErr = rows.Scan(&c.Id_medicamento, &c.Clave_med, &c.Descripcion)
		}

		if scanErr != nil {
			return nil, fmt.Errorf("error al escanear resultado: %w", scanErr)
		}
		claves = append(claves, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar resultados: %w", err)
	}

	return claves, nil
}

func (r *ClaveRepository) SearchAllInInventory(s string, id int32) ([]*claveEntity.ClaveEntity, error) {
	query := `
		SELECT 
			m.id_medicamento, 
			m.clave_med, 
			m.descripcion,
			d.cantidad
		FROM medicamentos m
		INNER JOIN inventarios i 
			ON i.id_cendis = ?
		INNER JOIN inventario_detalle d 
			ON d.id_inventario = i.id_inventario 
			AND d.id_medicamento = m.id_medicamento
		WHERE 
			(m.clave_med LIKE ? OR m.descripcion LIKE ?)
		ORDER BY 
			CASE 
				WHEN m.clave_med = ? THEN 1
				WHEN m.clave_med LIKE ? THEN 2
				ELSE 3
			END,
			m.clave_med
		LIMIT 50
	`

	searchTerm := "%" + s + "%"
	args := []interface{}{id, searchTerm, searchTerm, s, s + "%"}

	return r.scanClaves(query, args, true)
}

func (r *ClaveRepository) SearchAllClaves(s string) ([]*claveEntity.ClaveEntity, error) {
	query := `
		SELECT 
			id_medicamento, 
			clave_med, 
			descripcion
		FROM medicamentos
		WHERE 
			clave_med LIKE ? 
			OR descripcion LIKE ?
		ORDER BY 
			CASE 
				WHEN clave_med = ? THEN 1
				WHEN clave_med LIKE ? THEN 2
				ELSE 3
			END,
			clave_med
		LIMIT 50
	`

	searchTerm := "%" + s + "%"
	args := []interface{}{searchTerm, searchTerm, s, s + "%"}

	return r.scanClaves(query, args, false)
}
