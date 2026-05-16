package utils

import (
	"DIMISA/src/core/config"
	"strings"
)

func EsMedicamento(clave string) bool {
	for _, prefix := range config.MedPrefixes {
		if strings.HasPrefix(clave, prefix) {
			return true
		}
	}
	for _, exact := range config.MedExactKeys {
		if clave == exact {
			return true
		}
	}
	return false
}
