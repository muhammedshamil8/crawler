# Go Web Crawler 🕸️

A concurrent, high-performance Web Crawler written in **Go**, developed as part of the **Boot.dev** backend development path.

---

## 🎯 Learning Goals

This project was built to practice foundational Go concepts and web development techniques:
- 🛠️ **Local Go Development & Tooling**: Environment setup, modules (`go.mod`), dependencies, and package management.
- 🌐 **HTTP Requests in Go**: Making robust outbound HTTP requests, handling response headers, body reading, and status code verification.
- 📄 **HTML Parsing in Go**: Extracting structured data (headings, paragraphs, links, images) from raw HTML using Go's standard library and `goquery`.
- 🧪 **Unit Testing**: Writing comprehensive table-driven tests with Go's built-in `testing` package.
- ⚡ **Concurrency Controls**: Utilizing goroutines, channels (semaphore pattern), and `sync.WaitGroup` to safely execute concurrent crawls up to specified concurrency limits.

---

## ✨ Features

- **Concurrent Web Crawling**: Multi-threaded crawling powered by goroutines and channels.
- **Configurable Limits**: Restrict max concurrent requests and total max pages crawled.
- **Domain-Scoped Crawling**: Restricts crawler to only follow internal links within the same base host domain.
- **URL Normalization**: Normalizes URLs (stripping protocols and trailing slashes) to prevent redundant crawls.
- **Data Extraction**:
  - Main headings (`<h1>` / `<h2>`)
  - First paragraph content (`<p>` within `<main>` or body)
  - All outgoing hyperlinks (`<a>`)
  - All image URLs (`<img>`)
- **Structured JSON Export**: Formats crawl results into a sorted `report.json` file.

---

## 📁 Project Structure

```text
crawler/
├── main.go               # Entry point and CLI argument processing
├── extract_page_data.go  # Core crawler logic, HTTP fetching, and HTML parsing
├── normalize_url.go      # URL normalization utility functions
├── normalize_url_test.go # Comprehensive unit tests for parser and normalization logic
├── json_report.go        # JSON formatting and export handling
├── go.mod                # Go module definition
├── go.sum                # Go checksum database
└── README.md             # Project documentation
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) `1.22+` (or `1.26+`) installed on your system.

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/muhammedshamil8/crawler.git
   cd crawler
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   ```

---

## 💻 Usage

Run the crawler using `go run .` with the target URL, maximum concurrency limit, and maximum pages to crawl:

```bash
go run . <URL> <maxConcurrency> <maxPages>
```

### Example

```bash
go run . https://wagslane.dev 5 25
```

Output:
```text
starting crawl of: https://wagslane.dev
Crawling: wagslane.dev
Crawling: wagslane.dev/posts
...
============== CRAWLING COMPLETED ==============
URL: https://wagslane.dev, Heading: Wagslane, First Paragraph: Welcome to my blog...
```

This will generate a sorted `report.json` in the root directory detailing all scraped page information.

### Build Executable

To build a binary executable:

```bash
go build -o crawler
./crawler https://wagslane.dev 5 25
```

---

## 🧪 Running Tests

Execute the unit test suite with verbose output:

```bash
go test -v ./...
```

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for details.
