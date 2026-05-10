// src/core/utils/query_builder.go
package utils

import (
	"strings"
)

func BuildPrefixWhere(prefixes []string, exactKeys []string, exclusions []string) (string, []interface{}) {
	args := []interface{}{}
	parts := []string{}

	for _, p := range prefixes {
		parts = append(parts, "clave_med LIKE ?")
		args = append(args, p+"%")
	}

	// Exact keys adicionales (ej: SC-0000080 en med)
	if len(exactKeys) > 0 {
		placeholders := make([]string, len(exactKeys))
		for i, k := range exactKeys {
			placeholders[i] = "?"
			args = append(args, k)
		}
		parts = append(parts, "clave_med IN ("+strings.Join(placeholders, ", ")+")")
	}

	prefixClause := "(" + strings.Join(parts, " OR ") + ")"

	if len(exclusions) == 0 {
		return prefixClause, args
	}

	placeholders := make([]string, len(exclusions))
	for i, e := range exclusions {
		placeholders[i] = "?"
		args = append(args, e)
	}

	return prefixClause + " AND clave_med NOT IN (" + strings.Join(placeholders, ", ") + ")", args
}
