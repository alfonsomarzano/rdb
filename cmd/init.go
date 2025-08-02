package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var (
	layout string
	types  string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new RDB repository",
	Long: `Initialize a new RDB repository in the current directory or at the specified path.

Creates the directory tree and .rdb structure with the specified layout and asset types.

Examples:
  rdb init --layout tree --types "text,audio,texture,shader,mesh"
  rdb init --layout flat --types "text,audio"`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Local flags
	initCmd.Flags().StringVar(&layout, "layout", "tree", "repository layout (tree or flat)")
	initCmd.Flags().StringVar(&types, "types", "text,audio,texture,shader,mesh", "comma-separated list of asset types (optional)")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine repository path - always use current working directory
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}
	
	// Convert to absolute path
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	
	// Safety check: prevent operations in system directories
	if strings.Contains(strings.ToLower(absPath), "c:\\windows\\system32") {
		return fmt.Errorf("cannot create RDB repository in system directory: %s", absPath)
	}
	
	// Check if repository already exists
	if repo.IsRepository(absPath) {
		return fmt.Errorf("RDB repository already exists at %s", absPath)
	}
	
	// Create repository directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}
	
	// Parse asset types
	assetTypes := strings.Split(types, ",")
	for i, t := range assetTypes {
		assetTypes[i] = strings.TrimSpace(t)
	}
	
	// Validate layout
	if layout != "tree" && layout != "flat" {
		return fmt.Errorf("invalid layout: %s (must be 'tree' or 'flat')", layout)
	}
	
	// Create new repository
	r := repo.NewRepository(absPath)
	
	// Initialize repository
	if err := r.Init(layout, assetTypes); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	
	fmt.Printf("Initialized RDB repository at %s\n", absPath)
	fmt.Printf("Layout: %s\n", layout)
	fmt.Printf("Asset types: %s\n", strings.Join(assetTypes, ", "))
	
	return nil
} 