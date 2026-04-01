package cmd

import (
	"fmt"
	"os"

	"github.com/acrucettanieto/tag/internal/store"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <session-id>",
	Short: "Remove a pinned session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Load()
		if err != nil {
			return fmt.Errorf("load store: %w", err)
		}
		if !s.Remove(args[0]) {
			fmt.Fprintf(os.Stderr, "error: session %q not found\n", args[0])
			os.Exit(1)
		}
		if err := s.Save(); err != nil {
			return fmt.Errorf("save store: %w", err)
		}
		return nil
	},
}
