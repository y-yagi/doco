package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/doco/ent"
)

func Add(database string, stdout, stderr io.Writer) error {
	client, err := getEntClient(database)
	if err != nil {
		return err
	}
	defer client.Close()

	entry := ent.Entry{}
	if err = editEntry(client, &entry); err != nil {
		return fmt.Errorf("adding failed: %v", err)
	}

	fmt.Fprintf(stdout, "Added '%s'\n", entry.Title)
	return nil
}
