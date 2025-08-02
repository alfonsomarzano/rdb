package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Change directory to asset folder",
	Long: `Change directory to asset folder by searching for asset type names.

Searches through asset descriptions to find matching folders.

Examples:
  rdb cd text        # Go to text-related folders (Strings, Misc Text Files, etc.)
  rdb cd image       # Go to image-related folders (Flash Images, Images, etc.)
  rdb cd sound       # Go to sound-related folders (Sound Effects, Music, etc.)
  rdb cd xml         # Go to XML-related folders (XML Treasure Data, etc.)`,
	RunE: runCd,
}

func init() {
	rootCmd.AddCommand(cdCmd)
}

func runCd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what to search for (e.g., 'rdb cd text')")
	}
	
	searchTerm := strings.ToLower(args[0])
	
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
	
	// Asset type mapping with searchable descriptions
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
	
	// Find matching assets
	var matches []struct {
		id   int
		name string
	}
	
	for id, name := range assetTypes {
		if strings.Contains(strings.ToLower(name), searchTerm) {
			matches = append(matches, struct {
				id   int
				name string
			}{id, name})
		}
	}
	
	if len(matches) == 0 {
		return fmt.Errorf("no asset folders found matching '%s'", searchTerm)
	}
	
	// If multiple matches, check if user provided a selection number
	if len(matches) > 1 {
		if len(args) > 1 {
			// User provided a selection number
			selectionStr := args[1]
			var selection int
			if _, err := fmt.Sscanf(selectionStr, "%d", &selection); err != nil {
				return fmt.Errorf("invalid selection number: %s", selectionStr)
			}
			
			if selection < 1 || selection > len(matches) {
				return fmt.Errorf("selection number must be between 1 and %d", len(matches))
			}
			
			// Use the selected match
			match := matches[selection-1]
			assetPath := filepath.Join(r.Path, "assets", fmt.Sprintf("%d", match.id))
			
			if _, err := os.Stat(assetPath); err != nil {
				return fmt.Errorf("asset folder not found: %s", assetPath)
			}
			
			fmt.Printf("Changing directory to: %s (%s)\n", assetPath, match.name)
			if err := os.Chdir(assetPath); err != nil {
				return fmt.Errorf("failed to change directory: %w", err)
			}
			
			return nil
		}
		
		// Show multiple matches and let user choose
		fmt.Printf("Multiple matches found for '%s':\n", searchTerm)
		for i, match := range matches {
			fmt.Printf("  %d. %07d - %s\n", i+1, match.id, match.name)
		}
		fmt.Printf("\nPlease specify which one (e.g., 'rdb cd text 1' for the first match)\n")
		return nil
	}
	
	// Single match - change to that directory
	match := matches[0]
	assetPath := filepath.Join(r.Path, "assets", fmt.Sprintf("%d", match.id))
	
	if _, err := os.Stat(assetPath); err != nil {
		return fmt.Errorf("asset folder not found: %s", assetPath)
	}
	
	fmt.Printf("Changing directory to: %s (%s)\n", assetPath, match.name)
	if err := os.Chdir(assetPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	
	return nil
} 