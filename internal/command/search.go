package command

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/internal/config"
)

func Search(text string, cfg config.Config, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	entries, err := getEntriesByTitle(client, text)
	if err != nil {
		return fmt.Errorf("search failed: %v", err)
	}

	selectedEntry, err := selectEntry(stderr, stdout, cfg.SelectCmd, entries)
	if err != nil {
		return err
	}

	if selectedEntry == nil {
		return nil
	}

	if err := clipboard.WriteAll(selectedEntry.Body); err == nil {
		fmt.Fprintf(stdout, "copied '%s' to clipboard\n", selectedEntry.Body)
	} else {
		fmt.Fprintf(stdout, "value is '%s'\n", selectedEntry.Body)
	}

	if cfg.AutomaticallyOpenBrowser && strings.HasPrefix(selectedEntry.Body, "http") {
		cmd := exec.Command(cfg.Browser, selectedEntry.Body)
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("command execute failed: %v", err)
		}
	}

	return nil
}
