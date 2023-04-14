package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/ent/entry"
)

func getEntriesByTitle(client *ent.Client, text string) ([]*ent.Entry, error) {
	return client.Entry.
		Query().
		Where(func(s *sql.Selector) {
			s.Where(sql.Like(entry.FieldTitle, "%"+text+"%"))
		}).
		Order(ent.Asc(entry.FieldID)).
		Select(entry.FieldID, entry.FieldTitle, entry.FieldBody, entry.FieldTag).
		All(context.Background())
}

func getEntries(client *ent.Client) ([]*ent.Entry, error) {
	return client.Entry.
		Query().
		Order(ent.Asc(entry.FieldID)).
		Select(entry.FieldID, entry.FieldTitle, entry.FieldBody, entry.FieldTag).
		All(context.Background())
}

func selectEntry(stderr, stdout io.Writer, selectCmd string, entries []*ent.Entry) (*ent.Entry, error) {
	var inbuf string
	dict := make(map[string]*ent.Entry)
	for _, entry := range entries {
		inbuf += fmt.Sprintf("%s\n", entry.Title)
		dict[entry.Title] = entry
	}

	if len(inbuf) == 0 {
		fmt.Fprintf(stdout, "noting found\n")
		return nil, nil
	}

	var outbuf bytes.Buffer
	cmd := exec.Command("sh", "-c", selectCmd)
	cmd.Stderr = stderr
	cmd.Stdout = &outbuf
	cmd.Stdin = strings.NewReader(inbuf)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("select command failed: %v", err)
	}

	selected := strings.TrimRight(outbuf.String(), "\n")
	if len(selectCmd) == 0 {
		return nil, nil
	}

	return dict[selected], nil
}
