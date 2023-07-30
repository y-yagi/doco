package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/expandedwriter"
)

type ListCommand struct {
	Command
	database string
	stdout   io.Writer
	stderr   io.Writer
}

func List(database string, stdout, stderr io.Writer) *ListCommand {
	return &ListCommand{database: database, stdout: stdout, stderr: stderr}
}

func (c *ListCommand) Run() error {
	client, err := getEntClient(c.database)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntries(client)
	if err != nil {
		return fmt.Errorf("get entries failed: %v", err)
	}

	w := expandedwriter.NewWriter(c.stdout)
	w.SetFields([]string{"Title", "Tag", "Body"})

	for _, entry := range entries {
		w.Append([]string{entry.Title, entry.Tag, entry.Body})
	}

	return w.Render()
}
