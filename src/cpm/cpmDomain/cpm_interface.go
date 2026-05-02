package cpmDomain

import cpmEntity "DIMISA/src/cpm/cpmDomain/cpmEntity"

type CpmInterface interface {
	GetCpm() (cpmEntity.CpmEntity, error)
}
