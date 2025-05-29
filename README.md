# JNLP Resource Downloader (jrdl)

A command-line utility that downloads JAR files listed in Java Network Launch Protocol (JNLP) files.

## Description

This utility parses JNLP files and downloads all the JAR resources referenced within them to a local directory. It's useful for offline access to Java applications or for archiving purposes.

## Installation

Build the application using Go:

```bash
go build -o jrdl .
```

## Usage

```
jrdl [OPTIONS] <jnlp-file> [download-dir]
```

### Arguments

- `<jnlp-file>` - Path to the JNLP file to download (required)
- `[download-dir]` - Directory to download the JAR files to (optional, defaults to "downloads")

### Options

- `--failed-download-exit` - Exit with non-zero code if a download fails
- `--failed-file-creation-exit` - Exit with non-zero code if a file cannot be created
- `--failed-file-write-exit` - Exit with non-zero code if a file cannot be written

### Examples

Download JARs from a JNLP file to the default directory:
```bash
jrdl application.jnlp
```

Download JARs to a specific directory:
```bash
jrdl application.jnlp /path/to/output
```

Exit immediately on any download failure:
```bash
jrdl --failed-download-exit application.jnlp
```

## How It Works

1. Parses the provided JNLP file to extract JAR resource URLs
2. Creates a download directory named after the JNLP application title
3. Downloads each JAR file from the resolved URLs
4. Saves the JAR files to the specified directory

## Error Handling

By default, the application continues downloading even if individual files fail. Use the error flags to change this behavior:

- Download failures are logged but don't stop execution unless `--failed-download-exit` is used
- File creation failures are logged but don't stop execution unless `--failed-file-creation-exit` is used  
- File write failures are logged but don't stop execution unless `--failed-file-write-exit` is used

## JNLP File Structure

The application expects JNLP files with the following structure:

```xml
<jnlp codebase="http://example.com/" href="app.jnlp">
  <information>
    <title>Application Name</title>
  </information>
  <resources>
    <jar href="app.jar"/>
    <jar href="lib/dependency.jar"/>
  </resources>
</jnlp>
```

## Requirements

- Go 1.16 or later
- Internet connection for downloading JAR files

## License

This project is provided as-is without any specific license.
