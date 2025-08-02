package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/rdb/cli/internal/repo"
	"github.com/spf13/cobra"
)

var (
	logOneline   bool
	logMaxCount  int
	logSince     string
	logUntil     string
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history",
	Long: `Show the commit history.

Examples:
  rdb log
  rdb log --oneline
  rdb log --max-count 10
  rdb log --since "2024-01-01"`,
	RunE: runLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
	
	// Local flags
	logCmd.Flags().BoolVar(&logOneline, "oneline", false, "abbreviate commit output")
	logCmd.Flags().IntVar(&logMaxCount, "max-count", 0, "limit number of commits")
	logCmd.Flags().StringVar(&logSince, "since", "", "show commits more recent than date")
	logCmd.Flags().StringVar(&logUntil, "until", "", "show commits older than date")
}

func runLog(cmd *cobra.Command, args []string) error {
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
	
	// Get current commit
	currentCommit, err := r.GetCurrentCommit()
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}
	
	// Parse date filters
	var sinceTime, untilTime time.Time
	if logSince != "" {
		sinceTime, err = time.Parse("2006-01-02", logSince)
		if err != nil {
			return fmt.Errorf("invalid since date format: %w", err)
		}
	}
	if logUntil != "" {
		untilTime, err = time.Parse("2006-01-02", logUntil)
		if err != nil {
			return fmt.Errorf("invalid until date format: %w", err)
		}
	}
	
	// Show commit history
	if err := showCommitHistory(r, currentCommit, logOneline, logMaxCount, sinceTime, untilTime); err != nil {
		return fmt.Errorf("failed to show commit history: %w", err)
	}
	
	return nil
}

func showCommitHistory(r *repo.Repository, startCommit string, oneline bool, maxCount int, since, until time.Time) error {
	// TODO: Implement proper commit history traversal
	// For now, just show the current commit
	
	objType, data, err := r.ReadObject(startCommit)
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
	
	// Apply filters
	if !since.IsZero() && commit.Timestamp.Before(since) {
		return nil
	}
	if !until.IsZero() && commit.Timestamp.After(until) {
		return nil
	}
	
	// Format output
	if oneline {
		fmt.Printf("%s %s\n", commit.ID[:8], commit.Message)
	} else {
		fmt.Printf("commit %s\n", commit.ID)
		fmt.Printf("Author: %s\n", commit.Author)
		fmt.Printf("Date:   %s\n", commit.Timestamp.Format(time.RFC3339))
		fmt.Printf("\n    %s\n\n", commit.Message)
	}
	
	return nil
} 