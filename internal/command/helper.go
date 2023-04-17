package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/manifoldco/promptui"
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

func inputEntryByPrompt(entry *ent.Entry) error {
	var err error
	emptyVaidate := func(input string) error {
		if len(input) < 1 {
			return errors.New("value must have more than 1 character")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Title",
		Validate: emptyVaidate,
		Default:  entry.Title,
	}

	if entry.Title, err = prompt.Run(); err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	prompt = promptui.Prompt{
		Label:    "Body",
		Validate: emptyVaidate,
		Default:  entry.Body,
	}

	if entry.Body, err = prompt.Run(); err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	prompt = promptui.Prompt{
		Label:   "Tag",
		Default: entry.Tag,
	}

	if entry.Tag, err = prompt.Run(); err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	return nil
}
