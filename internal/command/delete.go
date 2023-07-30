package command

import (
	"context"
	"fmt"
	"io"

	"github.com/manifoldco/promptui"
	"github.com/y-yagi/doco/internal/config"
)

type DeleteCommand struct {
	Command
	text   string
	cfg    config.Config
	stdout io.Writer
	stderr io.Writer
}

func Delete(text string, cfg config.Config, stdout, stderr io.Writer) *DeleteCommand {
	return &DeleteCommand{text: text, cfg: cfg, stdout: stdout, stderr: stderr}
}

func (c *DeleteCommand) Run() error {
	client, err := getEntClient(c.cfg.DataBase)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntriesByTitle(client, c.text)
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

	prompt := promptui.Prompt{
		Label:     "Are you really want to delete '" + selectedEntry.Title + "'",
		IsConfirm: true,
	}

	if _, err = prompt.Run(); err != nil {
		// When users choose "no", go through here.
		return nil
	}

	if err = client.Entry.DeleteOne(selectedEntry).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	fmt.Fprintf(c.stdout, "Deleted '%s'\n", selectedEntry.Title)
	return nil
}
