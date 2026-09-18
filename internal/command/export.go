package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type backupEntry struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Tag   string `json:"tag"`
}

// backupVersion is the version of the export format. Increase it when the format
// changes incompatibly, so that import can reject data it doesn't understand.
const backupVersion = 1

type backupData struct {
	Version int           `json:"version"`
	Entries []backupEntry `json:"entries"`
}

type ExportCommand struct {
	Command
	database string
	path     string
	stdout   io.Writer
	stderr   io.Writer
}

func Export(database, path string, stdout, stderr io.Writer) *ExportCommand {
	return &ExportCommand{database: database, path: path, stdout: stdout, stderr: stderr}
}

func (c *ExportCommand) Run() error {
	client, err := getEntClient(c.database)
	if err != nil {
		return err
	}
	defer client.Close()

	entries, err := getEntries(client)
	if err != nil {
		return fmt.Errorf("get entries failed: %v", err)
	}

	data := backupData{Version: backupVersion, Entries: make([]backupEntry, 0, len(entries))}
	for _, entry := range entries {
		data.Entries = append(data.Entries, backupEntry{Title: entry.Title, Body: entry.Body, Tag: entry.Tag})
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	if c.path == "-" {
		if _, err = c.stdout.Write(b); err != nil {
			return fmt.Errorf("write failed: %v", err)
		}
		return nil
	}

	if err = os.WriteFile(c.path, b, 0600); err != nil {
		return fmt.Errorf("write failed: %v", err)
	}

	fmt.Fprintf(c.stdout, "Data is exported to %s\n", c.path)

	return nil
}
