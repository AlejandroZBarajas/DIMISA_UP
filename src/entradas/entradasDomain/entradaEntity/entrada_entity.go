package entradaEntity

type DetalleEntrada struct {
	Id_medicamento int32 `json:"id_medicamento"`
	Cantidad       int32 `json:"cantidad"`
}

type EntradaRequest struct {
	Id_cendis    int32            `json:"id_cendis"`
	Id_colectivo int32            `json:"id_colectivo"`
	Detalles     []DetalleEntrada `json:"detalles"`
}

type InventarioRequest struct {
	Id_cendis int32            `json:"id_cendis"`
	Detalles  []DetalleEntrada `json:"detalles"`
}
