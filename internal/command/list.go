package command

import (
	"fmt"
	"io"

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

	for _, entry := range entries {
		fmt.Fprintf(stdout, "%s: %s\n", entry.Title, entry.Body)
	}

	return nil
}
