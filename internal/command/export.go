package command

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type ExportCommand struct {
	Command
	database string
	stdout   io.Writer
	stderr   io.Writer
}

func Export(database string, stdout, stderr io.Writer) *ExportCommand {
	return &ExportCommand{database: database, stdout: stdout, stderr: stderr}
}

func (c *ExportCommand) Run() error {
	sql, err := c.generateBackupSQL()
	if err != nil {
		return err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	files := map[github.GistFilename]github.GistFile{
		"back of Doco": {Filename: github.String("export.sql"), Content: github.String(sql)},
	}
	gist, _, err := client.Gists.Create(ctx, &github.Gist{Description: github.String("backup of Doco"), Public: github.Bool(false), Files: files})
	if err != nil {
		return fmt.Errorf("creating the gist failed: %v", err)
	}

	fmt.Fprintf(c.stdout, "Data is exported to %s\n", *gist.GitPullURL)

	return nil
}

func (c *ExportCommand) generateBackupSQL() (string, error) {
	client, err := getEntClient(c.database)
	if err != nil {
		return "", err
	}

	defer client.Close()

	entries, err := getEntries(client)
	if err != nil {
		return "", fmt.Errorf("get entries failed: %v", err)
	}

	sql := ""
	for _, entry := range entries {
		sql += fmt.Sprintf("INSERT INTO entries (title, body, tag) VALUES('%s', '%s', '%s');\n", entry.Title, entry.Body, entry.Tag)
	}

	return sql, nil
}
