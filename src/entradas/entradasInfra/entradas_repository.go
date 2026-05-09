package entradasInfra

import (
	"DIMISA/src/entradas/entradasDomain/entradaEntity"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type EntradasRepository struct {
	DB *sql.DB
}

func NewEntradasRepository(db *sql.DB) *EntradasRepository {
	return &EntradasRepository{DB: db}
}

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

/*
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
*/

func (r *EntradasRepository) CapturarEntrada(entrada *entradaEntity.EntradaRequest) error {
	log.Printf("[CapturarEntrada] START id_colectivo=%d id_cendis=%d id_usuario=%d detalles=%d",
		entrada.Id_colectivo, entrada.Id_cendis, entrada.Id_usuario, len(entrada.Detalles))

	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("[CapturarEntrada] ERROR Begin tx: %v", err)
		return err
	}
	defer tx.Rollback()

	if len(entrada.Detalles) > 0 {
		var idInventario int

		_, err = tx.Exec(`
			INSERT IGNORE INTO inventarios (id_cendis, updated_at)
			VALUES (?, NOW())
		`, entrada.Id_cendis)
		if err != nil {
			log.Printf("[CapturarEntrada] ERROR INSERT IGNORE inventarios: %v", err)
			return err
		}

		err = tx.QueryRow(`
			SELECT id_inventario FROM inventarios WHERE id_cendis = ?
		`, entrada.Id_cendis).Scan(&idInventario)
		if err != nil {
			log.Printf("[CapturarEntrada] ERROR SELECT id_inventario: %v", err)
			return err
		}
		log.Printf("[CapturarEntrada] id_inventario=%d", idInventario)

		// ── inventario_detalle ────────────────────────────────────────
		invPlaceholders := make([]string, 0, len(entrada.Detalles))
		invArgs := make([]interface{}, 0, len(entrada.Detalles)*4)

		for _, d := range entrada.Detalles {
			log.Printf("[CapturarEntrada] detalle → id_medicamento=%d cantidad=%d", d.Id_medicamento, d.Cantidad)
			invPlaceholders = append(invPlaceholders, "(?, ?, ?, ?, NOW())")
			invArgs = append(invArgs, idInventario, d.Id_medicamento, d.Cantidad, entrada.Id_usuario)
		}

		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO inventario_detalle 
				(id_inventario, id_medicamento, cantidad, updated_by, updated_at)
			VALUES %s
			ON DUPLICATE KEY UPDATE
				cantidad   = cantidad + VALUES(cantidad),
				updated_by = VALUES(updated_by),
				updated_at = NOW()
		`, strings.Join(invPlaceholders, ", ")), invArgs...)
		if err != nil {
			log.Printf("[CapturarEntrada] ERROR bulk insert inventario_detalle: %v", err)
			return err
		}
		log.Printf("[CapturarEntrada] inventario_detalle OK")

		_, err = tx.Exec(`
			DELETE FROM inventario_detalle
			WHERE id_inventario = ? AND cantidad <= 0
		`, idInventario)
		if err != nil {
			log.Printf("[CapturarEntrada] ERROR DELETE cantidad<=0: %v", err)
			return err
		}

		_, err = tx.Exec(`
			UPDATE inventarios SET updated_at = NOW() WHERE id_inventario = ?
		`, idInventario)
		if err != nil {
			log.Printf("[CapturarEntrada] ERROR UPDATE inventarios: %v", err)
			return err
		}
		log.Printf("[CapturarEntrada] inventarios updated_at OK")

		// ── entradas_colectivo ────────────────────────────────────────
		recibido := make(map[int32]int32, len(entrada.Detalles))
		for _, d := range entrada.Detalles {
			recibido[d.Id_medicamento] += d.Cantidad
		}
		log.Printf("[CapturarEntrada] mapa recibido: %v", recibido)

		rows, err := tx.Query(`
			SELECT id_medicamento, cantidad
			FROM colectivo_detalle
			WHERE id_colectivo = ?
		`, entrada.Id_colectivo)
		if err != nil {
			log.Printf("[CapturarEntrada] ERROR SELECT colectivo_detalle: %v", err)
			return err
		}

		ecPlaceholders := make([]string, 0)
		ecArgs := make([]interface{}, 0)

		for rows.Next() {
			var idMed int32
			var solicitada int32
			if err := rows.Scan(&idMed, &solicitada); err != nil {
				log.Printf("[CapturarEntrada] ERROR Scan colectivo_detalle: %v", err)
				rows.Close()
				return err
			}
			cantRecibida := recibido[idMed]
			log.Printf("[CapturarEntrada] colectivo_detalle → id_medicamento=%d solicitada=%d recibida=%d",
				idMed, solicitada, cantRecibida)

			ecPlaceholders = append(ecPlaceholders, "(?, ?, ?, ?, ?)")
			ecArgs = append(ecArgs,
				entrada.Id_colectivo,
				entrada.Id_cendis,
				idMed,
				solicitada,
				cantRecibida,
			)
		}

		rows.Close()
		if err := rows.Err(); err != nil {
			log.Printf("[CapturarEntrada] ERROR rows.Err colectivo_detalle: %v", err)
			return err
		}
		log.Printf("[CapturarEntrada] colectivo_detalle filas procesadas=%d", len(ecPlaceholders))

		if len(ecPlaceholders) > 0 {
			_, err = tx.Exec(fmt.Sprintf(`
				INSERT INTO entradas_colectivo
					(id_colectivo, id_cendis, id_medicamento, cantidad_solicitada, cantidad_recibida)
				VALUES %s
				ON DUPLICATE KEY UPDATE
					cantidad_recibida = VALUES(cantidad_recibida)
			`, strings.Join(ecPlaceholders, ", ")), ecArgs...)
			if err != nil {
				log.Printf("[CapturarEntrada] ERROR bulk insert entradas_colectivo: %v", err)
				return err
			}
			log.Printf("[CapturarEntrada] entradas_colectivo OK")
		} else {
			log.Printf("[CapturarEntrada] WARN colectivo_detalle vacío, no se insertó en entradas_colectivo")
		}
	}

	_, err = tx.Exec(`
		UPDATE colectivos SET capturado = 1 WHERE id_colectivo = ?
	`, entrada.Id_colectivo)
	if err != nil {
		log.Printf("[CapturarEntrada] ERROR UPDATE colectivos: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[CapturarEntrada] ERROR Commit: %v", err)
		return err
	}

	log.Printf("[CapturarEntrada] DONE id_colectivo=%d", entrada.Id_colectivo)
	return nil
}
