package command

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"golang.design/x/clipboard"

	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/internal/config"
)

func Search(field, text string, cfg config.Config, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	entries, err := getEntriesBy(client, field, text)
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

	if cfg.AutomaticallyOpenBrowser && strings.HasPrefix(selectedEntry.Body, "http") {
		cmd := exec.Command(cfg.Browser, selectedEntry.Body)
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("command execute failed: %v", err)
		}
	}

	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(stdout, "value is '%s'\n", selectedEntry.Body)
		return nil
	}

	clipch := clipboard.Write(clipboard.FmtText, []byte(selectedEntry.Body))
	fmt.Fprintf(stdout, "copied '%s' to clipboard\n", selectedEntry.Body)

	// This is needed to work on Linux.
	// Ref: https://github.com/golang-design/clipboard/issues/15
	select {
	case <-clipch:
	case <-time.After(1 * time.Second):
	}

	return nil
}
