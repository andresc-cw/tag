package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/acrucettanieto/tag/internal/store"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <session-id> <project-path> <description>",
	Short: "Pin a Claude Code session",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Load()
		if err != nil {
			return fmt.Errorf("load store: %w", err)
		}
		s.Upsert(store.Pin{
			ID:          args[0],
			Project:     args[1],
			Description: args[2],
			PinnedAt:    time.Now().UTC(),
		})
		if err := s.Save(); err != nil {
			return fmt.Errorf("save store: %w", err)
		}
		fmt.Fprintf(os.Stderr, "pinned %s\n", args[0][:min(8, len(args[0]))])
		return nil
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
