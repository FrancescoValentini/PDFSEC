package core

import (
	"fmt"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// defaultPermissions returns the default set of PDF permissions (print, extract, annotate/fill forms).
func DefaultPermissions() model.PermissionFlags {
	return model.PermissionPrintRev3 |
		model.PermissionExtract |
		model.PermissionModAnnFillForm
}

// permissionMap returns a key-value mapping for the PDF's permissions
// https://pkg.go.dev/github.com/pdfcpu/pdfcpu@v0.12.1/pkg/pdfcpu/model#PermissionFlags
func PermissionMap() map[string]model.PermissionFlags {
	return map[string]model.PermissionFlags{
		"print":     model.PermissionPrintRev3,
		"print-low": model.PermissionPrintRev2,
		"modify":    model.PermissionModify,
		"copy":      model.PermissionExtract,
		"extract":   model.PermissionExtractRev3,
		"annotate":  model.PermissionModAnnFillForm,
		"fill":      model.PermissionFillRev3,
		"assemble":  model.PermissionAssembleRev3,
	}
}

// parsePermissions parses a comma-separated permission string into pdfcpu permission flags.
func ParsePermissions(s string) (model.PermissionFlags, error) {
	if strings.TrimSpace(s) == "" {
		return DefaultPermissions(), nil
	}

	permMap := PermissionMap()

	var perms model.PermissionFlags
	items := strings.Split(s, ",")
	for _, item := range items {
		key := strings.TrimSpace(strings.ToLower(item))
		p, ok := permMap[key]
		if !ok {
			return 0, fmt.Errorf("unknown permission %q - valid choices: print, print-low, copy, extract, modify, annotate, fill, assemble", item)
		}
		perms |= p
	}
	return perms, nil
}
