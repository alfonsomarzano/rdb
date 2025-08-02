package repo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Repository represents an RDB repository
type Repository struct {
	Path   string
	Config *Config
}

// Config represents the repository configuration
type Config struct {
	Core struct {
		Layout   string `json:"layout"`   // "tree" or "flat"
		AutoCRLF string `json:"autocrlf"` // "true", "false", or "input"
	} `json:"core"`
	Types []string `json:"types,omitempty"`
}

// Asset represents a typed asset in the repository
type Asset struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	
	// Metadata
	Tags       []string               `json:"tags,omitempty"`
	Version    int                    `json:"version,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Dependencies []Dependency         `json:"dependencies,omitempty"`
	
	// Content
	Paths []AssetPath `json:"paths,omitempty"`
	ETag  string      `json:"etag,omitempty"`
}

// AssetPath represents a logical path to content
type AssetPath struct {
	Logical string `json:"logical"`
	Object  string `json:"object"` // SHA256 hash
	Size    int64  `json:"size"`
}

// Dependency represents an asset dependency
type Dependency struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

// Commit represents a commit in the repository
type Commit struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Branch    string    `json:"branch"`
	Parent    string    `json:"parent,omitempty"`
	Tree      string    `json:"tree"` // SHA256 of the tree object
}

// Tree represents a directory tree
type Tree struct {
	Entries []TreeEntry `json:"entries"`
}

// TreeEntry represents an entry in a tree
type TreeEntry struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "blob", "tree", or "asset"
	Object   string `json:"object"` // SHA256 hash
	Size     int64  `json:"size,omitempty"`
	AssetID  int    `json:"asset_id,omitempty"`
	AssetType string `json:"asset_type,omitempty"`
}

// NewRepository creates a new repository at the given path
func NewRepository(path string) *Repository {
	return &Repository{
		Path: path,
		Config: &Config{},
	}
}

// Init initializes a new RDB repository
func (r *Repository) Init(layout string, types []string) error {
	// Create .rdb directory structure
	rdbPath := filepath.Join(r.Path, ".rdb")
	
	dirs := []string{
		rdbPath,
		filepath.Join(rdbPath, "refs", "heads"),
		filepath.Join(rdbPath, "refs", "tags"),
		filepath.Join(rdbPath, "objects"),
		filepath.Join(rdbPath, "locks"),
		filepath.Join(rdbPath, "hooks"),
		filepath.Join(rdbPath, "remotes"),
		filepath.Join(rdbPath, "packs"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	// Create initial config
	r.Config.Core.Layout = layout
	r.Config.Core.AutoCRLF = "true"
	r.Config.Types = types
	
	if err := r.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	// Create HEAD file pointing to main branch
	headPath := filepath.Join(rdbPath, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main"), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD: %w", err)
	}
	
	// Create assets directory
	assetsPath := filepath.Join(r.Path, "assets")
	if err := os.MkdirAll(assetsPath, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}
	
	// Create asset directories for all predefined asset types
	assetIDs := []int{
		1000624, // Flash Images
		1030002, // Strings
		1010042, // Loading Screens
		1000083, // XML Treasure Data
		1000087, // XML Zone Transition Points
		1000090, // XML Resurrection Points
		1000635, // USM Video Files
		1000636, // Images
		1070003, // Playfields
		1010013, // Maps
		1010210, // (no name specified)
		1010211, // (no name specified)
		1000623, // Misc Text Files
		1066603, // Unknown Textures
		1020001, // (no name specified)
		1020002, // Sound Effects
		1020005, // Music
		1020006, // Sounds - Tones
		1010207, // Particle Effects
		1000010, // File Names Index / FME Files
		1000007, // PhysX XML
		1020003, // Dialog Audio
		1010008, // Miscellaneous Images
	}
	
	for _, assetID := range assetIDs {
		assetDir := filepath.Join(assetsPath, strconv.Itoa(assetID))
		if err := os.MkdirAll(assetDir, 0755); err != nil {
			return fmt.Errorf("failed to create asset directory %d: %w", assetID, err)
		}
	}
	
	// Create initial commit
	if err := r.createInitialCommit(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}
	
	return nil
}

// SaveConfig saves the repository configuration
func (r *Repository) SaveConfig() error {
	configPath := filepath.Join(r.Path, ".rdb", "config.json")
	
	data, err := json.MarshalIndent(r.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	
	return nil
}

// LoadConfig loads the repository configuration
func (r *Repository) LoadConfig() error {
	configPath := filepath.Join(r.Path, ".rdb", "config.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	
	if err := json.Unmarshal(data, r.Config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return nil
}

// createInitialCommit creates the initial commit for the repository
func (r *Repository) createInitialCommit() error {
	// Create empty tree
	tree := &Tree{Entries: []TreeEntry{}}
	treeHash, err := r.writeObject("tree", tree)
	if err != nil {
		return fmt.Errorf("failed to write tree object: %w", err)
	}
	
	// Create initial commit
	commit := &Commit{
		ID:        generateID(),
		Author:    "RDB <rdb@localhost>",
		Timestamp: time.Now(),
		Message:   "Initial commit",
		Branch:    "main",
		Tree:      treeHash,
	}
	
	commitHash, err := r.writeObject("commit", commit)
	if err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}
	
	// Update HEAD to point to main branch
	headPath := filepath.Join(r.Path, ".rdb", "refs", "heads", "main")
	if err := os.WriteFile(headPath, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD: %w", err)
	}
	
	return nil
}

// writeObject writes an object to the repository
func (r *Repository) writeObject(objType string, obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal object: %w", err)
	}
	
	// Calculate SHA256 hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])
	
	// Create object path
	objPath := filepath.Join(r.Path, ".rdb", "objects", hashStr[:2], hashStr[2:])
	if err := os.MkdirAll(filepath.Dir(objPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create object directory: %w", err)
	}
	
	// Write object with type prefix
	content := fmt.Sprintf("%s %d\000", objType, len(data))
	content += string(data)
	
	if err := os.WriteFile(objPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write object: %w", err)
	}
	
	return hashStr, nil
}

// WriteObject writes an object to the repository (public method)
func (r *Repository) WriteObject(objType string, obj interface{}) (string, error) {
	return r.writeObject(objType, obj)
}

// generateID generates a unique ID
func generateID() string {
	hash := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(hash[:])[:16]
}

// GenerateID generates a unique ID (public method)
func GenerateID() string {
	return generateID()
}

// IsRepository checks if the given path is an RDB repository
func IsRepository(path string) bool {
	rdbPath := filepath.Join(path, ".rdb")
	configPath := filepath.Join(rdbPath, "config.json")
	
	_, err := os.Stat(configPath)
	return err == nil
}

// OpenRepository opens an existing repository
func OpenRepository(path string) (*Repository, error) {
	if !IsRepository(path) {
		return nil, fmt.Errorf("not an RDB repository: %s", path)
	}
	
	repo := NewRepository(path)
	if err := repo.LoadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load repository config: %w", err)
	}
	
	return repo, nil
}

// GetCurrentBranch returns the current branch name
func (r *Repository) GetCurrentBranch() (string, error) {
	headPath := filepath.Join(r.Path, ".rdb", "HEAD")
	
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}
	
	head := string(data)
	if len(head) > 5 && head[:5] == "ref: " {
		// Extract branch name from "ref: refs/heads/branch"
		branch := head[16:] // Skip "ref: refs/heads/"
		return branch, nil
	}
	
	return "", fmt.Errorf("HEAD is not pointing to a branch")
}

// GetCurrentCommit returns the current commit hash
func (r *Repository) GetCurrentCommit() (string, error) {
	branch, err := r.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	
	refPath := filepath.Join(r.Path, ".rdb", "refs", "heads", branch)
	
	data, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("failed to read branch ref: %w", err)
	}
	
	return string(data), nil
}

// readObject reads an object from the repository
func (r *Repository) readObject(hash string) (string, []byte, error) {
	objPath := filepath.Join(r.Path, ".rdb", "objects", hash[:2], hash[2:])
	
	file, err := os.Open(objPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open object: %w", err)
	}
	defer file.Close()
	
	// Read the entire file
	data, err := io.ReadAll(file)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read object: %w", err)
	}
	
	// Find the null byte separator
	nullIndex := -1
	for i, b := range data {
		if b == 0 {
			nullIndex = i
			break
		}
	}
	
	if nullIndex == -1 {
		return "", nil, fmt.Errorf("invalid object format: no null separator found")
	}
	
	// Parse header
	headerStr := string(data[:nullIndex])
	var objType string
	var size int
	_, err = fmt.Sscanf(headerStr, "%s %d", &objType, &size)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse object header: %w", err)
	}
	
	// Extract object data
	objData := data[nullIndex+1:]
	
	// Verify size
	if len(objData) != size {
		return "", nil, fmt.Errorf("object size mismatch: expected %d, got %d", size, len(objData))
	}
	
	return objType, objData, nil
}

// ReadObject reads an object from the repository (public method)
func (r *Repository) ReadObject(hash string) (string, []byte, error) {
	return r.readObject(hash)
} 