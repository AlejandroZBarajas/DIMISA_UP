package inventariosInfra

import (
	inventarioEntity "DIMISA/src/inventarioODT/inventariosDomain/inventariosEntity"
	"database/sql"
)

type InventariosRepository struct {
	DB *sql.DB
}

func NewInventariosRepository(db *sql.DB) *InventariosRepository {
	return &InventariosRepository{DB: db}
}

func (r *InventariosRepository) GetAllInventarios() ([]inventarioEntity.InventarioEntity, error) {
	query := `
		SELECT 
			i.id_inventario,
			i.id_cendis,
			c.cendis_nombre,
			d.id_medicamento,
			m.clave_med,
			m.descripcion,
			d.cantidad
		FROM inventarios i
		INNER JOIN cendis c ON c.id_cendis = i.id_cendis
		LEFT JOIN inventario_detalle d ON d.id_inventario = i.id_inventario
		LEFT JOIN medicamentos m ON m.id_medicamento = d.id_medicamento
		ORDER BY i.id_inventario
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanInventarios(rows)
}

func (r *InventariosRepository) GetInventarioByCendisID(id_cendis int32) (inventarioEntity.InventarioEntity, error) {
	query := `
		SELECT 
			i.id_inventario,
			i.id_cendis,
			c.cendis_nombre,
			d.id_medicamento,
			m.clave_med,
			m.descripcion,
			d.cantidad
		FROM inventarios i
		INNER JOIN cendis c ON c.id_cendis = i.id_cendis
		LEFT JOIN inventario_detalle d ON d.id_inventario = i.id_inventario
		LEFT JOIN medicamentos m ON m.id_medicamento = d.id_medicamento
		WHERE i.id_cendis = ?
	`

	rows, err := r.DB.Query(query, id_cendis)
	if err != nil {
		return inventarioEntity.InventarioEntity{}, err
	}
	defer rows.Close()

	inventarios, err := scanInventarios(rows)
	if err != nil {
		return inventarioEntity.InventarioEntity{}, err
	}
	if len(inventarios) == 0 {
		return inventarioEntity.InventarioEntity{}, sql.ErrNoRows
	}

	return inventarios[0], nil
}

// scanInventarios agrupa filas por inventario y construye los detalles
func scanInventarios(rows *sql.Rows) ([]inventarioEntity.InventarioEntity, error) {
	inventariosMap := make(map[int32]*inventarioEntity.InventarioEntity)
	var order []int32

	for rows.Next() {
		var (
			idInventario  int32
			idCendis      int32
			cendisNombre  string
			idMedicamento sql.NullInt32
			claveMed      sql.NullString
			descripcion   sql.NullString
			cantidad      sql.NullInt32
		)

		if err := rows.Scan(
			&idInventario,
			&idCendis,
			&cendisNombre,
			&idMedicamento,
			&claveMed,
			&descripcion,
			&cantidad,
		); err != nil {
			return nil, err
		}

		if _, exists := inventariosMap[idInventario]; !exists {
			inventariosMap[idInventario] = &inventarioEntity.InventarioEntity{
				Id:        idInventario,
				Id_cendis: idCendis,
				Cendis:    cendisNombre,
				Detalles:  []inventarioEntity.DetalleInventario{},
			}
			order = append(order, idInventario)
		}

		if idMedicamento.Valid {
			detalle := inventarioEntity.DetalleInventario{
				Id_medicamento: idMedicamento.Int32,
				Clave:          claveMed.String,
				Descripcion:    descripcion.String,
				Cantidad:       cantidad.Int32,
			}
			inventariosMap[idInventario].Detalles = append(
				inventariosMap[idInventario].Detalles,
				detalle,
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]inventarioEntity.InventarioEntity, 0, len(order))
	for _, id := range order {
		result = append(result, *inventariosMap[id])
	}

	return result, nil
}
