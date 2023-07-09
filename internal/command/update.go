package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/doco/internal/config"
)

func Update(text string, cfg config.Config, stdout, stderr io.Writer) error {
	client, err := getEntClient(cfg.DataBase)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntriesByTitle(client, text)
	if err != nil {
		return fmt.Errorf("search failed: %v", err)
	}

	e, err := selectEntry(stderr, stdout, cfg.SelectCmd, entries)
	if err != nil {
		return err
	}

	if e == nil {
		return nil
	}

	if err = editEntry(client, e); err != nil {
		return fmt.Errorf("updating failed: %v", err)
	}

	fmt.Fprintf(stdout, "Updted '%s'\n", e.Title)
	return nil
}
