package command

import (
	"bytes"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	db := filepath.Join(dir, "doco.db")

	var stdout, stderr bytes.Buffer
	if err := Migrate(db, &stdout, &stderr).Run(); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	return db
}

func addEntry(t *testing.T, database, title, body, tag string) {
	t.Helper()
	client, err := getEntClient(database)
	if err != nil {
		t.Fatalf("getEntClient failed: %v", err)
	}
	defer client.Close()

	if _, err := client.Entry.Create().SetTitle(title).SetBody(body).SetTag(tag).Save(t.Context()); err != nil {
		t.Fatalf("create entry failed: %v", err)
	}
}

func allEntries(t *testing.T, database string) []*struct{ Title, Body, Tag string } {
	t.Helper()
	client, err := getEntClient(database)
	if err != nil {
		t.Fatalf("getEntClient failed: %v", err)
	}
	defer client.Close()

	entries, err := getEntries(client)
	if err != nil {
		t.Fatalf("getEntries failed: %v", err)
	}

	result := make([]*struct{ Title, Body, Tag string }, 0, len(entries))
	for _, e := range entries {
		result = append(result, &struct{ Title, Body, Tag string }{Title: e.Title, Body: e.Body, Tag: e.Tag})
	}
	return result
}

func TestExportImportRoundTrip(t *testing.T) {
	src := newTestDB(t)
	dst := newTestDB(t)

	addEntry(t, src, "it's a title", "line1\nline2", "tag1")
	addEntry(t, src, "日本語のタイトル", "本文です 🎉", "タグ")

	tmp := filepath.Join(t.TempDir(), "export.json")

	var exportOut, exportErr bytes.Buffer
	if err := Export(src, tmp, &exportOut, &exportErr).Run(); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if got := exportOut.String(); got != "Data is exported to "+tmp+"\n" {
		t.Fatalf("unexpected export output: %q", got)
	}

	var stdin bytes.Buffer
	var importOut, importErr bytes.Buffer
	if err := Import(dst, tmp, &stdin, &importOut, &importErr).Run(); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	if got := importOut.String(); got != "Data is imported (imported: 2, skipped: 0)\n" {
		t.Fatalf("unexpected import output: %q", got)
	}

	entries := allEntries(t, dst)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Title != "it's a title" || entries[0].Body != "line1\nline2" || entries[0].Tag != "tag1" {
		t.Fatalf("unexpected entry: %+v", entries[0])
	}
	if entries[1].Title != "日本語のタイトル" || entries[1].Body != "本文です 🎉" || entries[1].Tag != "タグ" {
		t.Fatalf("unexpected entry: %+v", entries[1])
	}
}

func TestImportSkipsExisting(t *testing.T) {
	src := newTestDB(t)
	dst := newTestDB(t)

	addEntry(t, src, "title1", "body1", "tag1")
	addEntry(t, src, "title2", "body2", "tag2")
	addEntry(t, dst, "title1", "existing body", "existing tag")

	tmp := filepath.Join(t.TempDir(), "export.json")
	var buf bytes.Buffer
	if err := Export(src, tmp, &buf, &buf).Run(); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	var stdin bytes.Buffer
	var out, errOut bytes.Buffer
	if err := Import(dst, tmp, &stdin, &out, &errOut).Run(); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	if got := out.String(); got != "Data is imported (imported: 1, skipped: 1)\n" {
		t.Fatalf("unexpected import output: %q", got)
	}
	if got := errOut.String(); got != "warning: skipped 'title1' (already exists)\n" {
		t.Fatalf("unexpected import warning: %q", got)
	}

	entries := allEntries(t, dst)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// re-importing again should skip everything
	var out2, errOut2 bytes.Buffer
	if err := Import(dst, tmp, &stdin, &out2, &errOut2).Run(); err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if got := out2.String(); got != "Data is imported (imported: 0, skipped: 2)\n" {
		t.Fatalf("unexpected import output: %q", got)
	}
}

func TestImportDeduplicatesWithinFile(t *testing.T) {
	dst := newTestDB(t)

	json := `{"version":1,"entries":[{"title":"dup","body":"first","tag":"a"},{"title":"dup","body":"second","tag":"b"}]}`
	stdin := bytes.NewBufferString(json)
	var out, errOut bytes.Buffer
	if err := Import(dst, "-", stdin, &out, &errOut).Run(); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	if got := out.String(); got != "Data is imported (imported: 1, skipped: 1)\n" {
		t.Fatalf("unexpected import output: %q", got)
	}
	if got := errOut.String(); got != "warning: skipped 'dup' (duplicated in the file)\n" {
		t.Fatalf("unexpected import warning: %q", got)
	}

	entries := allEntries(t, dst)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Body != "first" {
		t.Fatalf("expected first duplicate to be kept, got body %q", entries[0].Body)
	}
}

func TestImportInvalidJSON(t *testing.T) {
	dst := newTestDB(t)

	stdin := bytes.NewBufferString("not json")
	var out bytes.Buffer
	if err := Import(dst, "-", stdin, &out, &out).Run(); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}

	entries := allEntries(t, dst)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestImportUnsupportedVersion(t *testing.T) {
	dst := newTestDB(t)

	stdin := bytes.NewBufferString(`{"version":2,"entries":[{"title":"t","body":"b","tag":""}]}`)
	var out bytes.Buffer
	if err := Import(dst, "-", stdin, &out, &out).Run(); err == nil {
		t.Fatalf("expected error for unsupported version")
	}

	entries := allEntries(t, dst)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestImportEmptyTitleRollsBack(t *testing.T) {
	dst := newTestDB(t)

	json := `{"version":1,"entries":[{"title":"ok","body":"body","tag":""},{"title":"","body":"body2","tag":""}]}`
	stdin := bytes.NewBufferString(json)
	var out bytes.Buffer
	if err := Import(dst, "-", stdin, &out, &out).Run(); err == nil {
		t.Fatalf("expected error for empty title")
	}

	entries := allEntries(t, dst)
	if len(entries) != 0 {
		t.Fatalf("expected nothing inserted, got %d entries", len(entries))
	}
}

func TestExportImportViaStdio(t *testing.T) {
	src := newTestDB(t)
	dst := newTestDB(t)

	addEntry(t, src, "title", "body", "tag")

	var exported bytes.Buffer
	var exportErr bytes.Buffer
	if err := Export(src, "-", &exported, &exportErr).Run(); err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if exported.Len() == 0 {
		t.Fatalf("expected exported data on stdout")
	}
	if exportErr.Len() != 0 {
		t.Fatalf("expected nothing extra written for stdout export, got %q", exportErr.String())
	}

	var out bytes.Buffer
	if err := Import(dst, "-", &exported, &out, &out).Run(); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	entries := allEntries(t, dst)
	if len(entries) != 1 || entries[0].Title != "title" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}
