package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/manifoldco/promptui"
	"github.com/pelletier/go-toml/v2"
	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/ent/entry"
)

type TmpEntry struct {
	Title string `toml:"title"`
	Body  string `toml:"body"`
	Tag   string `toml:"tag"`
}

func getEntriesByTitle(client *ent.Client, text string) ([]*ent.Entry, error) {
	return getEntriesBy(client, entry.FieldTitle, text)
}

func getEntriesBy(client *ent.Client, field, text string) ([]*ent.Entry, error) {
	return client.Entry.
		Query().
		Where(func(s *sql.Selector) {
			s.Where(sql.Like(field, "%"+text+"%")).Where(sql.NEQ(field, ""))
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
	lengthVaidate := func(input string) error {
		if len(input) < 1 {
			return errors.New("value must have more than 1 character")
		}
		return nil
	}

	emptyVaidate := func(_input string) error {
		return nil
	}

	getInputFromPrompt := func(label string, validate promptui.ValidateFunc, defau string) (string, error) {
		prompt := promptui.Prompt{
			Label:    label,
			Validate: validate,
			Default:  defau,
		}

		return prompt.Run()
	}

	if entry.Title, err = getInputFromPrompt("Title", lengthVaidate, entry.Title); err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	if entry.Body, err = getInputFromPrompt("Body", lengthVaidate, entry.Body); err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	if entry.Tag, err = getInputFromPrompt("Tag", emptyVaidate, entry.Tag); err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	return nil
}

func inputEntryByEditor(entry *ent.Entry, editor string) error {
	file, err := os.CreateTemp("", "doco")
	if err != nil {
		return fmt.Errorf("create a tmp file failed: %v", err)
	}
	defer os.Remove(file.Name())

	tmpEntry := TmpEntry{Title: entry.Title, Body: entry.Body, Tag: entry.Tag}
	if err = toml.NewEncoder(file).Encode(tmpEntry); err != nil {
		return fmt.Errorf("encode failed: %v", err)
	}

	if err = runEditor(editor, file.Name()); err != nil {
		return err
	}

	if err = readEntryFromFile(file, &tmpEntry); err != nil {
		return err
	}

	entry.Title = tmpEntry.Title
	entry.Body = tmpEntry.Body
	entry.Tag = tmpEntry.Tag

	return nil
}

func runEditor(editor, filename string) error {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor starting failed: %v", err)
	}

	return nil
}

func readEntryFromFile(file *os.File, e *TmpEntry) error {
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("file seek failed: %v", err)
	}

	if err := toml.NewDecoder(file).Decode(e); err != nil {
		return fmt.Errorf("decode failed: %v", err)
	}

	if len(e.Title) == 0 || len(e.Body) == 0 {
		return fmt.Errorf("title and body can't be empty")
	}

	return nil
}
