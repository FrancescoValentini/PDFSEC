package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

const (
	appName    = "PDFSEC"
	appVersion = "1.0.0"
	fileSuffix = "-protected"
)

// Config holds all parsed CLI options.
type Config struct {
	Passin         string
	OwnerPassword  string
	Permissions    string
	InputFile      string
	OutputFile     string
	PrintOwner     bool
	Verbose        bool
	ShowVersion    bool
	RndUsrPassword bool
	RndOwnPassword bool
}

func printHelp() {
	fmt.Printf(`%s v%s - Encrypt PDF files with AES-256 and fine-grained permissions.

USAGE:
  %s [flags] <input.pdf> <output.pdf>

FLAGS:
  --passin <pwd>            User (open) password. Prompted interactively if omitted.
  --owner-password <pwd>    Owner password. Prompted interactively if omitted.
  --permissions <list>      Comma-separated permissions to grant (default: print,copy,annotate).
                            Choices: print, print-low, copy, extract, modify, annotate, fill, assemble
  --print-owner-password    Print the owner password to stderr when randomly generated.
  --verbose                 Print logs
  --version                 Print version and exit.
  --help                    Show this help message.

EXAMPLES:
  # Encrypt with interactive password prompts
  %s document.pdf encrypted.pdf

  # Encrypt with passwords supplied directly
  %s --passin hunter2 --owner-password s3cr3t document.pdf encrypted.pdf

  # Restrict to print-only, no copying
  %s --passin "" --permissions print document.pdf encrypted.pdf

  # Verbose mode
  %s --verbose --passin "" document.pdf encrypted.pdf

PERMISSIONS:
  print       Allow high-quality printing
  print-low   Allow low-quality (draft) printing only
  copy        Allow copying text and images
  extract     Allow content extraction (accessibility)
  modify      Allow document modifications
  annotate    Allow adding annotations and filling forms
  fill        Allow filling form fields
  assemble    Allow inserting, rotating, or deleting pages

AUTHOR:
  Francesco Valentini (C) 2026
`, appName, appVersion, appName, appName, appName, appName, appName)
}

// parseFlags parses command-line flags, validates required arguments, and builds the Config struct.
func parseFlags() (*Config, error) {
	cfg := &Config{}

	flag.Usage = printHelp

	flag.StringVar(&cfg.Passin, "passin", "", "User (open) password")
	flag.StringVar(&cfg.OwnerPassword, "owner-password", "", "Owner password")
	flag.StringVar(&cfg.Permissions, "permissions", "", "Comma-separated permission list")
	flag.BoolVar(&cfg.PrintOwner, "print-owner-password", false, "Print owner password if randomly generated")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Print logs")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Print version and exit")
	flag.Parse()

	if cfg.ShowVersion {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		printHelp()
		return nil, fmt.Errorf("error: input and output file paths are required")
	}

	cfg.InputFile = args[0]
	cfg.InputFile = args[0]
	if len(args) >= 2 {
		cfg.OutputFile = args[1]
	}
	return cfg, nil
}

// getPassword retrieves user and owner passwords either from CLI flags or via interactive TTY prompts.
func getPassword(cfg *Config, log *Logger) bool {
	var globRnd bool = false
	if cfg.Passin == "" {
		log.Info("No user password provided via flag - prompting")
		pwd, rnd, err := ReadPasswordFromTTY("user")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading user password,\ntry again, if it still fails try providing the password via the appropriate command line flag.")
			os.Exit(1)
		}
		cfg.Passin = pwd
		cfg.RndUsrPassword = rnd
		globRnd = rnd
	} else {
		log.Info("User password supplied via flag")
	}

	if cfg.OwnerPassword == "" {
		log.Info("No owner password provided via flag - prompting")
		pwd, rnd, err := ReadPasswordFromTTY("owner")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading owner password,\ntry again, if it still fails try providing the password via the appropriate command line flag.")
			os.Exit(1)
		}
		cfg.OwnerPassword = pwd
		cfg.RndOwnPassword = rnd
		globRnd = rnd
	} else {
		log.Info("Owner password supplied via flag")
	}
	return globRnd
}

// defaultPermissions returns the default set of PDF permissions (print, extract, annotate/fill forms).
func defaultPermissions() model.PermissionFlags {
	return model.PermissionPrintRev3 |
		model.PermissionExtract |
		model.PermissionModAnnFillForm
}

// permissionMap returns a key-value mapping for the PDF's permissions
// https://pkg.go.dev/github.com/pdfcpu/pdfcpu@v0.12.1/pkg/pdfcpu/model#PermissionFlags
func permissionMap() map[string]model.PermissionFlags {
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
func parsePermissions(s string, log *Logger) (model.PermissionFlags, error) {
	if strings.TrimSpace(s) == "" {
		log.Info("No permissions specified - using defaults (print, copy, annotate)")
		return defaultPermissions(), nil
	}

	permMap := permissionMap()

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

	log.Info(fmt.Sprintf("Permissions set: %s", strings.TrimSpace(s)))
	return perms, nil
}

// printGeneratedPasswords outputs randomly generated user/owner passwords when applicable.
func printGeneratedPasswords(cfg *Config) {
	if cfg.RndUsrPassword {
		fmt.Fprintf(os.Stderr, "Reader password: %s\n", cfg.Passin)
	}
	if cfg.RndOwnPassword && cfg.PrintOwner {
		fmt.Fprintf(os.Stderr, "Owner password:  %s\n", cfg.OwnerPassword)
	}
}

// resolveOutputFile determines the output filename, adding a suffix if none is provided.
func resolveOutputFile(input, output string) string {
	if output != "" {
		return output
	}
	ext := filepath.Ext(input)
	name := strings.TrimSuffix(input, ext)
	return name + fileSuffix + ext
}

// encryptPDF encrypts the input PDF using AES-256 with the provided passwords and permissions.
func encryptPDF(inputFile, outputFile, userPwd, ownerPwd string, permissions model.PermissionFlags, log *Logger) error {
	log.Info(fmt.Sprintf("Reading input file: %s", filepath.Base(inputFile)))

	conf := model.NewAESConfiguration(userPwd, ownerPwd, 256)
	conf.Permissions = permissions

	log.Info("Encrypting with AES-256...")
	if err := api.EncryptFile(inputFile, outputFile, conf); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Encrypted file written: %s", filepath.Base(outputFile)))
	return nil
}

// checkFileExists returns an error if the file does not exist or if the path points to a folder
func checkFileExists(path string) error {
	filePointer, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	if filePointer.IsDir() {
		return fmt.Errorf("expected a single pdf file not a folder.")
	}
	return err
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log := NewLogger(cfg.Verbose)
	log.Info(fmt.Sprintf("Starting %s v%s", appName, appVersion))

	if err := checkFileExists(cfg.InputFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	cfg.OutputFile = resolveOutputFile(cfg.InputFile, cfg.OutputFile)
	printNewLine := getPassword(cfg, log)

	if printNewLine {
		fmt.Fprint(os.Stderr, "\n")
	}

	perms, err := parsePermissions(cfg.Permissions, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid permissions: %v\n", err)
		os.Exit(1)
	}

	if err := encryptPDF(cfg.InputFile, cfg.OutputFile, cfg.Passin, cfg.OwnerPassword, perms, log); err != nil {
		log.Error(fmt.Sprintf("Encryption failed: %v", err))
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printGeneratedPasswords(cfg)

	log.Info("Done.")
}
