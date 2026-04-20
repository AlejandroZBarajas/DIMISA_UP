package salidasInfra

import (
	salidaEntity "DIMISA/src/salidas/salidasDomain/entity"
	"database/sql"
	"fmt"
)

type SalidasRepository struct {
	DB *sql.DB
}

func NewSalidasRepository(db *sql.DB) *SalidasRepository {
	return &SalidasRepository{DB: db}
}

/* func (repo *SalidasRepository) CreateSalida(salida *salidaEntity.SalidaEntity) (int32, error) {
	if len(salida.Claves) == 0 {
		return 0, fmt.Errorf("la salida debe contener al menos un medicamento")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("error al iniciar transacción: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Validar inventario suficiente ANTES de crear la salida
	for _, detalle := range salida.Claves {
		var cantidadActual int32
		queryCheck := `
			SELECT cantidad_actual
			FROM inventarios
			WHERE id_cendis = ? AND id_medicamento = ?
		`
		err = tx.QueryRow(queryCheck, salida.Id_cendis, detalle.Id_medicamento).Scan(&cantidadActual)

		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("medicamento %d no existe en inventario del cendis %d",
				detalle.Id_medicamento, salida.Id_cendis)
		}
		if err != nil {
			return 0, fmt.Errorf("error al verificar inventario: %w", err)
		}

		if cantidadActual < detalle.Cantidad {
			return 0, fmt.Errorf("inventario insuficiente para medicamento %d: disponible=%d, solicitado=%d",
				detalle.Id_medicamento, cantidadActual, detalle.Cantidad)
		}
	}

	// 2. Crear la salida
	querySalida := `
		INSERT INTO salidas (id_area, id_cendis, id_usuario, fecha, editable, pendiente, tipo_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := tx.Exec(
		querySalida,
		salida.Id_area,
		salida.Id_cendis,
		salida.Id_usuario,
		salida.Fecha,
		salida.Editable,
		salida.Pendiente,
		salida.Tipo_id,
	)

	if err != nil {
		return 0, fmt.Errorf("error al crear salida: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener id de salida: %w", err)
	}

	salida.Id_salida = int32(id)

	// 3. Insertar detalles y descontar inventario
	queryDetalle := `
		INSERT INTO salidas_detalle (id_salida, id_medicamento, cantidad)
		VALUES (?, ?, ?)
	`

	queryDescuento := `
		UPDATE inventarios
		SET cantidad_actual = cantidad_actual - ?,
		    updated_at = CURDATE()
		WHERE id_cendis = ? AND id_medicamento = ?
	`

	for _, detalle := range salida.Claves {
		if detalle.Cantidad <= 0 {
			return 0, fmt.Errorf("la cantidad debe ser mayor a 0 para el medicamento %d", detalle.Id_medicamento)
		}

		// Insertar detalle
		_, err = tx.Exec(queryDetalle, salida.Id_salida, detalle.Id_medicamento, detalle.Cantidad)
		if err != nil {
			return 0, fmt.Errorf("error al insertar detalle de salida: %w", err)
		}

		// Descontar del inventario
		_, err = tx.Exec(queryDescuento, detalle.Cantidad, salida.Id_cendis, detalle.Id_medicamento)
		if err != nil {
			return 0, fmt.Errorf("error al descontar inventario: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return salida.Id_salida, nil
} */

func (repo *SalidasRepository) UpdateSalida(id_salida int32, salida *salidaEntity.SalidaEntity) error {

	if len(salida.Claves) == 0 {
		return fmt.Errorf("la salida debe contener al menos un medicamento")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var editable bool
	err = tx.QueryRow("SELECT editable FROM salidas WHERE id_salida = ?", id_salida).Scan(&editable)

	if err != nil {
		return fmt.Errorf("salida no encontrada: %w", err)
	}

	if !editable {
		return fmt.Errorf("la salida no puede ser editada")
	}

	querySalida := `
		UPDATE salidas
		SET fecha = ?
		WHERE id_salida = ?
	`

	result, err := tx.Exec(
		querySalida,
		salida.Fecha,
		id_salida,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar salida: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se pudo actualizar la salida")
	}

	queryDeleteDetalles := `DELETE FROM salidas_detalle WHERE id_salida = ?`
	_, err = tx.Exec(queryDeleteDetalles, id_salida)
	if err != nil {
		return fmt.Errorf("error al eliminar detalles anteriores: %w", err)
	}

	queryDetalle := `
		INSERT INTO salidas_detalle (id_salida, id_medicamento, cantidad)
		VALUES (?, ?, ?)
	`

	for _, detalle := range salida.Claves {
		if detalle.Cantidad <= 0 {
			return fmt.Errorf("la cantidad debe ser mayor a 0 para el medicamento %d", detalle.Id_medicamento)
		}

		_, err = tx.Exec(
			queryDetalle,
			id_salida,
			detalle.Id_medicamento,
			detalle.Cantidad,
		)

		if err != nil {
			return fmt.Errorf("error al insertar detalle de salida: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return nil
}

func (repo *SalidasRepository) GetSalidasByCendis(id_cendis int32) (*[]salidaEntity.SalidaEntity, error) {
	querySalidas := `
		SELECT id_salida, id_area, id_cendis, id_usuario, fecha, created_at, editable, pendiente
		FROM salidas
		WHERE id_cendis = ?
		ORDER BY created_at DESC
	`

	rows, err := repo.DB.Query(querySalidas, id_cendis)
	if err != nil {
		return nil, fmt.Errorf("error al consultar salidas: %w", err)
	}
	defer rows.Close()

	var salidas []salidaEntity.SalidaEntity

	for rows.Next() {
		var salida salidaEntity.SalidaEntity
		err := rows.Scan(
			&salida.Id_salida,
			&salida.Id_area,
			&salida.Id_cendis,
			&salida.Id_usuario,
			&salida.Fecha,
			&salida.Created_at,
			&salida.Editable,
			&salida.Pendiente,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear salida: %w", err)
		}

		queryDetalles := `
			SELECT id_salida_detalle, id_salida, id_medicamento, cantidad
			FROM salidas_detalle
			WHERE id_salida = ?
		`

		detalleRows, err := repo.DB.Query(queryDetalles, salida.Id_salida)
		if err != nil {
			return nil, fmt.Errorf("error al consultar detalles de salida: %w", err)
		}

		var detalles []salidaEntity.SalidaDetalleEntity
		for detalleRows.Next() {
			var detalle salidaEntity.SalidaDetalleEntity
			err := detalleRows.Scan(
				&detalle.Id_salidaDetalle,
				&detalle.Id_salida,
				&detalle.Id_medicamento,
				&detalle.Cantidad,
			)
			if err != nil {
				detalleRows.Close()
				return nil, fmt.Errorf("error al escanear detalle: %w", err)
			}
			detalles = append(detalles, detalle)
		}
		detalleRows.Close()

		salida.Claves = detalles
		salidas = append(salidas, salida)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar salidas: %w", err)
	}

	return &salidas, nil
}

func (repo *SalidasRepository) DeleteSalida(id_salida int32) error {
	query := `DELETE FROM salidas WHERE id_salida = ?`

	result, err := repo.DB.Exec(query, id_salida)
	if err != nil {
		return fmt.Errorf("error al eliminar salida: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar eliminación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("salida no encontrada")
	}

	return nil
}

func (repo *SalidasRepository) GetSalidasPendientes(id_cendis int32) (*[]salidaEntity.SalidaEntity, error) {
	return nil, fmt.Errorf("método no implementado aún")
}

func (repo *SalidasRepository) AddToSalida(id_cendis, id_area, tipo int32, claves *[]salidaEntity.SalidaDetalleEntity) error {
	return fmt.Errorf("método no implementado aún")

}

func (repo *SalidasRepository) CerrarSalida(id_salida int32) error {
	// Verificar que la salida existe y está pendiente
	var pendiente bool
	var editable bool

	err := repo.DB.QueryRow(
		"SELECT pendiente, editable FROM salidas WHERE id_salida = ?",
		id_salida,
	).Scan(&pendiente, &editable)

	if err != nil {
		return fmt.Errorf("salida no encontrada: %w", err)
	}

	// Validar que la salida esté pendiente
	if !pendiente {
		return fmt.Errorf("la salida ya fue cerrada anteriormente")
	}

	// Validar que la salida sea editable
	if !editable {
		return fmt.Errorf("la salida no es editable")
	}

	// Actualizar: editable = false, pendiente = false
	query := `
		UPDATE salidas 
		SET editable = 0, pendiente = 0
		WHERE id_salida = ?
	`

	result, err := repo.DB.Exec(query, id_salida)
	if err != nil {
		return fmt.Errorf("error al cerrar salida: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se pudo cerrar la salida")
	}

	return nil
}

func (repo *SalidasRepository) CreateSalida(salida *salidaEntity.SalidaEntity) (int32, error) {
	if len(salida.Claves) == 0 {
		return 0, fmt.Errorf("la salida debe contener al menos un medicamento")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// 1. Obtener id_inventario
	var idInventario int
	err = tx.QueryRow(`
		SELECT id_inventario 
		FROM inventarios 
		WHERE id_cendis = ?
	`, salida.Id_cendis).Scan(&idInventario)

	if err != nil {
		return 0, fmt.Errorf("inventario no encontrado para cendis %d", salida.Id_cendis)
	}

	// 2. Validación acumulada con lock
	type StockError struct {
		IdMedicamento int32 `json:"id_medicamento"`
		Disponible    int32 `json:"disponible"`
		Solicitado    int32 `json:"solicitado"`
	}

	var errores []StockError

	for _, detalle := range salida.Claves {
		if detalle.Cantidad <= 0 {
			return 0, fmt.Errorf("la cantidad debe ser mayor a 0 para el medicamento %d", detalle.Id_medicamento)
		}

		var disponible int32

		err = tx.QueryRow(`
			SELECT cantidad
			FROM inventario_detalle
			WHERE id_inventario = ? AND id_medicamento = ?
			FOR UPDATE
		`, idInventario, detalle.Id_medicamento).Scan(&disponible)

		if err == sql.ErrNoRows {
			errores = append(errores, StockError{
				IdMedicamento: detalle.Id_medicamento,
				Disponible:    0,
				Solicitado:    detalle.Cantidad,
			})
			continue
		}

		if err != nil {
			return 0, fmt.Errorf("error al verificar inventario: %w", err)
		}

		if disponible < detalle.Cantidad {
			errores = append(errores, StockError{
				IdMedicamento: detalle.Id_medicamento,
				Disponible:    disponible,
				Solicitado:    detalle.Cantidad,
			})
		}
	}

	// 3. Abortamos si hay errores
	if len(errores) > 0 {
		return 0, fmt.Errorf("stock insuficiente: %+v", errores)
	}

	// 4. Crear salida
	querySalida := `
		INSERT INTO salidas (id_area, id_cendis, id_usuario, fecha, editable, pendiente, tipo_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := tx.Exec(
		querySalida,
		salida.Id_area,
		salida.Id_cendis,
		salida.Id_usuario,
		salida.Fecha,
		salida.Editable,
		salida.Pendiente,
		salida.Tipo_id,
	)
	if err != nil {
		return 0, fmt.Errorf("error al crear salida: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener id de salida: %w", err)
	}

	salida.Id_salida = int32(id)

	// 5. Insertar detalles
	queryDetalle := `
		INSERT INTO salidas_detalle (id_salida, id_medicamento, cantidad)
		VALUES (?, ?, ?)
	`

	for _, detalle := range salida.Claves {
		_, err = tx.Exec(queryDetalle, salida.Id_salida, detalle.Id_medicamento, detalle.Cantidad)
		if err != nil {
			return 0, fmt.Errorf("error al insertar detalle de salida: %w", err)
		}
	}

	// 6. Descontar inventario con protección
	queryDescuento := `
		UPDATE inventario_detalle
		SET 
			cantidad = cantidad - ?,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = ?
		WHERE 
			id_inventario = ? 
			AND id_medicamento = ?
			AND cantidad >= ?
	`

	for _, detalle := range salida.Claves {
		res, err := tx.Exec(
			queryDescuento,
			detalle.Cantidad,
			salida.Id_usuario,
			idInventario,
			detalle.Id_medicamento,
			detalle.Cantidad,
		)
		if err != nil {
			return 0, fmt.Errorf("error al descontar inventario: %w", err)
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			return 0, fmt.Errorf(
				"conflicto concurrente al descontar medicamento %d",
				detalle.Id_medicamento,
			)
		}
	}

	// 7. Commit
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return salida.Id_salida, nil
}
