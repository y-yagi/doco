package config

type Config struct {
	DataBase                 string `toml:"database"`
	SelectCmd                string `toml:"selectcmd"`
	Browser                  string `toml:"browser"`
	AutomaticallyOpenBrowser bool   `toml:automatically_open_browser`
}
