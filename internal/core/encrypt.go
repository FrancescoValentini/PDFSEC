package core

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type EncryptOptions struct {
	InputFile     string
	OutputFile    string
	UserPassword  string
	OwnerPassword string
	Permissions   model.PermissionFlags
}

// encryptPDF encrypts the input PDF using AES-256 with the provided passwords and permissions.
func EncryptPDF(opts EncryptOptions) error {
	conf := model.NewAESConfiguration(
		opts.UserPassword,
		opts.OwnerPassword,
		256,
	)

	conf.Permissions = model.PermissionFlags(opts.Permissions)

	return api.EncryptFile(
		opts.InputFile,
		opts.OutputFile,
		conf,
	)
}
