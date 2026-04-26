package inventarioEntity

type DetalleInventario struct {
	Id_medicamento int32  `json:"id_medicamento"`
	Clave          string `json:"clave"`
	Descripcion    string `json:"descripcion"`
	Cantidad       int32  `json:"cantidad"`
}

type InventarioEntity struct {
	Id        int32               `json:"id"`
	Id_cendis int32               `json:"id_cendis"`
	Cendis    string              `json:"cendis"`
	Detalles  []DetalleInventario `json:"detalles"`
}
