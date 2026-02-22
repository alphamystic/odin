package main

import (
  "github.com/alphamystic/odin/lib/loadbalancer"
)

func main() {
	lb := loadbalancer.NewLBHolder()

	lb.AddBackend("API", "http://localhost:5000", "API-1", 100)
	lb.AddBackend("API", "https://localhost:5001", "API-2", 100)
	lb.AddBackend("UI", "http://localhost:4000", "UI-1", 100)
	lb.AddBackend("UI", "https://localhost:4001", "UI-2", 100)

	lb.StartServer()
}
