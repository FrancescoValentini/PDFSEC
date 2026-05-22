# PDFSEC

PDFSEC is a command-line utility for securing PDF documents with AES-256 and fine-grained permission control.

## Features

- AES-256 encryption with separate user and owner passwords
- Granular permission flags (print, copy, annotate, and more)
- If no password is provided, a random password will be generated automatically using a human-friendly character set
- Permission-only mode: restrict what readers can do without requiring a password to open the file

## Usage

```
pdfsec [flags] <input.pdf> [output.pdf]
```

If no output file is given, the encrypted file is saved as `input-protected.pdf`.

### Flags

| Flag                     | Description                                                                                                 |
| ------------------------ | ----------------------------------------------------------------------------------------------------------- |
| `--passin <pwd>`         | User (open) password. Prompted interactively if omitted.                                                    |
| `--owner-password <pwd>` | Owner password. Prompted interactively if omitted.                                                          |
| `--permissions <list>`   | Comma-separated permissions to grant (see below).                                                           |
| `--print-owner-password` | Print the owner password to stderr when randomly generated.                                                 |
| `--no-reader-password`   | Set PDF permissions using only an owner password, without requiring a reader password to open the document. |
| `--verbose`              | Show detailed logs.                                                                                         |
| `--version`              | Print version and exit.                                                                                     |
| `--help`                 | Show help message.                                                                                          |

### Permissions

| Value       | Description                        |
| ----------- | ---------------------------------- |
| `print`     | High-quality printing              |
| `print-low` | Draft/low-quality printing only    |
| `copy`      | Copy text and images               |
| `extract`   | Content extraction (accessibility) |
| `modify`    | General document modifications     |
| `annotate`  | Add annotations and fill forms     |
| `fill`      | Fill form fields only              |
| `assemble`  | Insert, rotate, or delete pages    |

Default permissions when `--permissions` is omitted: **print, copy, annotate**.


## Examples

**Interactive mode** — prompts for both passwords:
```bash
pdfsec document.pdf encrypted.pdf
```

**Supply passwords directly:**
```bash
pdfsec --passin hunter2 --owner-password s3cr3t document.pdf encrypted.pdf
```

**Read-only: allow printing, no copying:**
```bash
pdfsec --passin "" --permissions print document.pdf encrypted.pdf
```

**No user password, restrict to print and annotate only:**
```bash
pdfsec --passin "" --permissions print,annotate document.pdf
```

**Permissions only, no password required to open:**
```bash
pdfsec --no-reader-password --permissions print document.pdf encrypted.pdf
```

> [!WARNING] 
> When using `--no-reader-password`, the PDF is not fully encrypted. Only PDF permissions are set, and some readers may ignore or bypass them.

**Print randomly generated owner password to stderr:**
```bash
pdfsec --print-owner-password document.pdf encrypted.pdf
```

**Verbose output:**
```bash
pdfsec --verbose --no-reader-password document.pdf encrypted.pdf
```

---

## Dependencies

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) — PDF processing library for Go

## ⚠️ Disclaimer

This is a personal project, built for my own use and shared as-is.

I am **not responsible** for any damage, data loss, or corrupted files that may result from using this tool. Always keep a backup of your original PDFs before encrypting them. Use at your own risk.