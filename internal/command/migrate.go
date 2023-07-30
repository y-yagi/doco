package command

import (
	"context"
	"fmt"
	"io"
)

type MigrateCommand struct {
	Command
	database string
	stdout   io.Writer
	stderr   io.Writer
}

func Migrate(database string, stdout, stderr io.Writer) *MigrateCommand {
	return &MigrateCommand{database: database, stdout: stdout, stderr: stderr}
}

func (c *MigrateCommand) Run() error {
	client, err := getEntClient(c.database)
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
