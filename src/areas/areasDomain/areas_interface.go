package areasDomain

import (
	"DIMISA/src/areas/areasDomain/areaEntity"
)

type AreasInterface interface {
	CreateArea(area *areaEntity.AreaEntity) error
	GetAllAreas() ([]*areaEntity.AreaEntity, error)
	UpdateArea(area *areaEntity.AreaEntity) error
	DeleteArea(id int32) error
	GetAreaByID(id int32) (*areaEntity.AreaEntity, error)
	GetFreeAreas() ([]*areaEntity.AreaEntity, error)
	GetAreasByCendis(id int32) ([]*areaEntity.AreaEntity, error)
}
