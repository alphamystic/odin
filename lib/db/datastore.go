package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alphamystic/odin/lib/utils"
)

type DataStore struct {
	mu          sync.RWMutex            // Controls access to categories and fileLocks maps
	fileLocks   map[string]*sync.Mutex  // Maps folder paths to dedicated Mutexes
	categories  map[string]bool         // Registry of allowed subcategories
	defaultPerm os.FileMode             // Default file creation permissions
}

// NewDataStore initializes the storage driver with an initial list of allowed categories.
func NewDataStore(perm os.FileMode, initialCategories ...string) *DataStore {
	ds := &DataStore{
		fileLocks:   make(map[string]*sync.Mutex),
		categories:  make(map[string]bool),
		defaultPerm: perm,
	}

	// Register any boot-time categories passed in
	for _, cat := range initialCategories {
		if !utils.CheckifStringIsEmpty(cat) {
			ds.categories[cat] = true
		}
	}

	return ds
}

// RegisterCategory dynamically appends a new allowed subcategory at runtime.
func (ds *DataStore) RegisterCategory(category string) {
	if !utils.CheckifStringIsEmpty(category) {
		return
	}
	ds.mu.Lock()
	ds.categories[category] = true
	ds.mu.Unlock()
}

// ListCategories returns a slice containing all currently tracked categories.
func (ds *DataStore) ListCategories() []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	list := make([]string, 0, len(ds.categories))
	for cat := range ds.categories {
		list = append(list, cat)
	}
	return list
}

// IsCategoryAllowed checks if the requested subcategory path is registered.
func (ds *DataStore) IsCategoryAllowed(category string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.categories[category]
}

// Write saves a JSON file into: <baseDir>/<category>/<id>/<fileName>.json
func (ds *DataStore) Write(baseDir, category, id, fileName string, data interface{}) error {
	if !ds.IsCategoryAllowed(category) {
		return fmt.Errorf("unauthorized directory path: category '%s' is not registered", category)
	}
	if err := validatePaths(baseDir, category, id, fileName); err != nil {
		return err
	}

	targetDir := filepath.Clean(filepath.Join(baseDir, category, id))
	finalPath := filepath.Join(targetDir, fileName+".json")
	tempPath := finalPath + ".tmp"

	// Fetch or assign structural lock safely
	lock := ds.getLock(targetDir)
	lock.Lock()
	defer lock.Unlock()

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure %s: %w", targetDir, err)
	}

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}
	b = append(b, byte('\n'))

	if err := os.WriteFile(tempPath, b, ds.defaultPerm); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempPath, finalPath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to commit file via rename: %w", err)
	}

	return nil
}

// Read retrieves and unmarshals a specific JSON file.
func (ds *DataStore) Read(baseDir, category, id, fileName string, v interface{}) error {
	if !ds.IsCategoryAllowed(category) {
		return fmt.Errorf("unauthorized directory path: category '%s' is not registered", category)
	}
	if err := validatePaths(baseDir, category, id, fileName); err != nil {
		return err
	}

	recordPath := filepath.Clean(filepath.Join(baseDir, category, id, fileName+".json"))

	b, err := os.ReadFile(recordPath)
	if err != nil {
		return fmt.Errorf("failed to read record at %s: %w", recordPath, err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return nil
}

// ReadAll grabs all valid file contents found under a specific <baseDir>/<category>/<id> node.
func (ds *DataStore) ReadAll(baseDir, category, id string) ([]string, error) {
	if !ds.IsCategoryAllowed(category) {
		return nil, fmt.Errorf("unauthorized directory path: category '%s' is not registered", category)
	}
	if !utils.CheckifStringIsEmpty(baseDir) || !utils.CheckifStringIsEmpty(category) || !utils.CheckifStringIsEmpty(id) {
		return nil, fmt.Errorf("parameters cannot be empty")
	}

	targetDir := filepath.Clean(filepath.Join(baseDir, category, id))
	files, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection directory: %w", err)
	}

	var records []string
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		b, err := os.ReadFile(filepath.Join(targetDir, file.Name()))
		if err != nil {
			utils.NoticeError(fmt.Sprintf("Failed to read file %s: %v", file.Name(), err))
			continue
		}
		records = append(records, string(b))
	}

	return records, nil
}

// Delete removes either a single document or an entire collection cluster.
func (ds *DataStore) Delete(baseDir, category, id, fileName string) error {
	if !ds.IsCategoryAllowed(category) {
		return fmt.Errorf("unauthorized directory path: category '%s' is not registered", category)
	}
	if !utils.CheckifStringIsEmpty(baseDir) || !utils.CheckifStringIsEmpty(category) || !utils.CheckifStringIsEmpty(id) {
		return fmt.Errorf("base directory, category, and id are required for deletion")
	}

	var targetPath string
	targetDir := filepath.Clean(filepath.Join(baseDir, category, id))

	if fileName == "" {
		targetPath = targetDir
	} else {
		targetPath = filepath.Join(targetDir, fileName+".json")
	}

	lock := ds.getLock(targetDir)
	lock.Lock()
	defer lock.Unlock()

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("target path does not exist: %s", targetPath)
	}

	return os.RemoveAll(targetPath)
}

// Internal helper to fetch or assign unique structural locks
func (ds *DataStore) getLock(pathKey string) *sync.Mutex {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	m, exists := ds.fileLocks[pathKey]
	if !exists {
		m = &sync.Mutex{}
		ds.fileLocks[pathKey] = m
	}
	return m
}

func validatePaths(base, cat, id, file string) error {
	if !utils.CheckifStringIsEmpty(base) || !utils.CheckifStringIsEmpty(cat) ||
	   !utils.CheckifStringIsEmpty(id) || !utils.CheckifStringIsEmpty(file) {
		return fmt.Errorf("invalid operation: paths and identifiers cannot be empty strings")
	}
	return nil
}






//
// package main
//
// import (
// 	"fmt"
// 	"log"
// 	"db"
// )
//
// func main() {
// 	// 1. Initialize with your static required categories
// 	store := db.NewDataStore(0644, "recondata", "enumdata", "temp")
//
// 	// 2. Add subcategories dynamically as your system expands
// 	store.RegisterCategory("osintdata")
//
// 	// 3. Print out tracking info
// 	fmt.Println("Allowed Categories:", store.ListCategories())
//
// 	// 4. Test execution logic
// 	payload := map[string]string{"status": "alive"}
//
// 	// Succeeds because 'recondata' is tracked
// 	err := store.Write("/opt/odin", "recondata", "id-001", "scan", payload)
// 	if err != nil {
// 		log.Printf("Error: %v", err)
// 	}
//
// 	// Fails instantly before attempting disk operations because 'malwaredata' isn't registered
// 	err = store.Write("/opt/odin", "malwaredata", "id-001", "sample", payload)
// 	if err != nil {
// 		fmt.Println("Blocked Write Error:", err)
// 	}
// }