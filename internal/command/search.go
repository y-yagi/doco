package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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
		browser := c.cfg.Browser

		if len(os.Getenv("BROWSER")) > 0 {
			browser = os.Getenv("BROWSER")
		}
		browserArgs := ""

		if len(os.Getenv("BROWSER_ARGS")) > 0 {
			browserArgs = os.Getenv("BROWSER_ARGS")
		}

		cmd := exec.Command(browser, selectedEntry.Body, browserArgs)
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

	if _, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(selectedEntry.Body)); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %v", err)
	}
	fmt.Fprintf(c.stdout, "copied '%s' to clipboard\n", selectedEntry.Body)

	return nil
}
