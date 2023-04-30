package command

import (
	"context"
	"fmt"
	"io"
)

func Migrate(database string, stdout, stderr io.Writer) error {
	client, err := getEntClient(database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	return nil
}
