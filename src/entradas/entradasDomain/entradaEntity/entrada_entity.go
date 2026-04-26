package entradaEntity

type DetalleEntrada struct {
	Id_medicamento int32 `json:"id_medicamento"`
	Cantidad       int32 `json:"cantidad"`
}

type EntradaRequest struct {
	Id_cendis    int32            `json:"id_cendis"`
	Id_usuario   int32            `json:"id_usuario"`
	Id_colectivo int32            `json:"id_colectivo"`
	Detalles     []DetalleEntrada `json:"detalles"`
}

type InventarioRequest struct {
	Id_cendis  int32            `json:"id_cendis"`
	Id_usuario int32            `json:"id_usuario"`
	Detalles   []DetalleEntrada `json:"detalles"`
}
