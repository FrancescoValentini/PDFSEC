<# 
MIT License

Copyright (c) 2025 Francesco Valentini

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


Powershell script to quickly build a cross-platform GOLANG project
Author: Francesco Valentini
#> 

# Project Variables
$ProjectName = "pdfsec"
$OutputDir = "build"
$ReportFile = "build_report.md"

# List of common GOOS/GOARCH combinations
$Targets = @(
    "windows/amd64",
    "windows/386",
    "windows/arm",
    "linux/amd64",
    "linux/386",
    "linux/arm",
    "linux/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "freebsd/amd64",
    "openbsd/amd64"
)

# Create output folder if it does not exist
if (-not (Test-Path -Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# Initialize report
$reportLines = @()
$reportLines += "# $ProjectName"
$reportLines += "| FILE | PLATFORM | ARCHITECTURE | SHA256 |"
$reportLines += "|------|----------|--------------|--------|"

foreach ($target in $Targets) {
    $parts = $target -split "/"
    $GOOS = $parts[0]
    $GOARCH = $parts[1]

    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH

    $ext = ""
    if ($GOOS -eq "windows") {
        $ext = ".exe"
    }

    $fileName = "$ProjectName-$GOOS-$GOARCH$ext"
    $outputFile = "$OutputDir\$fileName"
    Write-Host "Building for: $GOOS/$GOARCH - $outputFile"

    go build -ldflags "-s -w" -o $outputFile

    if ($LASTEXITCODE -ne 0) {
        Write-Warning "Build failed for $GOOS/$GOARCH"
        continue
    }

    # Compute SHA256 hash
    $sha256 = Get-FileHash -Algorithm SHA256 -Path $outputFile | Select-Object -ExpandProperty Hash

    # Append to report
    $reportLines += "| $fileName | $GOOS | $GOARCH | $sha256 |"
}

# Write the report to file
$reportContent = $reportLines -join "`n"
Set-Content -Path $ReportFile -Value $reportContent -Encoding UTF8

Write-Host "Done. Report generated at $ReportFile"