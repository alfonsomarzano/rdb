package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var (
	buildOutput      string
	buildIncludeDrafts bool
	buildCompression string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Create .rdbdata package",
	Long: `Create a .rdbdata ZIP package from the current commit.

Examples:
  rdb build
  rdb build --out my-package.rdbdata
  rdb build --include-drafts --compression deflate`,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
	
	// Local flags
	buildCmd.Flags().StringVar(&buildOutput, "out", "", "output file (default: ./dist/<repo-name>-<branch>-<short-commit>.rdbdata)")
	buildCmd.Flags().BoolVar(&buildIncludeDrafts, "include-drafts", false, "include draft assets")
	buildCmd.Flags().StringVar(&buildCompression, "compression", "store", "compression method (store or deflate)")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Always use current working directory
	repoPath := "."
	
	// Convert to absolute path
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	
	// Safety check: prevent operations in system directories
	if strings.Contains(strings.ToLower(absPath), "c:\\windows\\system32") {
		return fmt.Errorf("cannot operate on RDB repository in system directory: %s", absPath)
	}
	
	// Check if repository exists
	if !repo.IsRepository(absPath) {
		return fmt.Errorf("not an RDB repository: %s", absPath)
	}
	
	// Open repository
	r, err := repo.OpenRepository(absPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}
	
	// Get current branch and commit
	branch, err := r.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	
	commit, err := r.GetCurrentCommit()
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}
	
	// Determine output file
	outputFile := buildOutput
	if outputFile == "" {
		repoName := filepath.Base(r.Path)
		shortCommit := commit[:8]
		outputFile = fmt.Sprintf("./dist/%s-%s-%s.rdbdata", repoName, branch, shortCommit)
	}
	
	// Create output directory if needed
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Create package
	if err := createPackage(r, outputFile, commit, branch, buildIncludeDrafts, buildCompression); err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}
	
	fmt.Printf("Created package: %s\n", outputFile)
	
	return nil
}

// Manifest represents the package manifest
type Manifest struct {
	SchemaVersion string    `json:"schemaVersion"`
	CreatedAt     time.Time `json:"createdAt"`
	Commit        struct {
		ID        string    `json:"id"`
		Author    string    `json:"author"`
		Timestamp time.Time `json:"timestamp"`
		Message   string    `json:"message"`
		Branch    string    `json:"branch"`
	} `json:"commit"`
	Assets []AssetEntry `json:"assets"`
}

// AssetEntry represents an asset in the manifest
type AssetEntry struct {
	Type  string              `json:"type"`
	ID    int                 `json:"id"`
	Name  string              `json:"name,omitempty"`
	Paths []repo.AssetPath    `json:"paths,omitempty"`
	Meta  interface{}         `json:"meta,omitempty"`
	ETag  string              `json:"etag,omitempty"`
}

func createPackage(r *repo.Repository, outputFile, commitHash, branch string, includeDrafts bool, compression string) error {
	// Create ZIP file
	zipFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// Set compression method
	var method uint16
	switch compression {
	case "store":
		method = zip.Store
	case "deflate":
		method = zip.Deflate
	default:
		return fmt.Errorf("invalid compression method: %s", compression)
	}
	
	// Read commit object
	objType, data, err := r.ReadObject(commitHash)
	if err != nil {
		return fmt.Errorf("failed to read commit object: %w", err)
	}
	
	if objType != "commit" {
		return fmt.Errorf("object is not a commit")
	}
	
	var commit repo.Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		return fmt.Errorf("failed to unmarshal commit: %w", err)
	}
	
	// Create manifest
	manifest := &Manifest{
		SchemaVersion: "1.0",
		CreatedAt:     time.Now(),
	}
	
	manifest.Commit.ID = commit.ID
	manifest.Commit.Author = commit.Author
	manifest.Commit.Timestamp = commit.Timestamp
	manifest.Commit.Message = commit.Message
	manifest.Commit.Branch = branch
	
	// TODO: Add assets to manifest
	// For now, create empty manifest
	
	// Write manifest
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	manifestHeader := &zip.FileHeader{
		Name:   "rdb-manifest.json",
		Method: method,
	}
	manifestHeader.SetModTime(time.Now())
	
	manifestFile, err := zipWriter.CreateHeader(manifestHeader)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err)
	}
	
	if _, err := manifestFile.Write(manifestData); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	
	// TODO: Copy objects to package
	// For now, just create the basic structure
	
	return nil
} 