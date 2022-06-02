// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package configurator

import (
	"fmt"
	"os"

	v "badassops.ldap/vars"
	"github.com/BurntSushi/toml"
	"github.com/akamensky/argparse"
	"github.com/badassops/packages-go/is"
	"github.com/badassops/packages-go/print"
)

type (
	Config struct {
		ConfigFile    string
		Server        string
		Cmd           string
		Debug         bool
		AuthValues    Auth
		DefaultValues Defaults
		SudoValues    Sudo
		LogValues     LogConfig
		EnvValues     Envs
		GroupValues   Groups
		ServerValues  Server
		RedisValues   Redis
		LockPID       int
	}

	// the entries structure in the toml file
	GroupMap struct {
		Name string
		Gid  int
	}

	Auth struct {
		AllowUsers []string
		AllowMods  []string
	}

	Defaults struct {
		LockFile      string
		Shell         string
		ValidShells   []string
		UserSearch    string
		GroupSearch   string
		GroupName     string
		GroupId       int
		ShadowMin     int
		ShadowMax     int
		ShadowAge     int
		ShadowWarning int
		Wait          int
		PassLenght    int
		PassComplex   bool
		UidStart      int
		GidStart      int
	}

	Sudo struct {
		ExcludeSudo []string
		SudoersBase string
	}

	LogConfig struct {
		LogsDir       string
		LogFile       string
		LogMaxSize    int
		LogMaxBackups int
		LogMaxAge     int
	}

	Envs struct {
		ValidEnvs []string
	}

	Groups struct {
		SpecialGroups []string
		Groups        []string
		GroupsMap     []GroupMap
	}

	Server struct {
		Server      string
		BaseDN      string
		Admin       string
		AdminPass   string
		UserDN      string
		GroupDN     string
		EmailDomain string
		TLS         bool
		Enabled     bool
		ReadOnly    bool
	}

	Redis struct {
		Server  string
		Port    int
		Enabled bool
		TmpFile string
	}

	tomlConfig struct {
		Auth      Auth              `toml:"auth"`
		Defaults  Defaults          `toml:"defaults"`
		Sudo      Sudo              `toml:"sudo"`
		LogConfig LogConfig         `toml:"logconfig"`
		Envs      Envs              `toml:"envs"`
		Groups    Groups            `toml:"groups"`
		Servers   map[string]Server `toml:"servers"`
		Redis     Redis             `toml:"redis"`
	}
)

var (
	Is    = is.New()
	Print = print.New()
)

// function to initialize the configuration
func Configurator() *Config {
	// the rest of the values will be filled from the given configuration file
	return &Config{
		ConfigFile: "",
		Server:     "",
		Cmd:        "",
	}
}

func (c *Config) InitializeArgs() {
	HelpMessage := fmt.Sprintf("commands: %s, %s, %s, %s\n",
		Print.MessageYellow("search"),
		Print.MessageYellow("create"),
		Print.MessageYellow("modify"),
		Print.MessageYellow("delete"),
	)

	errored := 0
	allowedValues := []string{"create", "modify", "delete", "search"}
	parser := argparse.NewParser(v.MyProgname, v.MyDescription)
	configFile := parser.String("c", "configFile",
		&argparse.Options{
			Required: false,
			Help:     "Path to the configuration file to be use",
			Default:  "/usr/local/etc/ldap-tool/ldap-tool.ini",
		})

	server := parser.String("s", "server",
		&argparse.Options{
			Required: false,
			Help:     "Server profile name",
		})

	cmd := parser.Selector("C", "command", allowedValues,
		&argparse.Options{
			Required: false,
			Help:     HelpMessage,
			Default:  "search",
		})

	debug := parser.Flag("d", "debug",
		&argparse.Options{
			Required: false,
			Help:     "Enable debug",
			Default:  false,
		})

	showInfo := parser.Flag("i", "info",
		&argparse.Options{
			Required: false,
			Help:     "Show information",
		})

	showVersion := parser.Flag("v", "version",
		&argparse.Options{
			Required: false,
			Help:     "Show version",
		})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *showVersion {
		Print.ClearScreen()
		Print.PrintYellow(v.MyProgname + " version: " + v.MyVersion + "\n")
		os.Exit(0)
	}

	if *showInfo {
		Print.ClearScreen()
		Print.PrintYellow(v.MyDescription + "\n")
		Print.PrintCyan(v.MyInfo)
		os.Exit(0)
	}

	if len(*configFile) == 0 {
		Print.PrintRed("the flag -c/--config is required\n")
		errored = 1
	}

	if len(*server) == 0 {
		Print.PrintRed("the flag -s/--server is required\n")
		errored = 1
	}

	if len(*cmd) == 0 {
		Print.PrintRed("the flag -C/--command is required\n")
		errored = 1
	}

	if errored == 1 {
		Print.PrintRed("Aborting..\n")
		os.Exit(1)
	}

	if _, ok, _ := Is.IsExist(*configFile, "file"); !ok {
		Print.PrintRed("Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.ConfigFile = *configFile
	c.Server = *server
	c.Cmd = *cmd
	c.Debug = *debug
}

// function to add the values to the Config object from the configuration file
func (c *Config) InitializeConfigs() {
	var configValues tomlConfig
	if _, err := toml.DecodeFile(c.ConfigFile, &configValues); err != nil {
		Print.PrintRed("Error reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// from the configuration file
	c.AuthValues = configValues.Auth
	c.DefaultValues = configValues.Defaults
	c.SudoValues = configValues.Sudo
	c.LogValues = configValues.LogConfig
	c.EnvValues = configValues.Envs
	c.GroupValues = configValues.Groups
	c.ServerValues = configValues.Servers[c.Server]
	c.RedisValues = configValues.Redis
}
