package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository("/test/path")
	if repo.Path != "/test/path" {
		t.Errorf("Expected path /test/path, got %s", repo.Path)
	}
	if repo.Config == nil {
		t.Error("Expected config to be initialized")
	}
}

func TestIsRepository(t *testing.T) {
	// Test with non-existent path
	if IsRepository("/non/existent/path") {
		t.Error("Expected false for non-existent path")
	}
	
	// Test with existing directory but no .rdb
	tempDir := t.TempDir()
	if IsRepository(tempDir) {
		t.Error("Expected false for directory without .rdb")
	}
	
	// Test with valid repository
	repo := NewRepository(tempDir)
	if err := repo.Init("tree", []string{"text", "audio"}); err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	
	if !IsRepository(tempDir) {
		t.Error("Expected true for valid repository")
	}
}

func TestRepositoryInit(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewRepository(tempDir)
	
	// Test initialization
	if err := repo.Init("tree", []string{"text", "audio", "texture"}); err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	
	// Check that .rdb directory was created
	rdbPath := filepath.Join(tempDir, ".rdb")
	if _, err := os.Stat(rdbPath); os.IsNotExist(err) {
		t.Error(".rdb directory was not created")
	}
	
	// Check that config was saved
	configPath := filepath.Join(rdbPath, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config.json was not created")
	}
	
	// Check that assets directory was created
	assetsPath := filepath.Join(tempDir, "assets")
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		t.Error("assets directory was not created")
	}
	
	// Check config values
	if repo.Config.Core.Layout != "tree" {
		t.Errorf("Expected layout 'tree', got '%s'", repo.Config.Core.Layout)
	}
	
	if len(repo.Config.Types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(repo.Config.Types))
	}
}

func TestOpenRepository(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test opening non-existent repository
	_, err := OpenRepository(tempDir)
	if err == nil {
		t.Error("Expected error when opening non-existent repository")
	}
	
	// Create repository
	repo := NewRepository(tempDir)
	if err := repo.Init("tree", []string{"text"}); err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	
	// Test opening valid repository
	openedRepo, err := OpenRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to open repository: %v", err)
	}
	
	if openedRepo.Path != tempDir {
		t.Errorf("Expected path %s, got %s", tempDir, openedRepo.Path)
	}
	
	if len(openedRepo.Config.Types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(openedRepo.Config.Types))
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewRepository(tempDir)
	
	if err := repo.Init("tree", []string{"text"}); err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	
	branch, err := repo.GetCurrentBranch()
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}
	
	if branch != "main" {
		t.Errorf("Expected branch 'main', got '%s'", branch)
	}
}

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()
	
	if id1 == id2 {
		t.Error("Generated IDs should be unique")
	}
	
	if len(id1) != 16 {
		t.Errorf("Expected ID length 16, got %d", len(id1))
	}
} 