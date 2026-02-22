package reconjs

/*
 This package automates doing recon data foe webapps given a few pre installed tools.
 This should be preinstalled and the data return values include:
 Mainly Vulnerabilities,
 JS Files
 Endpoints/URLs
 api_keys
 AWSKEys
*/
import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
//	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

// TIMEOUT for external commands
const CmdTimeout = 3600 * 3 * time.Second

// ReconData is the central structure that holds all recon results for one domain
type ReconData struct {
	Domain      string         `json:"domain"`
	JSFiles     []string       `json:"js_files,omitempty"`
	Endpoints   []string       `json:"endpoints,omitempty"`
	Secrets     SecretResults  `json:"secrets,omitempty"`
	HTTPChecks  []HTTPCheck    `json:"http_checks,omitempty"`
	FuzzResults []FFUFResult   `json:"fuzz_results,omitempty"`
	IDORChecks  []IDORResult   `json:"idor_checks,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
}

// SecretResults groups secret-like findings
type SecretResults struct {
	APIKeys  []string `json:"api_keys,omitempty"`
	AWSKeys  []string `json:"aws_keys,omitempty"`
	JsonLeak []string `json:"json_leaks,omitempty"`
	URLs     []string `json:"urls,omitempty"`
	Raw      []string `json:"raw,omitempty"`
}

// HTTPCheck is a single httpx result line parsed minimally
type HTTPCheck struct {
	Line string `json:"line"`
}

// FFUFResult is a stub for fuzz results
type FFUFResult struct {
	Raw string `json:"raw"`
}

// IDORResult captures a simple IDOR check outcome
type IDORResult struct {
	Endpoint      string `json:"endpoint"`
	InjectedValue string `json:"injected_value"`
	Status        int    `json:"status"`
	BodySnippet   string `json:"body_snippet,omitempty"`
}

//
// ----------------- Utilities -----------------
//

// runCmd runs a shell command (with bash -c when cmd contains pipes) and returns stdout as string
func runCmd(cmdStr string, args ...string) (string, error) {
	// if args provided use exec.Command with args
	ctx, cancel := context.WithTimeout(context.Background(), CmdTimeout)
	defer cancel()

	var cmd *exec.Cmd
	if len(args) == 0 {
		// use bash -c for compound commands
		cmd = exec.CommandContext(ctx, "bash", "-c", cmdStr)
	} else {
		// direct exec
		cmd = exec.CommandContext(ctx, cmdStr, args...)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// return stderr for debugging if present
		if stderr.Len() > 0 {
			return stdout.String(), fmt.Errorf("cmd error: %w: %s", err, stderr.String())
		}
		return stdout.String(), fmt.Errorf("cmd error: %w", err)
	}
	return stdout.String(), nil
}

// uniqueFilter filters a slice into sorted unique slice
func uniqueFilter(items []string) []string {
	set := make(map[string]struct{}, len(items))
	for _, i := range items {
		i = strings.TrimSpace(i)
		if i == "" {
			continue
		}
		set[i] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// linesFromString splits into lines and trims
func linesFromString(s string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

// writeJSON writes a value to a file path as JSON (pretty)
func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

//
// ----------------- Discover JS files -----------------
//

// DiscoverJSForDomain runs the configured tools to collect JS files for one domain.
// It uses getallurls, waybackurls, katana, and a subfinder|httpx|getallurls pipeline.
// Returns unique, sorted JS file URLs.
func DiscoverJSForDomain(domain string) ([]string, error) {
	var results []string

	log.Printf("[discover] %s: running getallurls --subs", domain)
	if out, err := runCmd(fmt.Sprintf("getallurls --subs %s", domain)); err == nil {
		results = append(results, linesWithJS(out)...)
	} else {
		log.Printf("[discover] getallurls error for %s: %v", domain, err)
	}

	log.Printf("[discover] %s: running waybackurls", domain)
	if out, err := runCmd(fmt.Sprintf("waybackurls %s", domain)); err == nil {
		results = append(results, linesWithJS(out)...)
	} else {
		log.Printf("[discover] waybackurls error for %s: %v", domain, err)
	}

	log.Printf("[discover] %s: running katana", domain)
	// katana usually accepts -u <url>
	if out, err := runCmd(fmt.Sprintf("katana -u https://%s -jc -silent", domain)); err == nil {
		results = append(results, linesWithJS(out)...)
	} else {
		// Katana might be unavailable or require different args; don't fail entire flow
		log.Printf("[discover] katana error for %s: %v", domain, err)
	}

	log.Printf("[discover] %s: running subfinder|httpx|getallurls pipeline", domain)
	pipe := fmt.Sprintf("echo %s | subfinder -silent | httpx -silent | getallurls", domain)
	if out, err := runCmd(pipe); err == nil {
		results = append(results, linesWithJS(out)...)
	} else {
		log.Printf("[discover] pipeline error for %s: %v", domain, err)
	}

	// filter unique js and only those that have .js at end or in path
	unique := uniqueFilter(results)
	// ensure they contain ".js"
	var filtered []string
	for _, u := range unique {
		if strings.Contains(u, ".js") {
			filtered = append(filtered, u)
		}
	}
	return uniqueFilter(filtered), nil
}

func linesWithJS(out string) []string {
	lines := linesFromString(out)
	var filtered []string
	for _, l := range lines {
		if strings.Contains(strings.ToLower(l), ".js") {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

//
// ----------------- Extract Endpoints (LinkFinder) -----------------
//

// ExtractEndpointsFromJS runs LinkFinder on each JS file URL and returns endpoints discovered.
// Uses concurrency with a worker pool to avoid spamming local resources.
func ExtractEndpointsFromJS(jsURLs []string, concurrency int) ([]string, error) {
	if concurrency <= 0 {
		concurrency = 5
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	in := make(chan string, len(jsURLs))
	outList := make([]string, 0, 100)

	// worker function
	worker := func() {
		defer wg.Done()
		for js := range in {
			// LinkFinder invocation - we assume "linkfinder.py" is available with -i <url> -o cli
			cmd := fmt.Sprintf("python3 /opt/exploits/LinkFinder/LinkFinder/linkfinder.py -i %s -o cli", escapeShellArg(js))
			log.Printf("[linkfinder] %s", js)
			out, err := runCmd(cmd)
			if err != nil {
				log.Printf("[linkfinder] error for %s: %v", js, err)
				// continue; LinkFinder may fail for some urls
			}
			// collect lines that look like endpoints (we'll be permissive)
			lines := linesFromString(out)
			mu.Lock()
			outList = append(outList, lines...)
			mu.Unlock()
		}
	}

	// start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	// feed inputs
	for _, u := range jsURLs {
		in <- u
	}
	close(in)
	wg.Wait()

	return uniqueFilter(outList), nil
}

//
// ----------------- Extract Secrets (SecretFinder, gf, grep) -----------------
//

// ExtractSecretsFromJS runs SecretFinder and GF/grep patterns over the provided JS files.
// It returns aggregated secret-like findings.
func ExtractSecretsFromJS(jsURLs []string, concurrency int) (SecretResults, error) {
	if concurrency <= 0 {
		concurrency = 5
	}
	var mu sync.Mutex
	var wg sync.WaitGroup

	in := make(chan string, len(jsURLs))
	secrets := SecretResults{}

	// raw collector
	var rawCollection []string

	worker := func() {
		defer wg.Done()
		for js := range in {
			// SecretFinder
			secCmd := fmt.Sprintf("python3 SecretFinder.py -i %s -o cli", escapeShellArg(js))
			log.Printf("[secretfinder] %s", js)
			out, err := runCmd(secCmd)
			if err == nil && out != "" {
				lines := linesFromString(out)
				mu.Lock()
				rawCollection = append(rawCollection, lines...)
				mu.Unlock()
			} else if err != nil {
				log.Printf("[secretfinder] error for %s: %v", js, err)
			}

			// Run gf patterns by fetching the JS then piping into gf (safer than piping list of URLs)
			// We'll attempt to curl the JS content and run gf on it
			curlCmd := fmt.Sprintf("curl -sL %s", escapeShellArg(js))
			content, err := runCmd(curlCmd)
			if err != nil || strings.TrimSpace(content) == "" {
				// nothing to grep
				continue
			}

			// run GF apikeys
			gfCmd := fmt.Sprintf("echo %s | gf apikeys", escapeSingleQuotes(content))
			if out2, err2 := runCmd(gfCmd); err2 == nil {
				lines := linesFromString(out2)
				mu.Lock()
				rawCollection = append(rawCollection, lines...)
				mu.Unlock()
			}

			// run GF aws-keys
			gfCmd2 := fmt.Sprintf("echo %s | gf aws-keys", escapeSingleQuotes(content))
			if out3, err3 := runCmd(gfCmd2); err3 == nil {
				lines := linesFromString(out3)
				mu.Lock()
				rawCollection = append(rawCollection, lines...)
				mu.Unlock()
			}

			// basic regex grep for common secrets
			grepCmd := fmt.Sprintf("echo %s | grep -Eo '(apiKey|authToken|client_secret|accessToken|aws_secret_access_key)[\"'= ]+[^\"' ]+' || true", escapeSingleQuotes(content))
			if gOut, errg := runCmd(grepCmd); errg == nil {
				lines := linesFromString(gOut)
				mu.Lock()
				rawCollection = append(rawCollection, lines...)
				mu.Unlock()
			}
		}
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	for _, u := range jsURLs {
		in <- u
	}
	close(in)
	wg.Wait()

	// dedupe and classify some patterns (very naive separation)
	rawUnique := uniqueFilter(rawCollection)
	secrets.Raw = rawUnique

	// naive classification
	for _, r := range rawUnique {
		lr := strings.ToLower(r)
		switch {
		case strings.Contains(lr, "aws") || strings.Contains(lr, "aws_secret") || strings.Contains(lr, "awssecret"):
			secrets.AWSKeys = append(secrets.AWSKeys, r)
		case strings.Contains(lr, "apikey") || strings.Contains(lr, "api_key") || strings.Contains(lr, "auth"):
			secrets.APIKeys = append(secrets.APIKeys, r)
		case strings.Contains(lr, "{") || strings.Contains(lr, "json"):
			secrets.JsonLeak = append(secrets.JsonLeak, r)
		case strings.HasPrefix(strings.TrimSpace(r), "http"):
			secrets.URLs = append(secrets.URLs, r)
		default:
			// leave in raw
		}
	}
	// dedupe classified lists
	secrets.APIKeys = uniqueFilter(secrets.APIKeys)
	secrets.AWSKeys = uniqueFilter(secrets.AWSKeys)
	secrets.JsonLeak = uniqueFilter(secrets.JsonLeak)
	secrets.URLs = uniqueFilter(secrets.URLs)

	return secrets, nil
}

//
// ----------------- HTTP Checks (httpx) -----------------
//

// CheckEndpoints sends endpoints through httpx to see which are alive. It returns raw lines.
func CheckEndpointsWithHttpx(endpoints []string) ([]HTTPCheck, error) {
	if len(endpoints) == 0 {
		return nil, nil
	}
	input := strings.Join(endpoints, "\n")
	// Use httpx -silent -status-code -title -tech-detect
	cmd := exec.Command("bash", "-c", "httpx -silent -status-code -title -tech-detect")
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ctx, cancel := context.WithTimeout(context.Background(), CmdTimeout)
	defer cancel()
	cmd = exec.CommandContext(ctx, "bash", "-c", "httpx -silent -status-code -title -tech-detect")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// Even if httpx returns non-zero, capture output
		log.Printf("[httpx] error: %v; stderr: %s", err, stderr.String())
	}
	lines := linesFromString(stdout.String())
	out := make([]HTTPCheck, 0, len(lines))
	for _, l := range lines {
		out = append(out, HTTPCheck{Line: l})
	}
	return out, nil
}

//
// ----------------- FFUF fuzzing (simple wrapper) -----------------
//

// RunFFUF runs ffuf against a target using provided wordlist. It's a lightweight wrapper that writes endpoints
// to a temp file and uses it as -w input if desired. Returns raw output as a single string in FFUFResult.Raw
func RunFFUF(targetBase string, wordlist string, extraArgs string) (FFUFResult, error) {
	// Example: ffuf -u https://target/FUZZ -w wordlist -mc 200
	cmdStr := fmt.Sprintf("ffuf -u %s/FUZZ -w %s %s", targetBase, wordlist, extraArgs)
	log.Printf("[ffuf] running: %s", cmdStr)
	out, err := runCmd(cmdStr)
	return FFUFResult{Raw: out}, err
}

//
// ----------------- Simple IDOR testing -----------------
//

// SimpleIDORTest does a naive IDOR test by replacing any numeric sequence in endpoints with injectedValue
// and performing GET with an invalid Authorization header. It returns a list of IDORResult.
func SimpleIDORTest(endpoints []string, injectedValue string, timeoutSeconds int) ([]IDORResult, error) {
	results := []IDORResult{}
	if injectedValue == "" {
		injectedValue = "1234"
	}
	clTimeout := 10 * time.Second
	if timeoutSeconds > 0 {
		clTimeout = time.Duration(timeoutSeconds) * time.Second
	}

	client := &httpClient{Timeout: clTimeout}

	for _, ep := range endpoints {
		// naive replace: replace sequences of digits with injectedValue
		newEp := replaceDigitsWith(ep, injectedValue)
		if newEp == ep {
			// try simple qsreplace-like behavior: if it has {id} or :id, attempt to append ?id=injected
			if strings.Contains(ep, "{id}") {
				newEp = strings.Replace(ep, "{id}", injectedValue, -1)
			} else if strings.Contains(ep, ":id") {
				newEp = strings.Replace(ep, ":id", injectedValue, -1)
			} else {
				// skip
				continue
			}
		}
		// perform GET with Authorization header invalid using simple internal helper
		status, bodySnippet, err := client.GetWithHeader(newEp, "Authorization", "Bearer invalidtoken")
		if err != nil {
			log.Printf("[idor] request error for %s: %v", newEp, err)
		}
		res := IDORResult{
			Endpoint:      newEp,
			InjectedValue: injectedValue,
			Status:        status,
			BodySnippet:   bodySnippet,
		}
		results = append(results, res)
	}
	return results, nil
}

// helper - replace digits
func replaceDigitsWith(s, rep string) string {
	// simple run: replace all sequences of digits
	out := ""
	var cur bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			// accumulate digit sequence
			cur.WriteByte(c)
		} else {
			if cur.Len() > 0 {
				// flush and replace
				out += rep
				cur.Reset()
			}
			out += string(c)
		}
	}
	if cur.Len() > 0 {
		out += rep
	}
	return out
}

//
// ----------------- Simple http client (used by IDOR test) -----------------
type httpClient struct {
	Timeout time.Duration
}

func (h *httpClient) GetWithHeader(urlStr, headerName, headerVal string) (int, string, error) {
	if urlStr == "" {
		return 0, "", fmt.Errorf("empty url")
	}
	// use curl fallback (in case of issues) via runCmd to avoid bringing net/http complexity for various TLS settings
	// Build curl command: curl -s -o - -w "%{http_code}" -H "Authorization: Bearer invalidtoken" "<url>"
	cmd := fmt.Sprintf("curl -s -S -L -m %d -H %s: %s %s -w '\\nHTTP_STATUS:%{http_code}'", int(h.Timeout.Seconds()), quoteForHeader(headerName), quoteForHeader(headerVal), escapeShellArg(urlStr))
	out, err := runCmd(cmd)
	if err != nil {
		return 0, "", err
	}
	// parse trailing HTTP_STATUS
	parts := strings.Split(out, "HTTP_STATUS:")
	body := strings.Join(parts[:len(parts)-1], "HTTP_STATUS:")
	body = strings.TrimSpace(body)
	status := 0
	if len(parts) > 1 {
		fmt.Sscanf(strings.TrimSpace(parts[len(parts)-1]), "%d", &status)
	}
	snippet := ""
	if len(body) > 300 {
		snippet = body[:300]
	} else {
		snippet = body
	}
	return status, snippet, nil
}

//
// ----------------- Small helpers for shell escaping -----------------
func escapeShellArg(s string) string {
	// naive escape by wrapping in single quotes and replacing existing single quotes
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
func escapeSingleQuotes(s string) string {
	// replace single quotes with safely escaped sequence for echo '...'
	return strings.ReplaceAll(s, "'", "'\"'\"'")
}

func quoteForHeader(s string) string {
	// ensure header values are quoted properly for curl -H "Name: Value"
	return fmt.Sprintf("\"%s\"", s)
}


// WriteJSONSafe writes a ReconData-like structure to path; logs but returns error
func WriteJSONSafe(path string, v any) error {
	if err := writeJSON(path, v); err != nil {
		log.Printf("failed to write json to %s: %v", path, err)
		return err
	}
	return nil
}
