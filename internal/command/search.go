package command

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.design/x/clipboard"

	"github.com/y-yagi/doco/internal/config"
)

type SearchCommand struct {
	Command
	field  string
	text   string
	cfg    config.Config
	stdout io.Writer
	stderr io.Writer
}

func Search(field, text string, cfg config.Config, stdout, stderr io.Writer) *SearchCommand {
	return &SearchCommand{field: field, text: text, cfg: cfg, stdout: stdout, stderr: stderr}
}

func (c *SearchCommand) Run() error {
	client, err := getEntClient(c.cfg.DataBase)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntriesBy(client, c.field, c.text)
	if err != nil {
		return fmt.Errorf("search failed: %v", err)
	}

	selectedEntry, err := selectEntry(c.stderr, c.stdout, c.cfg.SelectCmd, entries)
	if err != nil {
		return err
	}

	if selectedEntry == nil {
		return nil
	}

	if c.cfg.AutomaticallyOpenBrowser && strings.HasPrefix(selectedEntry.Body, "http") {
		cmd := exec.Command(c.cfg.Browser, selectedEntry.Body)
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("command execute failed: %v", err)
		}
		fmt.Fprintf(c.stdout, "Open '%s'\n", selectedEntry.Body)
		return nil
	}

	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(c.stdout, "value is '%s'\n", selectedEntry.Body)
		return nil
	}

	clipch := clipboard.Write(clipboard.FmtText, []byte(selectedEntry.Body))
	fmt.Fprintf(c.stdout, "copied '%s' to clipboard\n", selectedEntry.Body)

	if runtime.GOOS == "linux" {
		// This is needed to work on Linux.
		// Ref: https://github.com/golang-design/clipboard/issues/15
		select {
		case <-clipch:
		case <-time.After(1 * time.Second):
		}
	}

	return nil
}
