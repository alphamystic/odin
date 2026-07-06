package sourcemap

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// DefaultTempDir holds the fall-back root folder configuration.
const DefaultTempDir = "runtime_temp"

// SourceMap represents the parsed raw JSON content schema.
type SourceMap struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
}

// ExtractedFile holds the relative internal path and contents of an unbundled asset.
type ExtractedFile struct {
	Path    string
	Content string
}

// Result tracks execution metrics and data payloads for the caller.
type Result struct {
	SourceURI   string          // Original target file, folder, or URL string
	SavedFolder string          // The path on disk where files were flushed
	Files       []ExtractedFile // In-memory map array of individual elements
}

// Config manages the structural input layout with clean fallbacks.
type Config struct {
	URL         string   // URL or path to a direct Sourcemap file
	JSURL       string   // URL to a JavaScript file containing a sourcemap
	Dir         string   // Local directory of .map files to process recursively
	Proxy       string   // Optional upstream proxy URL string
	InsecureTLS bool     // Bypasses certificate checks if set to true
	Headers     []string // Custom array strings formatting HTTP Requests ("Key: Value")
	BaseOutDir  string   // Target directory on disk. If empty, defaults to "./runtime_temp/run_<timestamp>"
}

// Extractor is the modular controller orchestrating the operations.
type Extractor struct {
	cfg      *Config
	proxyURL *url.URL
}

// NewExtractor initializes the architecture, normalizes default configurations, and sets fallbacks.
func NewExtractor(cfg Config) (*Extractor, error) {
	// Mutual exclusivity routing validations
	modes := 0
	if cfg.URL != "" { modes++ }
	if cfg.JSURL != "" { modes++ }
	if cfg.Dir != "" { modes++ }

	if modes == 0 {
		return nil, errors.New("invalid configuration: choose an input driver (URL, JSURL, or Dir)")
	}
	if modes > 1 {
		return nil, errors.New("invalid configuration: input configurations are mutually exclusive")
	}

	// Dynamic Path Fallback Generation
	if cfg.BaseOutDir == "" {
		timestamp := time.Now().Format("20060102_150405")
		cfg.BaseOutDir = filepath.Join(DefaultTempDir, fmt.Sprintf("run_%s", timestamp))
	}

	var parsedProxy *url.URL
	if cfg.Proxy != "" {
		p, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, fmt.Errorf("failed parsing downstream proxy configuration: %w", err)
		}
		parsedProxy = p
	}

	return &Extractor{
		cfg:      &cfg,
		proxyURL: parsedProxy,
	}, nil
}

// Extract triggers the underlying source assembly logic based on the config input.
func (e *Extractor) Extract() ([]Result, error) {
	var totalResults []Result

	// Operational Pathway 1: Bulk File-System Folders
	if e.cfg.Dir != "" {
		var mapFiles []string
		err := filepath.Walk(e.cfg.Dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".map") {
				mapFiles = append(mapFiles, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed recursive crawl sequence: %w", err)
		}
		if len(mapFiles) == 0 {
			return nil, fmt.Errorf("no source maps found in location target: %s", e.cfg.Dir)
		}

		for _, mapFile := range mapFiles {
			sm, err := e.fetchSourceMap(mapFile)
			if err != nil {
				log.Printf("[!] Parsing issue with target %s: %v. Skipping.", mapFile, err)
				continue
			}

			// Segment distinct outputs cleanly into subfolders within our runtime block
			subDirName := strings.TrimSuffix(filepath.Base(mapFile), filepath.Ext(mapFile))
			targetDir := filepath.Join(e.cfg.BaseOutDir, subDirName)

			res, err := e.buildArtifacts(sm, mapFile, targetDir)
			if err != nil {
				log.Printf("[!] Directory assembly error for %s: %v. Skipping.", mapFile, err)
				continue
			}
			totalResults = append(totalResults, res)
		}
		return totalResults, nil
	}

	// Operational Pathway 2: Single Targets (Web URL endpoints or explicit file maps)
	var sm SourceMap
	var err error
	var contextTarget string

	if e.cfg.URL != "" {
		contextTarget = e.cfg.URL
		sm, err = e.fetchSourceMap(e.cfg.URL)
	} else if e.cfg.JSURL != "" {
		contextTarget = e.cfg.JSURL
		sm, err = e.fetchSourceMapFromJS(e.cfg.JSURL)
	}

	if err != nil {
		return nil, err
	}

	res, err := e.buildArtifacts(sm, contextTarget, e.cfg.BaseOutDir)
	if err != nil {
		return nil, err
	}
	totalResults = append(totalResults, res)

	return totalResults, nil
}

func (e *Extractor) fetchSourceMap(source string) (m SourceMap, err error) {
	var body []byte
	u, err := url.ParseRequestURI(source)

	if err != nil || (u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "data") {
		body, err = os.ReadFile(source)
		if err != nil {
			return m, err
		}
	} else if u.Scheme == "data" {
		chunks := strings.Split(u.Opaque, ",")
		if len(chunks) < 2 {
			return m, errors.New("malformed data URI string sequence")
		}
		body, err = base64.StdEncoding.DecodeString(chunks[1])
		if err != nil {
			return m, fmt.Errorf("failed base64 sequence decode: %w", err)
		}
	} else {
		body, err = e.executeHTTPRequest(u.String())
		if err != nil {
			return m, err
		}
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, fmt.Errorf("json validation signature error: %w", err)
	}
	return m, nil
}

func (e *Extractor) fetchSourceMapFromJS(jsurl string) (m SourceMap, err error) {
	u, err := url.ParseRequestURI(jsurl)
	if err != nil {
		return m, fmt.Errorf("invalid JavaScript asset string URI: %w", err)
	}

	body, err := e.executeHTTPRequest(u.String())
	if err != nil {
		return m, err
	}

	var sourceMapLocation string
	// Check standard map headers
	if sourceMapLocation = e.extractHeaderMap(jsurl); sourceMapLocation == "" {
		re := regexp.MustCompile(`\/\/[@#] sourceMappingURL=(.*)`)
		matches := re.FindAllSubmatch(body, -1)
		if len(matches) > 0 {
			sourceMapLocation = string(matches[len(matches)-1][1])
		}
	}

	if sourceMapLocation == "" {
		return m, errors.New("unable to identify any valid mapping endpoints inside target script file")
	}

	mapURL, err := url.ParseRequestURI(sourceMapLocation)
	if err != nil {
		mapURL, err = u.Parse(sourceMapLocation)
		if err != nil {
			return m, err
		}
	}

	return e.fetchSourceMap(mapURL.String())
}

func (e *Extractor) executeHTTPRequest(targetURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{}
	if e.cfg.InsecureTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if e.proxyURL != nil {
		transport.Proxy = http.ProxyURL(e.proxyURL)
	}

	if len(e.cfg.Headers) > 0 {
		rawHeaders := strings.Join(e.cfg.Headers, "\r\n") + "\r\n\r\n"
		tpReader := textproto.NewReader(bufio.NewReader(strings.NewReader(rawHeaders)))
		mimeHeader, err := tpReader.ReadMIMEHeader()
		if err == nil {
			req.Header = http.Header(mimeHeader)
		}
	}

	client := http.Client{Transport: transport, Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status transaction return code: %d", res.StatusCode)
	}

	return io.ReadAll(res.Body)
}

func (e *Extractor) extractHeaderMap(urlStr string) string {
	req, _ := http.NewRequest("HEAD", urlStr, nil)
	tr := &http.Transport{}
	if e.cfg.InsecureTLS { tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} }
	if e.proxyURL != nil { tr.Proxy = http.ProxyURL(e.proxyURL) }

	client := http.Client{Transport: tr, Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if loc := resp.Header.Get("SourceMap"); loc != "" { return loc }
	return resp.Header.Get("X-SourceMap")
}

func (e *Extractor) buildArtifacts(sm SourceMap, origin string, targetFolder string) (Result, error) {
	res := Result{
		SourceURI:   origin,
		SavedFolder: targetFolder,
		Files:       make([]ExtractedFile, 0, len(sm.Sources)),
	}

	if len(sm.Sources) == 0 || len(sm.SourcesContent) == 0 {
		return res, errors.New("empty map context array records")
	}

	if err := os.MkdirAll(targetFolder, 0700); err != nil {
		return res, fmt.Errorf("failed to safely generate run-time directory blocks: %w", err)
	}

	for i, innerPath := range sm.Sources {
		if i >= len(sm.SourcesContent) { break }

		innerPath = "/" + innerPath
		if runtime.GOOS == "windows" {
			m1 := regexp.MustCompile(`[?%*|:"<>]`)
			innerPath = m1.ReplaceAllString(innerPath, "")
		}

		cleanPath := filepath.Clean(innerPath)
		filePayload := sm.SourcesContent[i]

		res.Files = append(res.Files, ExtractedFile{
			Path:    cleanPath,
			Content: filePayload,
		})

		destinationDiskPath := filepath.Join(targetFolder, cleanPath)
		if err := os.MkdirAll(filepath.Dir(destinationDiskPath), 0700); err == nil {
			_ = os.WriteFile(destinationDiskPath, []byte(filePayload), 0600)
		}
	}

	return res, nil
}