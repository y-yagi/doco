package command

import (
	"fmt"
	"io"

	"github.com/y-yagi/expandedwriter"
)

func List(database string, stdout, stderr io.Writer) error {
	client, err := getEntClient(database)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntries(client)
	if err != nil {
		return fmt.Errorf("get entries failed: %v", err)
	}

	w := expandedwriter.NewWriter(stdout)
	w.SetFields([]string{"Title", "Tag", "Body"})

	for _, entry := range entries {
		w.Append([]string{entry.Title, entry.Tag, entry.Body})
	}

	w.Render()

	return nil
}
