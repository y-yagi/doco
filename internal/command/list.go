package command

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/y-yagi/doco/ent"
)

func List(database string, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", database+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

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
