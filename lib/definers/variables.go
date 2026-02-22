package definers

import (
	"os"
	"log"
	"fmt"
	"bufio"
	"strings"
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Initialize database connection for the given domain
type DBConfig struct {
	Username string
	Password string
	DBName   string
	Host     string
	Port     int
}

// Initiate a new MysqlDB Connector
func NewMySQLConnector(dbconfig *DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			dbconfig.Username,
			dbconfig.Password,
			dbconfig.Host,
			dbconfig.Port,
			dbconfig.DBName,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating new Connector: %v", err)
	}
	return db, nil
}

// Initialize a db configurator
func InitializeConnector(mode string) *DBConfig {
	envMap := loadEnvFile(".env")

	cfg := &DBConfig{}
	var portStr string

	if strings.ToUpper(mode) == "PROD" {
		cfg.Username = envMap["DB_USER"]
		cfg.Password = envMap["DB_PASS_"]
		cfg.DBName = envMap["DB_NAME"]
		cfg.Host = envMap["DB_HOST"]
		portStr = envMap["DB_PORT"]
	} else {
		cfg.Username = envMap["DB_USER_LOCAL"]
		cfg.Password = envMap["DB_PASS_LOCAL"]
		cfg.DBName = envMap["DB_NAME_LOCAL"]
		cfg.Host = envMap["DB_HOST_LOCAL"]
		portStr = envMap["DB_PORT_LOCAL"]
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number in .env: %v", err)
	}
	cfg.Port = port

	return cfg
}

// Helper to read .env file without external packages
func loadEnvFile(path string) map[string]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	defer file.Close()

	envs := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		envs[key] = value
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	return envs
}
