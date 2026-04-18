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

func (r *EntradasRepository) CapturarEntrada(entrada *entradaEntity.EntradaRequest) error {
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
