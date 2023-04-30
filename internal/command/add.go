package command

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/y-yagi/doco/ent"
)

func Add(database string, stdout, stderr io.Writer) error {
	client, err := getEntClient(database)
	if err != nil {
		return err
	}
	defer client.Close()

	entry := ent.Entry{}
	editor := os.Getenv("DOCO_EDITOR")

	if len(editor) != 0 {
		if err = inputEntryByEditor(&entry, editor); err != nil {
			return err
		}
	} else {
		if err = inputEntryByPrompt(&entry); err != nil {
			return err
		}
	}

	_, err = client.Entry.Create().SetTitle(entry.Title).SetBody(entry.Body).SetTag(entry.Tag).Save(context.Background())
	if err != nil {
		return fmt.Errorf("adding failed: %v", err)
	}

	fmt.Fprintf(stdout, "Added '%s'\n", entry.Title)
	return nil
}
