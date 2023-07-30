package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/doco/internal/config"
)

type UpdateCommand struct {
	Command
	text   string
	cfg    config.Config
	stdout io.Writer
	stderr io.Writer
}

func Update(text string, cfg config.Config, stdout, stderr io.Writer) *UpdateCommand {
	return &UpdateCommand{text: text, cfg: cfg, stdout: stdout, stderr: stderr}
}

func (c *UpdateCommand) Run() error {
	client, err := getEntClient(c.cfg.DataBase)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntriesByTitle(client, c.text)
	if err != nil {
		return fmt.Errorf("search failed: %v", err)
	}

	e, err := selectEntry(c.stderr, c.stdout, c.cfg.SelectCmd, entries)
	if err != nil {
		return err
	}

	if e == nil {
		return nil
	}

	if err = editEntry(client, e); err != nil {
		return fmt.Errorf("updating failed: %v", err)
	}

	fmt.Fprintf(c.stdout, "Updted '%s'\n", e.Title)
	return nil
}
