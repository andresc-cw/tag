package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/acrucettanieto/tag/internal/store"
	"github.com/spf13/cobra"
)

var rawOutput bool

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Pick a tagged session with fzf and print the resume command",
	Long: `Opens an fzf picker showing all tagged sessions.
Press Enter to resume a session (cd + claude --resume).
Press ctrl-d to delete the selected tag without leaving the picker.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Load()
		if err != nil {
			return fmt.Errorf("load store: %w", err)
		}
		if len(s.Sessions) == 0 {
			fmt.Fprintln(os.Stderr, "no tagged sessions")
			return nil
		}

		// --raw: just print lines (used by fzf reload binding)
		if rawOutput {
			for _, p := range s.Sessions {
				fmt.Print(formatLine(p))
			}
			return nil
		}

		if _, err := exec.LookPath("fzf"); err != nil {
			return fmt.Errorf("fzf not found in PATH — install with: brew install fzf")
		}

		// Resolve absolute path so fzf's minimal shell can find the binary.
		tagBin, err := exec.LookPath("tag")
		if err != nil {
			return fmt.Errorf("tag binary not found in PATH: %w", err)
		}

		var buf bytes.Buffer
		for _, p := range s.Sessions {
			fmt.Fprint(&buf, formatLine(p))
		}

		fzf := exec.Command("fzf",
			"--with-nth=2..",
			"--delimiter=\t",
			"--height=40%",
			"--border",
			"--prompt=tag> ",
			"--reverse",
			"--header=enter:resume  ctrl-d:delete",
			fmt.Sprintf("--bind=ctrl-d:execute-silent(%s rm {1})+reload(%s ls --raw)", tagBin, tagBin),
		)
		fzf.Stdin = &buf
		fzf.Stderr = os.Stderr

		out, err := fzf.Output()
		if err != nil {
			if exit, ok := err.(*exec.ExitError); ok {
				// 130 = Escape pressed, 1 = no match / list emptied
				if exit.ExitCode() == 130 || exit.ExitCode() == 1 {
					return nil
				}
			}
			return fmt.Errorf("fzf: %w", err)
		}

		line := strings.TrimSpace(string(out))
		if line == "" {
			return nil
		}
		selectedID := strings.SplitN(line, "\t", 2)[0]

		var project string
		for _, p := range s.Sessions {
			if p.ID == selectedID {
				project = p.Project
				break
			}
		}

		fmt.Printf("cd %s && claude --resume %s\n", project, selectedID)
		return nil
	},
}

func formatLine(p store.Pin) string {
	shortID := p.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	base := filepath.Base(p.Project)
	date := p.PinnedAt.Format("2006-01-02")
	return fmt.Sprintf("%s\t%s  %-20s  %s  %s\n", p.ID, shortID, base, p.Description, date)
}

func init() {
	lsCmd.Flags().BoolVar(&rawOutput, "raw", false, "Print lines without launching fzf (used internally by reload binding)")
}
