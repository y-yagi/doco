package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/y-yagi/configure"
	"github.com/y-yagi/doco/ent/entry"
	"github.com/y-yagi/doco/internal/command"
	"github.com/y-yagi/doco/internal/config"
)

const cmd = "doco"

var (
	cfg           config.Config
	flags         *flag.FlagSet
	showVersion   bool
	migrateFlag   bool
	consoleFlag   bool
	addFlag       bool
	updateFlag    string
	deleteFlag    string
	listFlag      bool
	configureFlag bool
	tagFlag       string
	exportFlag    bool
	importFlag    string

	version = "devel"
)

func init() {
	err := configure.Load(cmd, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	needToSave := false

	if len(cfg.DataBase) == 0 {
		cfg.DataBase = filepath.Join(configure.ConfigDir(cmd), cmd+".db")
		needToSave = true
	}

	if len(cfg.SelectCmd) == 0 {
		cfg.SelectCmd = "peco"
		needToSave = true
	}

	if len(cfg.Browser) == 0 {
		cfg.Browser = "google-chrome"
		needToSave = true
	}

	if needToSave {
		if err = configure.Save(cmd, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	setFlags()
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func setFlags() {
	flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.BoolVar(&showVersion, "v", false, "print version number")
	flags.BoolVar(&migrateFlag, "migrate", false, "run DB migration")
	flags.BoolVar(&addFlag, "add", false, "add new entry")
	flags.BoolVar(&consoleFlag, "console", false, "run DB console")
	flags.StringVar(&deleteFlag, "delete", "", "delete entry")
	flags.BoolVar(&listFlag, "list", false, "show all entries")
	flags.StringVar(&updateFlag, "update", "", "update entry")
	flags.BoolVar(&configureFlag, "configure", false, "edit config")
	flags.StringVar(&tagFlag, "tag", "", "search entry by tag")
	flags.BoolVar(&exportFlag, "export", false, "export data to the Gist")
	flags.StringVar(&importFlag, "import", "", "import data from Gist")
	flags.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: %s [OPTIONS] TEXT\n\n", cmd)
	fmt.Fprintln(os.Stdout, "OPTIONS:")
	flags.PrintDefaults()
}

func msg(err error, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
		return 1
	}
	return 0
}

func run(args []string, stdout, stderr io.Writer) int {
	if err := flags.Parse(args[1:]); err != nil {
		return msg(err, stderr)
	}

	if showVersion {
		fmt.Fprintf(stdout, "%s %s\n", cmd, version)
		return 0
	}

	if migrateFlag {
		return msg(command.Migrate(cfg.DataBase, stdout, stderr), stderr)
	}

	if _, err := os.Stat(cfg.DataBase); err != nil {
		if err = command.Migrate(cfg.DataBase, stdout, stderr); err != nil {
			return msg(err, stderr)
		}
	}

	if addFlag {
		return msg(command.Add(cfg.DataBase, stdout, stderr), stderr)
	}

	if consoleFlag {
		return msg(command.Console(cfg.DataBase, stdout, stderr), stderr)
	}

	if len(deleteFlag) != 0 {
		return msg(command.Delete(deleteFlag, cfg, stdout, stderr), stderr)
	}

	if len(updateFlag) != 0 {
		return msg(command.Update(updateFlag, cfg, stdout, stderr), stderr)
	}

	if listFlag {
		return msg(command.List(cfg.DataBase, stdout, stderr), stderr)
	}

	if len(tagFlag) != 0 {
		return msg(command.Search(entry.FieldTag, tagFlag, cfg, stdout, stderr), stderr)
	}

	if exportFlag {
		return msg(command.Export(cfg.DataBase, stdout, stderr), stderr)
	}

	if len(importFlag) != 0 {
		return msg(command.Import(cfg.DataBase, importFlag, stdout, stderr), stderr)
	}

	if configureFlag {
		return msg(configure.Edit(cmd, os.Getenv("DOCO_EDITOR")), stderr)
	}

	if len(flags.Args()) != 1 {
		flags.Usage()
		return 1
	}

	return msg(command.Search(entry.FieldTitle, flags.Arg(0), cfg, stdout, stderr), stderr)
}
