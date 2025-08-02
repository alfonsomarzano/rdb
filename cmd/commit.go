package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var (
	commitMessage string
	commitAuthor  string
	amend         bool
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new commit",
	Long: `Create a new commit with the staged changes.

Examples:
  rdb commit -m "Add intro dialog line"
  rdb commit -m "Update textures" --author "John Doe <john@example.com>"
  rdb commit --amend`,
	RunE: runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	
	// Local flags
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "commit message")
	commitCmd.Flags().StringVar(&commitAuthor, "author", "", "author (format: 'Name <email>')")
	commitCmd.Flags().BoolVar(&amend, "amend", false, "amend the previous commit")
	
	// Mark required flags
	commitCmd.MarkFlagRequired("message")
}

func runCommit(cmd *cobra.Command, args []string) error {
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
	
	// Get current commit (for amend)
	var parentCommit string
	if amend {
		parentCommit, err = r.GetCurrentCommit()
		if err != nil {
			return fmt.Errorf("failed to get current commit: %w", err)
		}
	}
	
	// Determine author
	author := commitAuthor
	if author == "" {
		// TODO: Get from config or environment
		author = "RDB <rdb@localhost>"
	}
	
	// Create commit
	commit := &repo.Commit{
		ID:        repo.GenerateID(),
		Author:    author,
		Timestamp: time.Now(),
		Message:   commitMessage,
		Branch:    branch,
	}
	
	if amend && parentCommit != "" {
		commit.Parent = parentCommit
	}
	
	// TODO: Create tree from staged changes
	// For now, create an empty tree
	tree := &repo.Tree{Entries: []repo.TreeEntry{}}
	treeHash, err := r.WriteObject("tree", tree)
	if err != nil {
		return fmt.Errorf("failed to write tree object: %w", err)
	}
	commit.Tree = treeHash
	
	// Write commit object
	commitHash, err := r.WriteObject("commit", commit)
	if err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}
	
	// Update branch reference
	refPath := filepath.Join(r.Path, ".rdb", "refs", "heads", branch)
	if err := os.WriteFile(refPath, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("failed to update branch reference: %w", err)
	}
	
	if amend {
		fmt.Printf("Amended commit %s\n", commitHash[:8])
	} else {
		fmt.Printf("Created commit %s\n", commitHash[:8])
	}
	
	return nil
} 