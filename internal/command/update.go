package command

import (
	"context"
	"fmt"
	"io"

	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/internal/config"
)

func Update(text string, cfg config.Config, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

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

	if err = inputEntryByPrompt(e); err != nil {
		return err
	}

	_, err = client.Entry.UpdateOneID(e.ID).SetTitle(e.Title).SetBody(e.Body).SetTag(e.Tag).Save(context.Background())
	if err != nil {
		return fmt.Errorf("updating failed: %v", err)
	}

	fmt.Fprintf(stdout, "Updted '%s'\n", e.Title)
	return nil
}
