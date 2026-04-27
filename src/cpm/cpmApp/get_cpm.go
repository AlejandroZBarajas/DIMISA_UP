package cpmApp

import (
	cpmDomain "DIMISA/src/cpm/cpmDomain"
	cpmEntity "DIMISA/src/cpm/cpmDomain/cpmEntity"
)

type GetCpm struct {
	Repo cpmDomain.CpmInterface
}

func (uc *GetCpm) Execute() (cpmEntity.CpmEntity, error) {
	return uc.Repo.GetCpm()
}
