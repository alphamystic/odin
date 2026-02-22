package loadbalancer

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"sync"
	"context"
	"strings"
	"strconv"
	"syscall"
	"net/url"
	"net/http"
	"os/signal"
	"crypto/tls"
	"encoding/json"
	"net/http/httputil"

	"github.com/alphamystic/odin/lib/utils"
)

type Backend interface {
	SetAlive(bool)
	IsAlive() bool
	GetURL() *url.URL
	GetActiveConnections() int
	IsAvailable() bool
	CheckHealth() bool
	Serve(http.ResponseWriter, *http.Request)
}

type backend struct {
	url            *url.URL
	alive          bool
	mux            sync.RWMutex
	connections    int
	maxConnections int
	reverseProxy   *httputil.ReverseProxy
	name           string
	HealthPath     string
}

func (b *backend) GetActiveConnections() int {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.connections
}

func (b *backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}

func (b *backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) IsAvailable() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive && b.connections < b.maxConnections
}

//	func (b *backend) CheckHealth() bool {
//		resp, err := http.Get(b.url.String() + "/health")
//		if err != nil || resp.StatusCode != http.StatusOK {
//			b.SetAlive(false)
//			return false
//		}
//		b.SetAlive(true)
//		return true
//	}
func (b *backend) CheckHealth() bool {
	u := *b.url
	u.Path = b.HealthPath

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("Health check request creation failed for %s: %v", u.String(), err)
		b.SetAlive(false)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Dev only
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Health check failed for %s: %v", u.String(), err)
		b.SetAlive(false)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unhealthy backend %s: Status %d", u.String(), resp.StatusCode)
		b.SetAlive(false)
		return false
	}

	log.Printf("Health response from %s: %s", u.String(), resp.Status)
	b.SetAlive(true)
	return true
}
func (lb *LBHolder) SetHealthPath(backendName string, path string) {
	for _, b := range lb.backends {
		if b.name == backendName {
			b.HealthPath = path
			return
		}
	}
}

func getHTTPClient(scheme string) *http.Client {
	fmt.Println("")
	fmt.Println(scheme)

	if scheme == "https" {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // ⚠️ dev use only
				},
			},
			Timeout: 3 * time.Second,
		}
	}

	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

func checkHTTP(url string, b *backend) bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	return doHealthCheck(client, url, b)
}

func checkHTTPS(url string, b *backend) bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ⚠️ Use proper certs in prod
			},
		},
	}
	return doHealthCheck(client, url, b)
}

func doHealthCheck(client *http.Client, url string, b *backend) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create health check request for %s: %v", url, err)
		b.SetAlive(false)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Health check failed for %s: %v", url, err)
		b.SetAlive(false)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response from %s: %v", url, err)
		b.SetAlive(false)
		return false
	}

	log.Printf("Health response from %s: %s", url, string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unhealthy backend %s: Status %d", url, resp.StatusCode)
		b.SetAlive(false)
		return false
	}

	b.SetAlive(true)
	return true
}

func (b *backend) Serve(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Routing request %s %s to backend %s", req.Method, req.URL.Path, b.GetURL())

	b.mux.Lock()
	b.connections++
	b.mux.Unlock()

	defer func() {
		b.mux.Lock()
		b.connections--
		b.mux.Unlock()
	}()
	b.reverseProxy.ServeHTTP(rw, req)
}

func NewBackend(rawurl, name string, maxConn int) Backend {
	u, _ := url.Parse(rawurl)
	rp := httputil.NewSingleHostReverseProxy(u)

	// Fix TLS verification issue for reverse proxying to https
	if u.Scheme == "https" {
		rp.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ⚠️ dev only
			},
		}
	}

	return &backend{
		url:            u,
		alive:          true,
		reverseProxy:   rp,
		maxConnections: maxConn,
		name:           name,
	}
}

type Pool struct {
	Backends []Backend
	current  int
	mux      sync.Mutex
	label    string
}

func NewPool(label string) *Pool {
	return &Pool{label: label}
}

func (p *Pool) Add(b Backend) {
	p.mux.Lock()
	p.Backends = append(p.Backends, b)
	p.mux.Unlock()
	go func() {
		for {
			if !b.CheckHealth() {
				log.Printf("[HEALTH] %s is down", b.GetURL())
			}
			time.Sleep(60 * 30 * time.Second)
		}
	}()
}

func (p *Pool) NextAvailable() Backend {
	p.mux.Lock()
	defer p.mux.Unlock()
	n := len(p.Backends)
	start := p.current

	for i := 0; i < n; i++ {
		idx := (start + i) % n
		b := p.Backends[idx]
		if b.IsAvailable() {
			p.current = (idx + 1) % n
			return b
		}
	}

	log.Printf("[SCALE TRIGGER] %s: all backends overloaded or down", p.label)
	return nil
}

func (p *Pool) Serve(w http.ResponseWriter, r *http.Request) {
	b := p.NextAvailable()
	if b == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}
	b.Serve(w, r)
}

type LBHolder struct {
	APIPool  *Pool
	UIPool   *Pool
	backends map[string]*backend
	HTTPSAddr int
	HTTPAddr int
}

func NewLBHolder() *LBHolder {
	envs := utils.LoadEnvFile(".env")
	addr1 := envs["HTTPSADDR"]
	addr2 := envs["HTTPADDR"]
	httpsAddr, err := strconv.Atoi(addr1)
	if err != nil {
		log.Fatalf("Invalid HTTPS Port address")
	}
	httpAddr, err := strconv.Atoi(addr2)
	if err != nil {
		log.Fatalf("Invalid HTTP Port address")
	}
	return &LBHolder{
		APIPool:  NewPool("API"),
		UIPool:   NewPool("UI"),
		backends: make(map[string]*backend),
		HTTPSAddr: httpsAddr,
		HTTPAddr: httpAddr,
	}
}

func (lb *LBHolder) AddBackend(poolType, rawurl, name string, maxConn int) {
	b := NewBackend(rawurl, name, maxConn)

	// Convert from interface to concrete type
	realBackend, ok := b.(*backend)
	if !ok {
		log.Printf("Failed to cast backend %s to concrete type", name)
		return
	}

	// Set health path dynamically based on pool type
	switch poolType {
	case "API":
		realBackend.HealthPath = "/api/health"
		lb.APIPool.Add(realBackend)
	case "UI":
		realBackend.HealthPath = "/health"
		lb.UIPool.Add(realBackend)
	default:
		log.Printf("Unknown pool type: %s", poolType)
		return
	}

	// Store backend in central registry
	lb.backends[name] = realBackend
}

func isAPIPath(path string) bool {
	// Simple rule: All API endpoints are prefixed with /api
	return len(path) >= 5 && path[:5] == "/api"
}

func (lb *LBHolder) StartServer() {
	// Shared handler registration
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/health" || path == "/api/health" {
			// Try API first
			if b := lb.APIPool.NextAvailable(); b != nil {
				b.Serve(w, r)
				return
			}
			if b := lb.UIPool.NextAvailable(); b != nil {
				b.Serve(w, r)
				return
			}
			http.Error(w, "No healthy backends", http.StatusServiceUnavailable)
			return
		}

		// Route logic
		if strings.HasPrefix(path, "/api/") {
			if strings.HasPrefix(path, "/api/auth") {
				lb.RouteTo("UI", w, r)
			} else {
				lb.RouteTo("API", w, r)
			}
		} else {
			lb.RouteTo("UI", w, r)
		}
	})

	http.HandleFunc("/register-backend", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Type    string `json:"type"`
			URL     string `json:"url"`
			Name    string `json:"name"`
			MaxConn int    `json:"max_connections"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		lb.AddBackend(payload.Type, payload.URL, payload.Name, payload.MaxConn)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Backend added"))
	})

	// --- Two servers: HTTP and HTTPS ---
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d",lb.HTTPAddr),  // Plain HTTP
	}

	httpsServer := &http.Server{
		Addr: fmt.Sprintf(":%d",lb.HTTPSAddr), // HTTPS
	}

	// Start HTTP listener
	go func() {
		log.Println("Load balancer (HTTP) started on :",fmt.Sprintf(":%d",lb.HTTPAddr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS listener
	go func() {
		log.Println("Load balancer (HTTPS) started on :",fmt.Sprintf(":%d",lb.HTTPSAddr))
		if err := httpsServer.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown failed: %v", err)
	}
	if err := httpsServer.Shutdown(ctx); err != nil {
		log.Printf("HTTPS server shutdown failed: %v", err)
	}

	log.Println("Load balancer shut down cleanly.")
}



func (lb *LBHolder) RouteTo(poolType string, w http.ResponseWriter, r *http.Request) {
	switch poolType {
	case "API":
		lb.APIPool.Serve(w, r)
	case "UI":
		lb.UIPool.Serve(w, r)
	default:
		http.Error(w, "Unknown route type", http.StatusBadRequest)
	}
}
