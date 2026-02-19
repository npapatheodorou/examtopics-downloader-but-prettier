# ExamTopics Downloader (Enhanced Edition)

<p align="center">
  <a href="https://github.com/thatonecodes/examtopics-downloader">
    <img src="https://img.shields.io/badge/Forked%20from-thatonecodes-blue?style=flat-square" alt="Forked from thatonecodes/examtopics-downloader">
  </a>
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Built%20with-Go-blue?style=flat-square" alt="Built with Go">
  </a>
  <a href="https://github.com/npapatheodorou/examtopics-downloader-but-prettier/releases/latest">
    <img src="https://img.shields.io/github/v/release/npapatheodorou/examtopics-downloader-but-prettier?include_prereleases&label=Latest%20Release&style=flat-square" alt="Latest Release">
  </a>
</p>

> **This project is a fork of [thatonecodes/examtopics-downloader](https://github.com/thatonecodes/examtopics-downloader)**  
> Special thanks to [@thatonecodes](https://github.com/thatonecodes) for creating the original tool that made this possible.

---

## What is this?

**ExamTopics Downloader (Enhanced Edition)** is a user-friendly command-line tool that lets you download exam questions from [ExamTopics](https://examtopics.com/) and save them in a clean, readable format. Whether you're preparing for AWS, Azure, Google Cloud, CompTIA, or dozens of other certifications, this tool makes it easy to access and study exam questions offline.

---

## Features

| Feature | Description |
|---------|-------------|
| **Download Exams** | Fetch questions from any exam available on ExamTopics |
| **Clean HTML Output** | Beautiful, readable HTML format that's easy on the eyes |
| **Interactive Selection** | Browse and select exams with an easy-to-use menu |
| **Exam Simulation Mode** | Practice exams interactively with an HTML-based simulation |
| **Ready-to-Use .exe** | No need to install Go - download the pre-built executable and run it |
| **Cross-Platform** | Build from source for Windows, macOS, or Linux |

---

## Quick Start (Windows)

### Option 1: Use the Pre-built .exe (Recommended)

1. Go to the [Releases](https://github.com/npapatheodorou/examtopics-downloader-but-prettier/releases) page
2. Download the latest `examtopics-downloader-windows-amd64.exe`
3. Double-click to run - no installation needed!

### Option 2: Build from Source

If you have [Go installed](https://go.dev/dl/):

```bash
git clone https://github.com/npapatheodorou/examtopics-downloader-but-prettier.git
cd examtopics-downloader-but-prettier
go build -o examtopics-downloader.exe ./cmd/main.go
```

Or use the included build script:

```bash
build.bat
```

---

## How to Use

### Running the Tool

Simply double-click the `.exe` file (or run from terminal):

```
examtopics-downloader-windows-amd64.exe
```

### Step-by-Step

1. **Select a Provider**  
   Choose your certification vendor (e.g., AWS, Azure, CompTIA)

2. **Select an Exam**  
   Pick the specific exam or exam series you want to download

3. **Wait for Download**  
   The tool will fetch all questions and save them

4. **Open the Output**  
   Find the generated `.html` file in the same folder and open it in your browser

### Output Files

- **`provider_examname.html`** - The main exam output in HTML format
- Open the HTML file in any browser to view, print, or study

---

## Sample Workflow

```
============================================================
 ExamTopics Downloader - Interactive Exam Extractor
============================================================

[*] Loading providers from ExamTopics...
[OK] Done. Found 45 provider(s) in 4s.

--------------------------------------------------------
 Available Exams
--------------------------------------------------------
 Showing 1 of 1
 Filter: ""

   1. aws-saa-co03 (SAA-C03 - AWS Solutions Architect)

 Commands: [number] select | /text filter | / clear | /refresh refetch
Select> 1
[INFO] Starting extraction for aws / aws-saa-co03...
[OK] Successfully saved output: aws_saa-co03.html
```

---

## Improvements Over the Original

This fork includes several enhancements over the original [examtopics-downloader](https://github.com/thatonecodes/examtopics-downloader):

- **Fixed Missing Exams** - Added support for exams that were previously not accessible
- **Resolved Minor Issues** - Bug fixes and stability improvements
- **Better Usability** - Improved user experience for non-technical users
- **Compiled .exe Release** - No need to install Go; download and run
- **HTML Exam Simulation** - Practice exams in an interactive browser-based format
- **Enhanced Filtering** - Search and filter exams with the `/text` command

---

## Technical Details

### Requirements

- **Windows**: No additional requirements (pre-built .exe)
- **From Source**: Go 1.21+

### Building

```bash
# Build for Windows (amd64)
go build -o examtopics-downloader.exe ./cmd/main.go

# Build for other platforms
GOOS=linux GOARCH=amd64 go build -o examtopics-downloader ./cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o examtopics-downloader ./cmd/main.go
```

### Output Formats

The tool generates clean, styled HTML that works in any browser. The HTML includes:
- Question text with proper formatting
- Multiple choice answers (A, B, C, D...)
- Correct answer highlights
- Explanation sections
- Clean, modern styling

---

## Disclaimer

This tool is for **educational purposes only**. ExamTopics content is copyrighted material. Please support ExamTopics by using their site directly if you can.

---

## License

See [LICENSE](LICENSE) for details.

---

## Credits

- **Original Author**: [@thatonecodes](https://github.com/thatonecodes) - Thank you for creating this amazing tool!
- **This Fork**: Enhanced and maintained by the community

---

<p align="center">
  <strong>Good luck with your certification studies!</strong>
</p>
