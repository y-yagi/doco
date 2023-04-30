package command

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
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

	table := tablewriter.NewWriter(stdout)
	table.SetHeader([]string{"Title", "Tag", "Body"})

	for _, entry := range entries {
		table.Append([]string{entry.Title, entry.Tag, entry.Body})
	}

	table.Render()

	return nil
}
