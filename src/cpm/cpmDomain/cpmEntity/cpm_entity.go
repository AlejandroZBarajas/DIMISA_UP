package cpmEntity

type DetalleMes struct {
	Mes     int    `json:"mes"`
	Anio    int    `json:"anio"`
	Nombre  string `json:"nombre"`
	Consumo int32  `json:"consumo"`
}

type CpmDetalle struct {
	IdMedicamento   int32        `json:"id_medicamento"`
	Clave           string       `json:"clave"`
	Descripcion     string       `json:"descripcion"`
	Meses           []DetalleMes `json:"meses"`
	Sumatoria       int32        `json:"sumatoria"`
	PromedioMensual float64      `json:"promedio_mensual"`
	PromedioDiario  float64      `json:"promedio_diario"`
	Diez            float64      `json:"diez_pct"`
	ConsumoDiario   float64      `json:"consumo_diario"`
	ConsumoMensual  float64      `json:"consumo_mensual"`
}

type CpmEntity struct {
	Medicamentos []CpmDetalle `json:"medicamentos"`
	Material     []CpmDetalle `json:"material"`
}
