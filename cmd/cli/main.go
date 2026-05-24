package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pdfsec/internal/cli"
	"pdfsec/internal/core"
	"strings"
)

// Config holds all parsed CLI options.
type Config struct {
	Passin            string
	OwnerPassword     string
	Permissions       string
	InputFile         string
	OutputFile        string
	OwnerPasswordOnly bool
	PrintOwner        bool
	Verbose           bool
	ShowVersion       bool
	RndUsrPassword    bool
	RndOwnPassword    bool
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
  --no-reader-password      Set PDF permissions using only an owner password, without requiring a reader 
  							password to open the document

  --verbose                 Print logs
  --version                 Print version and exit.
  --help                    Show this help message.

EXAMPLES:
  # Encrypt with interactive password prompts
  %s document.pdf encrypted.pdf

  # Encrypt with passwords supplied directly
  %s --passin hunter2 --owner-password s3cr3t document.pdf encrypted.pdf

  # Restrict to print-only, no copying
  %s --no-reader-password --permissions print document.pdf encrypted.pdf

  # Verbose mode
  %s --verbose --no-reader-password document.pdf encrypted.pdf

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
`, core.AppName, core.AppVersion, core.AppName, core.AppName, core.AppName, core.AppName, core.AppName)
}

// parseFlags parses command-line flags, validates required arguments, and builds the Config struct.
func parseFlags() (*Config, error) {
	cfg := &Config{}

	flag.Usage = printHelp

	flag.StringVar(&cfg.Passin, "passin", "", "User (open) password")
	flag.StringVar(&cfg.OwnerPassword, "owner-password", "", "Owner password")
	flag.StringVar(&cfg.Permissions, "permissions", "", "Comma-separated permission list")
	flag.BoolVar(&cfg.PrintOwner, "print-owner-password", false, "Print owner password if randomly generated")
	flag.BoolVar(&cfg.OwnerPasswordOnly, "no-reader-password", false, "Set PDF permissions using only an owner password, without requiring a reader password to open the document")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Print logs")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Print version and exit")
	flag.Parse()

	if cfg.ShowVersion {
		fmt.Printf("%s version %s\n", core.AppName, core.AppVersion)
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
func getPassword(cfg *Config, log *cli.Logger) bool {
	var globRnd bool = false

	if !cfg.OwnerPasswordOnly {
		if cfg.Passin == "" {
			log.Info("No user password provided via flag - prompting")
			pwd, rnd, err := cli.ReadPasswordFromTTY("user")
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
	} else {
		log.Info("Using only an owner password for PDF permissions")
	}

	if cfg.OwnerPassword == "" {
		log.Info("No owner password provided via flag - prompting")
		pwd, rnd, err := cli.ReadPasswordFromTTY("owner")
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

// printGeneratedPasswords outputs randomly generated user/owner passwords when applicable.
func printGeneratedPasswords(cfg *Config) {
	if cfg.RndUsrPassword {
		fmt.Fprintf(os.Stderr, "Reader password: %s\n", cfg.Passin)
	}
	if cfg.RndOwnPassword && (cfg.PrintOwner || cfg.OwnerPasswordOnly) {
		fmt.Fprintf(os.Stderr, "Owner password:  %s\n", cfg.OwnerPassword)
	}
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log := cli.NewLogger(cfg.Verbose)
	log.Info(fmt.Sprintf("Starting %s v%s", core.AppName, core.AppVersion))

	if err := core.CheckFileExists(cfg.InputFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	perms, err := core.ParsePermissions(cfg.Permissions)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if strings.TrimSpace(cfg.Permissions) == "" {
		log.Info("No permissions specified - using defaults (print, copy, annotate)")
	} else {
		log.Info(fmt.Sprintf("Permissions set: %s", strings.TrimSpace(cfg.Permissions)))
	}

	printNewLine := getPassword(cfg, log)

	if printNewLine {
		fmt.Fprint(os.Stderr, "\n")
	}

	log.Info("Encrypting with AES-256...")
	err = core.EncryptPDF(core.EncryptOptions{
		InputFile: cfg.InputFile,
		OutputFile: core.ResolveOutputFile(
			cfg.InputFile,
			cfg.OutputFile,
		),
		UserPassword:  cfg.Passin,
		OwnerPassword: cfg.OwnerPassword,
		Permissions:   perms,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("Encrypted file written: %s", filepath.Base(cfg.OutputFile)))

	if cfg.OwnerPasswordOnly {
		fmt.Fprintln(os.Stderr, "WARNING: the PDF is not fully encrypted. Only PDF permissions are set, and some readers may ignore or bypass them.")
	}

	printGeneratedPasswords(cfg)

	log.Info("Done.")
}
