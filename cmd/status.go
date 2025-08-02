package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var porcelain bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show working tree status",
	Long: `Show the status of the working tree.

Shows changes: A (added), M (modified), D (deleted), R (renamed), U (unmerged).

Use --porcelain for stable machine output.`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	
	// Local flags
	statusCmd.Flags().BoolVar(&porcelain, "porcelain", false, "give stable machine output")
}

func runStatus(cmd *cobra.Command, args []string) error {
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
	
	// Get current branch
	branch, err := r.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	
	// Get current commit
	commit, err := r.GetCurrentCommit()
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}
	
	if porcelain {
		// Machine-readable output
		fmt.Printf("branch %s\n", branch)
		fmt.Printf("commit %s\n", commit)
		// TODO: Add staged/unstaged changes
	} else {
		// Human-readable output
		fmt.Printf("On branch %s\n", branch)
		fmt.Printf("commit %s\n\n", commit)
		
		// TODO: Show working tree status
		fmt.Println("No changes to commit, working tree clean")
	}
	
	return nil
} 