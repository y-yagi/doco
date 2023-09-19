package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/expandedwriter"
)

type ListCommand struct {
	Command
	database string
	text     string
	field    string
	stdout   io.Writer
	stderr   io.Writer
}

func List(database, field, text string, stdout, stderr io.Writer) *ListCommand {
	return &ListCommand{database: database, text: text, field: field, stdout: stdout, stderr: stderr}
}

func (c *ListCommand) Run() error {
	client, err := getEntClient(c.database)
	if err != nil {
		return err
	}
	defer client.Close()

	var entries []*ent.Entry

	if len(c.text) == 0 {
		entries, err = getEntries(client)
	} else {
		entries, err = getEntriesBy(client, c.field, c.text)
	}

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
