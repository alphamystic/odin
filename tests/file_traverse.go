package main

import (
	"os"
	"fmt"
	"io"
	"io/fs"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"
)

this (FileData) wiil be defined in a separete package known as definers imported as dfn, utils will be imported for timestamp here is how utils works type TimeStamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *TimeStamps) Touch() {
  currentTime := time.Now()
  formattedTime := currentTime.Format("2006-01-02 15:04:05")
  parsedTime, _ := time.Parse("2006-01-02 15:04:05", formattedTime)
  t.UpdatedAt = parsedTime.UTC()
  if t.CreatedAt.IsZero() {
    t.CreatedAt = t.UpdatedAt
  }
}
type FileData struct {
	InitialHAsh string Can be a hash or a uuid to identify the file with multiple times ieven after it's hash changes
	Name     string `json:"name"`
	CurrentHash     string `json:"hash"`
	PreviuosHash string
	PreviousHashes string // we can take all previous hashes and turn them t a string of cooma separted values that ay whenever the file hash changes we update it
	Directory string `json:"directory"`
	IsBinary bool//bool can be a binary/executable or just a regular fil
	IsScanned bool
	TimeFileCreated
	TimeFileUpdated
	utils.TimeStamps // the time the file was created and updated in the db
}

// Check if the file has an executable binary extension
func isExecutableBinary(filePath string) bool {
	extensions := []string{".so", ".run", ".exe", ".msc", ".dll"}// add a check for Zip Files too
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, validExt := range extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func traverseFileSystem(rootDir string) ([]FileData, error) {
	var files []FileData

	err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isExecutableBinary(path) {
			hash, err := calculateHash(path)
			if err != nil {
				return fmt.Errorf("error calculating hash for file %s: %v", path, err)
			}
			files = append(files, FileData{
				Name:     info.Name(),
				Hash:     hash,
				Directory: filepath.Dir(path),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func writeToFile(outputPath string, data []FileData) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

Add a new function to scan it using a yara rule and retuurn the output too.
You can also add a function to check if a file has strings basically an ip or a url then return them and the associated file
I intend to use this for threat hunting to look for alicous binaries
func main() {
	rootDir := "/home/sam/Documents/3l0racle/odin" // Adjust this to the root directory you want to scan
	outputPath := "binaries.json"

	files, err := traverseFileSystem(rootDir)
	if err != nil {
		fmt.Printf("Error traversing file system: %v\n", err)
		return
	}

	if err := writeToFile(outputPath, files); err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

	fmt.Printf("Binary files and hashes written to %s\n", outputPath)
}
