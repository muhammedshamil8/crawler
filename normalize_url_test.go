package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove scheme and trailing slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "no scheme or trailing slash",
			inputURL: "www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "h1 present",
			input:    "<html><body><h1>Test Title</h1></body></html>",
			expected: "Test Title",
		},
		{
			name:     "h2 present",
			input:    "<html><body><h2>Test Title</h2></body></html>",
			expected: "Test Title",
		},
		{
			name:     "no heading present",
			input:    "<html><body><p>No heading here</p></body></html>",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected heading: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "main present",
			input:    `<html><body><p>Outside paragraph.</p><main><p>Main paragraph.</p></main></body></html>`,
			expected: "Main paragraph.",
		},
		{
			name:     "no main present",
			input:    `<html><body><p>First paragraph.</p><p>Second paragraph.</p></body></html>`,
			expected: "First paragraph.",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected paragraph: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "absolute URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
		},
		{
			name:      "relative URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/path/to/page"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com/path/to/page"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputURL := tc.inputURL
			inputBody := tc.inputBody
			baseURL, err := url.Parse(inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getURLsFromHTML(inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := tc.expected
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "relative image URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "absolute image URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputURL := tc.inputURL
			inputBody := tc.inputBody
			baseURL, err := url.Parse(inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}
			actual, err := getImagesFromHTML(inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := tc.expected
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}


func TestExtractPageData(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  PageData
	}{
		{
			name:     "basic extraction",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<h1>Test Title</h1>
				<p>This is the first paragraph.</p>
				<a href="/link1">Link 1</a>
				<img src="/image1.jpg" alt="Image 1">
			</body></html>`,
			expected: PageData{
				URL:             "https://crawler-test.com",
				Heading:         "Test Title",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks:  []string{"https://crawler-test.com/link1"},
				ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
			},
		},
		{
			name:     "no heading or paragraph",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="/link1">Link 1</a>
				<img src="/image1.jpg" alt="Image 1">
			</body></html>`,
			expected: PageData{
				URL:             "https://crawler-test.com",
				Heading:         "",
				FirstParagraph: "",
				OutgoingLinks:  []string{"https://crawler-test.com/link1"},
				ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractPageData(tc.inputBody, tc.inputURL)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}	
		})
	}
}