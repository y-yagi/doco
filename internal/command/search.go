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
	"github.com/y-yagi/doco/internal/config"
	"golang.design/x/clipboard"
)

func Search(text string, cfg config.Config, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	entries, err := getEntries(client, text)
	if err != nil {
		return fmt.Errorf("search failed: %v", err)
	}

	selectedBody, err := selectEntry(stderr, stdout, cfg.SelectCmd, entries)
	if err != nil {
		return err
	}

	if len(selectedBody) == 0 {
		return nil
	}

	if err = clipboard.Init(); err == nil {
		clipboard.Write(clipboard.FmtText, []byte(selectedBody))
	}

	if cfg.AutomaticallyOpenBrowser && strings.HasPrefix(selectedBody, "http") {
		cmd := exec.Command(cfg.Browser, selectedBody)
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("command execute failed: %v", err)
		}
	}

	return nil
}

func getEntries(client *ent.Client, text string) ([]*ent.Entry, error) {
	return client.Entry.
		Query().
		Where(func(s *sql.Selector) {
			s.Where(sql.Like(entry.FieldTitle, "%"+text+"%"))
		}).
		Select(entry.FieldID, entry.FieldTitle, entry.FieldBody, entry.FieldTag).
		All(context.Background())
}

func selectEntry(stderr, stdout io.Writer, selectCmd string, entries []*ent.Entry) (string, error) {
	var inbuf string
	dict := make(map[string]string)
	for _, entry := range entries {
		inbuf += fmt.Sprintf("%s\n", entry.Title)
		dict[entry.Title] = entry.Body
	}

	if len(inbuf) == 0 {
		fmt.Fprintf(stdout, "noting found\n")
		return "", nil
	}

	var outbuf bytes.Buffer
	cmd := exec.Command("sh", "-c", selectCmd)
	cmd.Stderr = stderr
	cmd.Stdout = &outbuf
	cmd.Stdin = strings.NewReader(inbuf)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("select command failed: %v", err)
	}

	selected := strings.TrimRight(outbuf.String(), "\n")
	if len(selectCmd) == 0 {
		return "", nil
	}

	return dict[selected], nil
}
