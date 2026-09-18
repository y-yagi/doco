package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/ent/entry"
)

type ImportCommand struct {
	Command
	database string
	path     string
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
}

func Import(database, path string, stdin io.Reader, stdout, stderr io.Writer) *ImportCommand {
	return &ImportCommand{database: database, path: path, stdin: stdin, stdout: stdout, stderr: stderr}
}

func (c *ImportCommand) Run() error {
	var r io.Reader
	if c.path == "-" {
		r = c.stdin
	} else {
		f, err := os.Open(c.path)
		if err != nil {
			return fmt.Errorf("open failed: %v", err)
		}
		defer f.Close()
		r = f
	}

	var data backupData
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Errorf("decode failed: %v", err)
	}

	if data.Version != backupVersion {
		return fmt.Errorf("unsupported version: %d (supported: %d)", data.Version, backupVersion)
	}

	client, err := getEntClient(c.database)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction failed: %v", err)
	}

	imported, skipped, err := importEntries(ctx, tx, data.Entries)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("import failed: %v, rolling back failed: %v", err, rerr)
		}
		return fmt.Errorf("import failed: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}

	for _, s := range skipped {
		fmt.Fprintf(c.stderr, "warning: skipped '%s' (%s)\n", s.title, s.reason)
	}

	fmt.Fprintf(c.stdout, "Data is imported (imported: %d, skipped: %d)\n", imported, len(skipped))

	return nil
}

type skippedEntry struct {
	title  string
	reason string
}

func importEntries(ctx context.Context, tx *ent.Tx, entries []backupEntry) (int, []skippedEntry, error) {
	titles := make([]string, 0, len(entries))
	for _, e := range entries {
		titles = append(titles, e.Title)
	}

	existing, err := tx.Entry.Query().Where(entry.TitleIn(titles...)).Select(entry.FieldTitle).All(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("checking existing entries failed: %v", err)
	}

	exists := make(map[string]bool, len(existing))
	for _, e := range existing {
		exists[e.Title] = true
	}

	var skipped []skippedEntry
	seen := make(map[string]bool, len(entries))
	builders := make([]*ent.EntryCreate, 0, len(entries))
	for _, e := range entries {
		if exists[e.Title] {
			skipped = append(skipped, skippedEntry{title: e.Title, reason: "already exists"})
			continue
		}
		if seen[e.Title] {
			skipped = append(skipped, skippedEntry{title: e.Title, reason: "duplicated in the file"})
			continue
		}
		seen[e.Title] = true

		builders = append(builders, tx.Entry.Create().SetTitle(e.Title).SetBody(e.Body).SetTag(e.Tag))
	}

	if len(builders) > 0 {
		if _, err := tx.Entry.CreateBulk(builders...).Save(ctx); err != nil {
			return 0, nil, fmt.Errorf("creating entries failed: %v", err)
		}
	}

	return len(builders), skipped, nil
}
