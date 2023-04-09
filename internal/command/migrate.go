package command

import (
	"context"
	"fmt"
	"io"

	"github.com/y-yagi/doco/ent"
)

func Migrate(database string, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", database+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	return nil
}
