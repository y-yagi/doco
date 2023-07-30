package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/doco/ent"
)

type AddCommand struct {
	Command
	database string
	stdout   io.Writer
	stderr   io.Writer
}

func Add(database string, stdout, stderr io.Writer) *AddCommand {
	return &AddCommand{database: database, stdout: stdout, stderr: stderr}
}

func (c *AddCommand) Run() error {
	client, err := getEntClient(c.database)
	if err != nil {
		return err
	}
	defer client.Close()

	entry := ent.Entry{}
	if err = editEntry(client, &entry); err != nil {
		return fmt.Errorf("adding failed: %v", err)
	}

	fmt.Fprintf(c.stdout, "Added '%s'\n", entry.Title)
	return nil
}
