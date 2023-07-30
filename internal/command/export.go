package command

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func Export(database string, stdout, stderr io.Writer) error {
	sql, err := generateBackupSQL(database)
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

	fmt.Fprintf(stdout, "Data is exported to %s\n", *gist.GitPullURL)

	return nil
}

func generateBackupSQL(database string) (string, error) {
	client, err := getEntClient(database)
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
