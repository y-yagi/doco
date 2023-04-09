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
	"golang.design/x/clipboard"
)

func Search(text, database, browser, selectCmd string, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", database+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	entries, err := client.Entry.
		Query().
		Where(func(s *sql.Selector) {
			s.Where(sql.Like(entry.FieldTitle, "%"+text+"%"))
		}).
		Select(entry.FieldID, entry.FieldTitle, entry.FieldBody, entry.FieldTag).
		All(context.Background())

	if err != nil {
		return fmt.Errorf("search failed: %v", err)
	}

	var inbuf string
	var outbuf bytes.Buffer
	dict := make(map[string]string)

	for _, entry := range entries {
		inbuf += fmt.Sprintf("%s\n", entry.Title)
		dict[entry.Title] = entry.Body
	}

	cmd := exec.Command("sh", "-c", selectCmd)
	cmd.Stderr = stderr
	cmd.Stdout = &outbuf
	cmd.Stdin = strings.NewReader(inbuf)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("select command failed: %v", err)
	}

	selected := strings.TrimRight(outbuf.String(), "\n")

	if len(selected) > 0 {
		body := dict[selected]

		if err = clipboard.Init(); err == nil {
			clipboard.Write(clipboard.FmtText, []byte(body))
		}

		if strings.HasPrefix(body, "http") {
			cmd = exec.Command(browser, body)
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("command execute failed: %v", err)
			}
		}
	}
	return nil
}
