package entradaColectivoInfra

import (
	entity "DIMISA/src/entradaColectivo/entradaColectivoDomain/entradaColectivoEntity"
	"database/sql"
	"fmt"
	"log"
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
        ec.id_cendis,
        ec.id_medicamento,
        m.clave_med,
        m.descripcion,
        SUM(cd.piezas_esperadas)  AS piezas_esperadas,
        SUM(ec.cantidad_recibida) AS cantidad_recibida,
        MONTH(MIN(ec.created_at)),
        YEAR(MIN(ec.created_at))
    FROM entradas_colectivo ec
    INNER JOIN medicamentos m ON m.id_medicamento = ec.id_medicamento
    INNER JOIN colectivo_detalle cd ON cd.id_colectivo   = ec.id_colectivo
                                   AND cd.id_medicamento  = ec.id_medicamento
    INNER JOIN cendis c ON c.id_cendis = ec.id_cendis
    WHERE ec.id_cendis         = ?
    AND   MONTH(ec.created_at) = ?
    AND   YEAR(ec.created_at)  = ?
    GROUP BY ec.id_cendis, ec.id_medicamento, m.clave_med, m.descripcion
    ORDER BY ec.id_medicamento
	`
	log.Printf("[GetReporteMensual] START id_cendis=%d mes=%d anio=%d", idCendis, mes, anio)

	reporte := entity.EntradaColectivoReporte{
		Mes:      mes,
		Anio:     anio,
		Detalles: []entity.EntradaColectivoDetalle{},
	}

	err := r.DB.QueryRow(`
		SELECT cendis_nombre FROM cendis WHERE id_cendis = ?
	`, idCendis).Scan(&reporte.Cendis)
	if err != nil {
		return entity.EntradaColectivoReporte{}, fmt.Errorf("SELECT cendis_nombre: %w", err)
	}

	rows, err := r.DB.Query(query, idCendis, mes, anio)
	if err != nil {
		log.Printf("[GetReporteMensual] ERROR Query: %v", err)
		return entity.EntradaColectivoReporte{}, err
	}
	defer rows.Close()

	reporte = entity.EntradaColectivoReporte{
		Mes:      mes,
		Anio:     anio,
		Detalles: []entity.EntradaColectivoDetalle{},
	}

	var totalSolicitado, totalRecibido int32

	for rows.Next() {
		var d entity.EntradaColectivoDetalle
		if err := rows.Scan(
			&d.IdCendis,
			&d.IdMedicamento,
			&d.Clave,
			&d.Descripcion,
			&d.PiezasEsperadas,
			&d.CantidadRecibida,
			&d.Mes,
			&d.Anio,
		); err != nil {
			log.Printf("[GetReporteMensual] ERROR Scan: %v", err)
			return entity.EntradaColectivoReporte{}, err
		}
		log.Printf("[GetReporteMensual] detalle → id_medicamento=%d piezas_esperadas=%d recibida=%d", d.IdMedicamento, d.PiezasEsperadas, d.CantidadRecibida)

		d.Deficit = d.PiezasEsperadas - d.CantidadRecibida
		d.Estatus = calcularEstatus(d.PiezasEsperadas, d.CantidadRecibida)

		totalSolicitado += d.PiezasEsperadas
		totalRecibido += d.CantidadRecibida

		reporte.Detalles = append(reporte.Detalles, d)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[GetReporteMensual] ERROR rows.Err: %v", err)
		return entity.EntradaColectivoReporte{}, err
	}

	reporte.IdCendis = idCendis
	reporte.TotalSolicitado = totalSolicitado
	reporte.TotalRecibido = totalRecibido
	reporte.TotalDeficit = totalSolicitado - totalRecibido
	reporte.PctCumplimiento = calcularPct(totalSolicitado, totalRecibido)

	log.Printf("[GetReporteMensual] DONE detalles=%d", len(reporte.Detalles))
	return reporte, nil
}
func (r *EntradaColectivoRepository) GetDeficitCronico(idCendis int32, anio int) ([]entity.EntradaColectivoDetalle, error) {
	query := `
    SELECT
        ec.id_medicamento,
        m.clave_med,
        m.descripcion,
        ec.id_cendis,
        SUM(cd.piezas_esperadas)  AS total_solicitado,
        SUM(ec.cantidad_recibida) AS total_recibido
    FROM entradas_colectivo ec
    INNER JOIN medicamentos m ON m.id_medicamento = ec.id_medicamento
    INNER JOIN colectivo_detalle cd ON cd.id_colectivo   = ec.id_colectivo
                                   AND cd.id_medicamento  = ec.id_medicamento
    WHERE ec.id_cendis        = ?
    AND   YEAR(ec.created_at) = ?
    AND   cd.piezas_esperadas > 0
    GROUP BY ec.id_medicamento, m.clave_med, m.descripcion, ec.id_cendis
    HAVING SUM(cd.piezas_esperadas) > SUM(ec.cantidad_recibida)
    ORDER BY (SUM(cd.piezas_esperadas) - SUM(ec.cantidad_recibida)) DESC
`
	log.Printf("[GetDeficitCronico] START id_cendis=%d anio=%d", idCendis, anio)

	rows, err := r.DB.Query(query, idCendis, anio)
	if err != nil {
		log.Printf("[GetDeficitCronico] ERROR Query: %v", err)
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
			&d.PiezasEsperadas,
			&d.CantidadRecibida,
		); err != nil {
			log.Printf("[GetDeficitCronico] ERROR Scan: %v", err)
			return nil, err
		}

		d.Deficit = d.PiezasEsperadas - d.CantidadRecibida
		d.Estatus = calcularEstatus(d.PiezasEsperadas, d.CantidadRecibida)

		detalles = append(detalles, d)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[GetDeficitCronico] ERROR rows.Err: %v", err)
		return nil, err
	}

	log.Printf("[GetDeficitCronico] DONE detalles=%d", len(detalles))
	return detalles, nil
}
func (r *EntradaColectivoRepository) GetComparativoCendis(mes int, anio int) ([]entity.EntradaColectivoReporte, error) {
	query := `
    SELECT
        ec.id_cendis,
        c.cendis_nombre,
        SUM(cd.piezas_esperadas)  AS total_solicitado,
        SUM(ec.cantidad_recibida) AS total_recibido
    FROM entradas_colectivo ec
    INNER JOIN cendis c ON c.id_cendis = ec.id_cendis
    INNER JOIN colectivo_detalle cd ON cd.id_colectivo   = ec.id_colectivo
                                   AND cd.id_medicamento  = ec.id_medicamento
    WHERE MONTH(ec.created_at) = ?
    AND   YEAR(ec.created_at)  = ?
    AND   cd.piezas_esperadas  > 0
    GROUP BY ec.id_cendis, c.cendis_nombre
    HAVING SUM(cd.piezas_esperadas) > 0
    ORDER BY (SUM(cd.piezas_esperadas) - SUM(ec.cantidad_recibida)) DESC
`
	log.Printf("[GetComparativoCendis] START mes=%d anio=%d", mes, anio)

	rows, err := r.DB.Query(query, mes, anio)
	if err != nil {
		log.Printf("[GetComparativoCendis] ERROR Query: %v", err)
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
			log.Printf("[GetComparativoCendis] ERROR Scan: %v", err)
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
		log.Printf("[GetComparativoCendis] ERROR rows.Err: %v", err)
		return nil, err
	}

	log.Printf("[GetComparativoCendis] DONE reportes=%d", len(reportes))
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
	return math.Round(pct*100) / 100 // redondear a 2 decimales
}
