package command

import (
	"context"
	"fmt"
	"io"

	"github.com/y-yagi/doco/ent"
)

func Add(database string, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", database+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	entry := ent.Entry{}
	if err = inputEntryByPrompt(&entry); err != nil {
		return err
	}

	_, err = client.Entry.Create().SetTitle(entry.Title).SetBody(entry.Body).SetTag(entry.Tag).Save(context.Background())
	if err != nil {
		return fmt.Errorf("adding failed: %v", err)
	}

	fmt.Fprintf(stdout, "Added '%s'\n", entry.Title)
	return nil
}
