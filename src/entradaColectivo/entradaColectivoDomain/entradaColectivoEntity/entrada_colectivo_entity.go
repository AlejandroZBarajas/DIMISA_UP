package entradaColectivoEntity

type EntradaColectivoDetalle struct {
	IdEntradaColectivo int32  `json:"id_entrada_colectivo"`
	IdColectivo        int32  `json:"id_colectivo"`
	IdCendis           int32  `json:"id_cendis"`
	IdMedicamento      int32  `json:"id_medicamento"`
	Clave              string `json:"clave"`
	Descripcion        string `json:"descripcion"`
	CantidadSolicitada int32  `json:"cantidad_solicitada"`
	CantidadRecibida   int32  `json:"cantidad_recibida"`
	Deficit            int32  `json:"deficit"`
	Estatus            string `json:"estatus"` // "Completo" | "Parcial" | "No surtido"
	Mes                int    `json:"mes"`
	Anio               int    `json:"anio"`
}

type EntradaColectivoReporte struct {
	IdCendis        int32                     `json:"id_cendis"`
	Cendis          string                    `json:"cendis"`
	Mes             int                       `json:"mes"`
	Anio            int                       `json:"anio"`
	TotalSolicitado int32                     `json:"total_solicitado"`
	TotalRecibido   int32                     `json:"total_recibido"`
	TotalDeficit    int32                     `json:"total_deficit"`
	PctCumplimiento float64                   `json:"pct_cumplimiento"`
	Detalles        []EntradaColectivoDetalle `json:"detalles"`
}
