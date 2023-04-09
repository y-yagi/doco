package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/manifoldco/promptui"
	"github.com/y-yagi/doco/ent"
)

func Add(database string, stdout, stderr io.Writer) error {
	client, err := ent.Open("sqlite3", database+"?_fk=1")
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
