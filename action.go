package main

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
	"github.com/y-yagi/doco/ent"
	"github.com/y-yagi/doco/ent/entry"
	"golang.design/x/clipboard"
)

func migrate(stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	return nil
}

func add(stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	emptyVaidate := func(input string) error {
		if len(input) < 1 {
			return errors.New("value must have more than 1 character")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Title",
		Validate: emptyVaidate,
	}

	title, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompty failed: %v", err)
	}

	prompt = promptui.Prompt{
		Label:    "Body",
		Validate: emptyVaidate,
	}

	body, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompty failed: %v", err)
	}

	prompt = promptui.Prompt{
		Label: "Tag",
	}

	tag, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompty failed: %v", err)
	}

	entry, err := client.Entry.Create().SetTitle(title).SetBody(body).SetTag(tag).Save(context.Background())
	if err != nil {
		return fmt.Errorf("adding failed: %v", err)
	}

	fmt.Fprintf(stdout, "Added '%s'\n", entry.Title)
	return nil
}

func console(stdout, stderr io.Writer) error {
	cmd := exec.Command("sqlite3", cfg.DataBase)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func search(text string, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", cfg.DataBase+"?_fk=1")
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

	cmd := exec.Command("sh", "-c", cfg.SelectCmd)
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
			cmd = exec.Command(cfg.Browser, body)
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("command execute failed: %v", err)
			}
		}
	}
	return nil
}
