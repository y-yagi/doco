package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/y-yagi/configure"
)

const cmd = "doco"

type config struct {
	DataBase  string `toml:"database"`
	SelectCmd string `toml:"selectcmd"`
	Browser   string `toml:"browser"`
}

var (
	cfg         config
	flags       *flag.FlagSet
	showVersion bool
	migrateFlag bool
	consoleFlag bool
	addFlag     bool

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
		configure.Save(cmd, cfg)
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
	flags.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] TEXT\n\n", cmd)
	fmt.Fprintln(os.Stderr, "OPTIONS:")
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
		return msg(migrate(stdout, stderr), stderr)
	}

	if addFlag {
		return msg(add(stdout, stderr), stderr)
	}

	if consoleFlag {
		return msg(console(stdout, stderr), stderr)
	}

	if len(flags.Args()) != 1 {
		flags.Usage()
		return 1
	}

	return msg(search(flags.Arg(0), stdout, stderr), stderr)
}
