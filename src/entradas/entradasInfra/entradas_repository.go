package entradasInfra

import (
	"DIMISA/src/entradas/entradasDomain/entradaEntity"
	"database/sql"
	"fmt"
	"strings"
)

type EntradasRepository struct {
	DB *sql.DB
}

func NewEntradasRepository(db *sql.DB) *EntradasRepository {
	return &EntradasRepository{DB: db}
}

/* func (r *EntradasRepository) CapturarEntrada(entrada *entradaEntity.EntradaRequest) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO inventarios (id_cendis, id_medicamento, cantidad_actual, updated_at)
        VALUES (?, ?, ?, NOW())
        ON DUPLICATE KEY UPDATE
            cantidad_actual = cantidad_actual + VALUES(cantidad_actual),
            updated_at = NOW()
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, detalle := range entrada.Detalles {
		_, err := stmt.Exec(entrada.Id_cendis, detalle.Id_medicamento, detalle.Cantidad)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("UPDATE colectivos SET capturado = 1 WHERE id_colectivo = ?", entrada.Id_colectivo)
	if err != nil {
		return err
	}

	return tx.Commit()
}
*/
/*
func (r *EntradasRepository) CapturarInventario(inventario *entradaEntity.InventarioRequest) error {
	if len(inventario.Detalles) == 0 {
		return fmt.Errorf("detalles vacíos")
	}

	// Bulk insert con ON DUPLICATE KEY UPDATE (SET como carga inicial, no suma)
	placeholders := make([]string, 0, len(inventario.Detalles))
	args := make([]interface{}, 0, len(inventario.Detalles)*3)

	for _, d := range inventario.Detalles {
		placeholders = append(placeholders, "(?, ?, ?, NOW())")
		args = append(args, inventario.Id_cendis, d.Id_medicamento, d.Cantidad)
	}

	query := fmt.Sprintf(`
        INSERT INTO inventarios (id_cendis, id_medicamento, cantidad_actual, updated_at)
        VALUES %s
        ON DUPLICATE KEY UPDATE
            cantidad_actual = cantidad_actual + VALUES(cantidad_actual),
            updated_at      = NOW()
    `, strings.Join(placeholders, ", "))

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al capturar inventario: %w", err)
	}

	return nil
}
*/

func (r *EntradasRepository) CapturarInventario(inventario *entradaEntity.InventarioRequest) error {
	if len(inventario.Detalles) == 0 {
		return fmt.Errorf("detalles vacíos")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var idInventario int

	// 1. Obtener o crear inventario
	err = tx.QueryRow(`
		SELECT id_inventario 
		FROM inventarios 
		WHERE id_cendis = ?
	`, inventario.Id_cendis).Scan(&idInventario)

	if err != nil {
		if err == sql.ErrNoRows {
			// crear inventario
			res, err := tx.Exec(`
				INSERT INTO inventarios (id_cendis, updated_at)
				VALUES (?, NOW())
			`, inventario.Id_cendis)
			if err != nil {
				return fmt.Errorf("error creando inventario: %w", err)
			}

			lastID, err := res.LastInsertId()
			if err != nil {
				return err
			}
			idInventario = int(lastID)

		} else {
			return err
		}
	}

	// 2. Bulk insert detalles
	placeholders := make([]string, 0, len(inventario.Detalles))
	args := make([]interface{}, 0, len(inventario.Detalles)*4)

	for _, d := range inventario.Detalles {
		placeholders = append(placeholders, "(?, ?, ?, ?, NOW())")
		args = append(args,
			idInventario,
			d.Id_medicamento,
			d.Cantidad,
			inventario.Id_usuario,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO inventario_detalle 
		(id_inventario, id_medicamento, cantidad, updated_by, updated_at)
		VALUES %s
		ON DUPLICATE KEY UPDATE
			cantidad = cantidad + VALUES(cantidad),
			updated_by = VALUES(updated_by),
			updated_at = NOW()
	`, strings.Join(placeholders, ", "))

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error insertando detalles: %w", err)
	}

	// 3. actualizar timestamp del inventario
	_, err = tx.Exec(`
		UPDATE inventarios 
		SET updated_at = NOW()
		WHERE id_inventario = ?
	`, idInventario)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *EntradasRepository) CapturarEntrada(entrada *entradaEntity.EntradaRequest) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Si hay detalles, procesarlos — si no, solo marcar el colectivo
	if len(entrada.Detalles) > 0 {
		var idInventario int

		err = tx.QueryRow(`
			SELECT id_inventario 
			FROM inventarios 
			WHERE id_cendis = ?
		`, entrada.Id_cendis).Scan(&idInventario)

		if err != nil {
			if err == sql.ErrNoRows {
				res, err := tx.Exec(`
					INSERT INTO inventarios (id_cendis, updated_at)
					VALUES (?, NOW())
				`, entrada.Id_cendis)
				if err != nil {
					return err
				}
				lastID, err := res.LastInsertId()
				if err != nil {
					return err
				}
				idInventario = int(lastID)
			} else {
				return err
			}
		}

		placeholders := make([]string, 0, len(entrada.Detalles))
		args := make([]interface{}, 0, len(entrada.Detalles)*4)

		for _, d := range entrada.Detalles {
			placeholders = append(placeholders, "(?, ?, ?, ?, NOW())")
			args = append(args, idInventario, d.Id_medicamento, d.Cantidad, entrada.Id_usuario)
		}

		query := fmt.Sprintf(`
			INSERT INTO inventario_detalle 
			(id_inventario, id_medicamento, cantidad, updated_by, updated_at)
			VALUES %s
			ON DUPLICATE KEY UPDATE
				cantidad = cantidad + VALUES(cantidad),
				updated_by = VALUES(updated_by),
				updated_at = NOW()
		`, strings.Join(placeholders, ", "))

		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE inventarios 
			SET updated_at = NOW()
			WHERE id_inventario = ?
		`, idInventario)
		if err != nil {
			return err
		}
	}

	// Siempre marcar el colectivo como capturado
	_, err = tx.Exec(`
		UPDATE colectivos 
		SET capturado = 1 
		WHERE id_colectivo = ?
	`, entrada.Id_colectivo)
	if err != nil {
		return err
	}

	return tx.Commit()
}
