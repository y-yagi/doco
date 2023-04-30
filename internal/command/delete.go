package command

import (
	"context"
	"fmt"
	"io"

	"github.com/manifoldco/promptui"
	"github.com/y-yagi/doco/internal/config"
)

func Delete(text string, cfg config.Config, stdout, stderr io.Writer) error {
	client, err := getEntClient(cfg.DataBase)
	if err != nil {
		return err
	}
	defer client.Close()

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

	fmt.Fprintf(stdout, "Deleted '%s'\n", selectedEntry.Title)
	return nil
}
