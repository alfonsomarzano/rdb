package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var (
	addType string
	addID   int
	addName string
)

// Asset type mapping based on ID
var assetTypeMap = map[int]string{
	1000624: "flash_image",
	1030002: "string",
	1010042: "loading_screen",
	1000083: "xml_treasure",
	1000087: "xml_zone_transition",
	1000090: "xml_resurrection",
	1000635: "usm_video",
	1000636: "image",
	1070003: "playfield",
	1010013: "map",
	1010210: "image",
	1010211: "image",
	1000623: "text",
	1066603: "texture",
	1020001: "unknown",
	1020002: "sound_effect",
	1020005: "music",
	1020006: "sound_tone",
	1010207: "particle_effect",
	1000010: "file_index",
	1000007: "physx_xml",
	1020003: "dialog_audio",
	1010008: "misc_image",
}

// getAssetTypeFromID returns the asset type for a given ID
func getAssetTypeFromID(id int) string {
	if assetType, exists := assetTypeMap[id]; exists {
		return assetType
	}
	return "unknown"
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Stage files for commit",
	Long: `Stage new/changed files or folders for commit.

If metadata is missing, can create meta.json with the specified type, id, and name.
Asset type is automatically determined from the folder ID if not specified.

Examples:
  rdb add .\assets\1030002\ --id 1030002 --name "DialogLine_Intro"
  rdb add .\assets\42001\music.mp3 --id 42001`,
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	
	// Local flags
	addCmd.Flags().StringVar(&addType, "type", "", "asset type (optional, auto-determined from ID)")
	addCmd.Flags().IntVar(&addID, "id", 0, "asset ID")
	addCmd.Flags().StringVar(&addName, "name", "", "asset name")
}

func runAdd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no paths specified")
	}
	
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
	
	// Process each path, expanding glob patterns
	var allPaths []string
	for _, pattern := range args {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("invalid glob pattern %s: %w", pattern, err)
		}
		if len(matches) == 0 {
			fmt.Printf("Warning: no files match pattern %s\n", pattern)
			continue
		}
		
		// Filter out non-asset paths
		for _, match := range matches {
			absMatch, err := filepath.Abs(match)
			if err != nil {
				fmt.Printf("Warning: could not resolve path %s: %v\n", match, err)
				continue
			}
			if isAssetPath(r.Path, absMatch) {
				allPaths = append(allPaths, match)
			} else {
				fmt.Printf("Skipping non-asset path: %s\n", match)
			}
		}
	}
	
	// Process each resolved path
	if len(allPaths) == 0 {
		fmt.Println("No asset files found to add")
		return nil
	}
	
	for _, path := range allPaths {
		if err := addPath(r, path, addType, addID, addName); err != nil {
			return fmt.Errorf("failed to add %s: %w", path, err)
		}
	}
	
	return nil
}

func addPath(r *repo.Repository, path, assetType string, assetID int, assetName string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	
	// Check if path exists
	_, err = os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}
	
	// Determine asset type and ID from path if not specified
	if assetType == "" || assetID == 0 {
		relPath, err := filepath.Rel(r.Path, absPath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		
		// Parse path to extract type and ID
		// Expected format: assets/<id>/...
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) >= 2 && parts[0] == "assets" {
			if assetID == 0 {
				if id, err := strconv.Atoi(parts[1]); err == nil {
					assetID = id
				}
			}
		}
		
		// If we still don't have type/ID, try to find the containing asset directory
		if assetType == "" || assetID == 0 {
			assetDir := findAssetDirectory(r.Path, absPath)
			if assetDir != "" {
				relAssetPath, _ := filepath.Rel(r.Path, assetDir)
				assetParts := strings.Split(relAssetPath, string(filepath.Separator))
				if len(assetParts) >= 2 && assetParts[0] == "assets" {
					if assetID == 0 {
						if id, err := strconv.Atoi(assetParts[1]); err == nil {
							assetID = id
						}
					}
				}
			}
		}
	}
	
	// Automatically determine asset type from ID if not specified
	if assetType == "" && assetID != 0 {
		assetType = getAssetTypeFromID(assetID)
	}
	
	// Validate asset type
	if assetType == "" {
		return fmt.Errorf("asset type not specified and could not be determined from path")
	}
	
	// Create or update metadata
	if err := createOrUpdateMetadata(r, absPath, assetType, assetID, assetName); err != nil {
		return fmt.Errorf("failed to create/update metadata: %w", err)
	}
	
	fmt.Printf("Added %s (type: %s, id: %d)\n", path, assetType, assetID)
	
	return nil
}

func createOrUpdateMetadata(r *repo.Repository, path, assetType string, assetID int, assetName string) error {
	// Determine the asset directory
	_, err := filepath.Rel(r.Path, path)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}
	
	// If path is a file, find its containing asset directory
	_, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}
	
	var assetDir string
	// For now, assume it's a directory
	assetDir = path
	
	// Find the asset directory containing this file
	// Look for a directory with meta.json or create one
	dir := filepath.Dir(path)
	for {
		if dir == r.Path {
			break
		}
		
		metaPath := filepath.Join(dir, "meta.json")
		if _, err := os.Stat(metaPath); err == nil {
			assetDir = dir
			break
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	if assetDir == "" {
		// Create asset directory based on ID
		assetDir = filepath.Join(r.Path, "assets", strconv.Itoa(assetID))
		if err := os.MkdirAll(assetDir, 0755); err != nil {
			return fmt.Errorf("failed to create asset directory: %w", err)
		}
	}
	
	// Create or update meta.json
	metaPath := filepath.Join(assetDir, "meta.json")
	
	// TODO: Implement metadata creation/update logic
	// For now, just create a basic metadata file
	metadata := map[string]interface{}{
		"type": assetType,
		"id":   assetID,
	}
	
	if assetName != "" {
		metadata["name"] = assetName
	}
	
	// TODO: Write metadata to file
	fmt.Printf("Would create/update metadata at %s\n", metaPath)
	
	return nil
}

// findAssetDirectory finds the containing asset directory for a given path
func findAssetDirectory(repoPath, filePath string) string {
	dir := filepath.Dir(filePath)
	for {
		if dir == repoPath {
			break
		}
		
		// Check if this directory has a meta.json file
		metaPath := filepath.Join(dir, "meta.json")
		if _, err := os.Stat(metaPath); err == nil {
			return dir
		}
		
		// Check if this directory follows the assets/<id> pattern
		relPath, err := filepath.Rel(repoPath, dir)
		if err == nil {
			parts := strings.Split(relPath, string(filepath.Separator))
			if len(parts) >= 2 && parts[0] == "assets" {
				if _, err := strconv.Atoi(parts[1]); err == nil {
					return dir
				}
			}
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	return ""
}

// isAssetPath checks if a path is within the assets directory structure
func isAssetPath(repoPath, filePath string) bool {
	relPath, err := filepath.Rel(repoPath, filePath)
	if err != nil {
		return false
	}
	
	parts := strings.Split(relPath, string(filepath.Separator))
	
	if len(parts) < 2 {
		return false
	}
	
	// Must be under assets/<id>/...
	if parts[0] != "assets" {
		return false
	}
	
	// Check if the second part is a numeric ID
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return false
	}
	
	return true
} 