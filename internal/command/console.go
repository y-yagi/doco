package command

import (
	"io"
	"os"
	"os/exec"
)

func Console(database string, stdout, stderr io.Writer) error {
	cmd := exec.Command("sqlite3", database)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
