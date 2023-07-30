package command

import (
	"io"
	"os"
	"os/exec"
)

type ConsoleCommand struct {
	Command
	database string
	stdout   io.Writer
	stderr   io.Writer
}

func Console(database string, stdout, stderr io.Writer) *ConsoleCommand {
	return &ConsoleCommand{database: database, stdout: stdout, stderr: stderr}
}

func (c *ConsoleCommand) Run() error {
	cmd := exec.Command("sqlite3", c.database)
	cmd.Stdin = os.Stdin
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
