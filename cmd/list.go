package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var (
	listCd bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List asset types and folders",
	Long: `List all asset types with their corresponding folders.

Shows the mapping between asset IDs and their types, and optionally changes directory to a specific asset folder.

Examples:
  rdb list                    # List all asset types and folders
  rdb list --cd 1030002      # List and change to Strings folder
  rdb list --cd 1000624      # List and change to Flash Images folder`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	
	// Local flags
	listCmd.Flags().BoolVar(&listCd, "cd", false, "change directory to specified asset folder")
}

func runList(cmd *cobra.Command, args []string) error {
	// Always use current working directory
	repoPath := "."
	
	// Convert to absolute path
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
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
	
	// Asset type mapping
	assetTypes := map[int]string{
		1000624: "Flash Images",
		1030002: "Strings",
		1010042: "Loading Screens",
		1000083: "XML Treasure Data",
		1000087: "XML Zone Transition Points",
		1000090: "XML Resurrection Points",
		1000635: "USM Video Files",
		1000636: "Images",
		1070003: "Playfields",
		1010013: "Maps",
		1010210: "Image (no name)",
		1010211: "Image (no name)",
		1000623: "Misc Text Files",
		1066603: "Unknown Textures",
		1020001: "Unknown",
		1020002: "Sound Effects",
		1020005: "Music",
		1020006: "Sounds - Tones",
		1010207: "Particle Effects",
		1000010: "File Names Index / FME Files",
		1000007: "PhysX XML",
		1020003: "Dialog Audio",
		1010008: "Miscellaneous Images",
	}
	
	// Check if user wants to change directory
	if listCd && len(args) > 0 {
		targetID := args[0]
		if assetName, exists := assetTypes[parseAssetID(targetID)]; exists {
			assetPath := filepath.Join(r.Path, "assets", targetID)
			if _, err := os.Stat(assetPath); err == nil {
				fmt.Printf("Changing directory to: %s (%s)\n", assetPath, assetName)
				if err := os.Chdir(assetPath); err != nil {
					return fmt.Errorf("failed to change directory: %w", err)
				}
				return nil
			} else {
				return fmt.Errorf("asset folder not found: %s", targetID)
			}
		} else {
			return fmt.Errorf("unknown asset ID: %s", targetID)
		}
	}
	
	// List all asset types and folders
	fmt.Printf("Asset Types and Folders:\n")
	fmt.Printf("========================\n\n")
	
	assetsPath := filepath.Join(r.Path, "assets")
	
	for id, name := range assetTypes {
		folderPath := filepath.Join(assetsPath, fmt.Sprintf("%d", id))
		
		// Check if folder exists
		exists := ""
		if _, err := os.Stat(folderPath); err == nil {
			exists = "✓"
		} else {
			exists = "✗"
		}
		
		fmt.Printf("%s %07d - %s\n", exists, id, name)
	}
	
	fmt.Printf("\nUsage:\n")
	fmt.Printf("  rdb list --cd <id>    # Change to specific asset folder\n")
	fmt.Printf("  rdb list --cd 1030002 # Change to Strings folder\n")
	
	return nil
}

func parseAssetID(idStr string) int {
	// Try to parse as integer
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
		return id
	}
	return 0
} 