package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type ImportCommand struct {
	Command
	database string
	gistURL  string
	stdout   io.Writer
	stderr   io.Writer
}

func Import(database, gistURL string, stdout, stderr io.Writer) *ImportCommand {
	return &ImportCommand{database: database, gistURL: gistURL, stdout: stdout, stderr: stderr}
}

func (c *ImportCommand) Run() error {
	gistID, found := strings.CutPrefix(c.gistURL, "https://gist.github.com/")
	if !found {
		return errors.New("the URL need to start with 'https://gist.github.com'")
	}
	gistID, _ = strings.CutSuffix(gistID, ".git")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	gist, _, err := client.Gists.Get(ctx, gistID)
	if err != nil {
		return fmt.Errorf("getting the gist failed: %v", err)
	}

	sql := gist.Files[github.GistFilename("export.sql")].Content
	if err = c.insertBackupSQL(*sql); err != nil {
		return fmt.Errorf("inserting backup is failed: %v", err)
	}

	fmt.Fprintln(c.stdout, "Data is imported")

	return nil
}

func (c *ImportCommand) insertBackupSQL(sql string) error {
	client, err := getEntClient(c.database)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.ExecContext(context.Background(), sql)
	return err
}
