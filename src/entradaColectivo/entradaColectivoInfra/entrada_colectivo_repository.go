package entradaColectivoInfra

import (
	entity "DIMISA/src/entradaColectivo/entradaColectivoDomain/entradaColectivoEntity"
	"database/sql"
	"math"
)

type EntradaColectivoRepository struct {
	DB *sql.DB
}

func NewEntradaColectivoRepository(db *sql.DB) *EntradaColectivoRepository {
	return &EntradaColectivoRepository{DB: db}
}

func (r *EntradaColectivoRepository) GetReporteMensual(idCendis int32, mes int, anio int) (entity.EntradaColectivoReporte, error) {
	query := `
		SELECT
			ec.id_entrada_colectivo,
			ec.id_colectivo,
			ec.id_cendis,
			c.cendis_nombre,
			ec.id_medicamento,
			m.clave_med,
			m.descripcion,
			ec.cantidad_solicitada,
			ec.cantidad_recibida,
			MONTH(ec.created_at),
			YEAR(ec.created_at)
		FROM entradas_colectivo ec
		INNER JOIN cendis      c ON c.id_cendis      = ec.id_cendis
		INNER JOIN medicamentos m ON m.id_medicamento = ec.id_medicamento
		WHERE ec.id_cendis        = ?
		AND   MONTH(ec.created_at) = ?
		AND   YEAR(ec.created_at)  = ?
		ORDER BY ec.id_medicamento
	`

	rows, err := r.DB.Query(query, idCendis, mes, anio)
	if err != nil {
		return entity.EntradaColectivoReporte{}, err
	}
	defer rows.Close()

	reporte := entity.EntradaColectivoReporte{
		Mes:      mes,
		Anio:     anio,
		Detalles: []entity.EntradaColectivoDetalle{},
	}

	var totalSolicitado, totalRecibido int32

	for rows.Next() {
		var d entity.EntradaColectivoDetalle
		if err := rows.Scan(
			&d.IdEntradaColectivo,
			&d.IdColectivo,
			&d.IdCendis,
			&reporte.Cendis,
			&d.IdMedicamento,
			&d.Clave,
			&d.Descripcion,
			&d.CantidadSolicitada,
			&d.CantidadRecibida,
			&d.Mes,
			&d.Anio,
		); err != nil {
			return entity.EntradaColectivoReporte{}, err
		}

		d.Deficit = d.CantidadSolicitada - d.CantidadRecibida
		d.Estatus = calcularEstatus(d.CantidadSolicitada, d.CantidadRecibida)

		totalSolicitado += d.CantidadSolicitada
		totalRecibido += d.CantidadRecibida

		reporte.Detalles = append(reporte.Detalles, d)
	}
	if err := rows.Err(); err != nil {
		return entity.EntradaColectivoReporte{}, err
	}

	reporte.IdCendis = idCendis
	reporte.TotalSolicitado = totalSolicitado
	reporte.TotalRecibido = totalRecibido
	reporte.TotalDeficit = totalSolicitado - totalRecibido
	reporte.PctCumplimiento = calcularPct(totalSolicitado, totalRecibido)

	return reporte, nil
}

func (r *EntradaColectivoRepository) GetDeficitCronico(idCendis int32, anio int) ([]entity.EntradaColectivoDetalle, error) {
	query := `
		SELECT
			ec.id_medicamento,
			m.clave_med,
			m.descripcion,
			ec.id_cendis,
			SUM(ec.cantidad_solicitada) AS total_solicitado,
			SUM(ec.cantidad_recibida)   AS total_recibido
		FROM entradas_colectivo ec
		INNER JOIN medicamentos m ON m.id_medicamento = ec.id_medicamento
		WHERE ec.id_cendis        = ?
		AND   YEAR(ec.created_at) = ?
		GROUP BY ec.id_medicamento, m.clave_med, m.descripcion, ec.id_cendis
		HAVING SUM(ec.cantidad_solicitada) > SUM(ec.cantidad_recibida)
		ORDER BY SUM(ec.cantidad_solicitada - ec.cantidad_recibida) DESC
	`

	rows, err := r.DB.Query(query, idCendis, anio)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var detalles []entity.EntradaColectivoDetalle

	for rows.Next() {
		var d entity.EntradaColectivoDetalle
		if err := rows.Scan(
			&d.IdMedicamento,
			&d.Clave,
			&d.Descripcion,
			&d.IdCendis,
			&d.CantidadSolicitada,
			&d.CantidadRecibida,
		); err != nil {
			return nil, err
		}

		d.Deficit = d.CantidadSolicitada - d.CantidadRecibida
		d.Estatus = calcularEstatus(d.CantidadSolicitada, d.CantidadRecibida)

		detalles = append(detalles, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return detalles, nil
}

func (r *EntradaColectivoRepository) GetComparativoCendis(mes int, anio int) ([]entity.EntradaColectivoReporte, error) {
	query := `
		SELECT
			ec.id_cendis,
			c.cendis_nombre,
			SUM(ec.cantidad_solicitada) AS total_solicitado,
			SUM(ec.cantidad_recibida)   AS total_recibido
		FROM entradas_colectivo ec
		INNER JOIN cendis c ON c.id_cendis = ec.id_cendis
		WHERE MONTH(ec.created_at) = ?
		AND   YEAR(ec.created_at)  = ?
		GROUP BY ec.id_cendis, c.cendis_nombre
		ORDER BY SUM(ec.cantidad_solicitada - ec.cantidad_recibida) DESC
	`

	rows, err := r.DB.Query(query, mes, anio)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportes []entity.EntradaColectivoReporte

	for rows.Next() {
		var rep entity.EntradaColectivoReporte
		if err := rows.Scan(
			&rep.IdCendis,
			&rep.Cendis,
			&rep.TotalSolicitado,
			&rep.TotalRecibido,
		); err != nil {
			return nil, err
		}

		rep.Mes = mes
		rep.Anio = anio
		rep.TotalDeficit = rep.TotalSolicitado - rep.TotalRecibido
		rep.PctCumplimiento = calcularPct(rep.TotalSolicitado, rep.TotalRecibido)
		rep.Detalles = []entity.EntradaColectivoDetalle{}

		reportes = append(reportes, rep)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reportes, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func calcularEstatus(solicitada, recibida int32) string {
	switch {
	case recibida >= solicitada:
		return "Completo"
	case recibida == 0:
		return "No surtido"
	default:
		return "Parcial"
	}
}

func calcularPct(solicitada, recibida int32) float64 {
	if solicitada == 0 {
		return 0
	}
	pct := float64(recibida) / float64(solicitada) * 100
	return math.Round(pct*100) / 100
}
