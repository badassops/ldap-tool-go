// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package configurator

import (
	"fmt"
	"os"

	"badassops.ldap/constants"
	"badassops.ldap/utils"

	"github.com/akamensky/argparse"
	"github.com/BurntSushi/toml"
)

type (
	GroupMap struct {
			Name	string
			Gid		int
	}

	Config struct {
		ConfigFile		string
		Env				string
		Cmd				string
		// from the configuration file
		Shell			string
		ValidShells		[]string
		UserSearch		string
		GroupSearch		string
		GroupName		string
		GroupId			int
		ShadowMax		int
		ShadowWarning	int
		Wait			int
		PassLenght		int
		PassComplex		bool
		// from the configuration file
		LogsDir         string
		LogFile         string
		LogMaxSize      int
		LogMaxBackups   int
		LogMaxAge		int

		ValidEnvs		[]string

		Admins			[]string
		VPNs			[]string
		Groups			[]string
		GroupsMap		[]GroupMap

		Server			string
		BaseDN			string
		Admin			string
		AdminPass		string
		UserDN			string
		GroupDN			string
		EmailDomain		string
		TLS				bool
		Enabled			bool
		LockFile		string

		// passed by main
		LockPID			int
	}

	// the entries structure in the toml file
	Defaults struct {
		LockFile		string
		Shell			string
		UserSearch		string
		GroupSearch		string
		ValidShells		[]string
		GroupName		string
		GroupId			int
		ShadowMax		int
		ShadowWarning	int
		Wait			int
		PassLenght		int
		PassComplex		bool
	}

	LogConfig struct {
		LogsDir         string
		LogFile         string
		LogMaxSize      int
		LogMaxBackups   int
		LogMaxAge       int
	}

	Envs struct {
		ValidEnvs		[]string
	}

	Groups struct {
		Admins		[]string
		VPNs		[]string
		Groups		[]string
		GroupsMap	[]GroupMap
	}

	Server struct {
		Server			string
		BaseDN			string
		Admin			string
		AdminPass		string
		UserDN			string
		GroupDN			string
		EmailDomain		string
		TLS				bool
		Enabled			bool
	}

	tomlConfig struct {
		Defaults	Defaults			`toml:"defaults"`
		LogConfig	LogConfig			`toml:"logconfig"`
		Envs		Envs				`toml:"envs"`
		Groups		Groups				`toml:"groups"`
		Servers		map[string]Server	`toml:"servers"`
	}
)

// function to initialize the configuration
func Configurator() *Config {
	// the rest of the values will be filled from the given configuration file
	return &Config {
		ConfigFile:			"",
		Env:				"",
		Cmd:				"",
	}
}

func (c *Config) InitializeArgs() {
	baseCmd := fmt.Sprintf("base commands:\n\t\t\t %s, %s, %s\n",
				utils.CreateColorMsg(constants.Yellow, "create"),
				utils.CreateColorMsg(constants.Yellow, "modify"),
				utils.CreateColorMsg(constants.Yellow, "delete"),
	)
	searchCmd := fmt.Sprintf("\t\t     search commands:\n\t\t\t (user) %s, (group) %s, %s\n",
				utils.CreateColorMsg(constants.Green, "search"),
				utils.CreateColorMsg(constants.Green, "group"),
				utils.CreateColorMsg(constants.Green, "admin"),
	)
	searchAllCMD := fmt.Sprintf("\t\t     get all records users and groups commands:\n\t\t\t (user) %s, (group) %s, %s\n",
				utils.CreateColorMsg(constants.Blue, "users"),
				utils.CreateColorMsg(constants.Blue, "groups"),
				utils.CreateColorMsg(constants.Blue, "admins"),
	)

	HelpMessage := fmt.Sprintf("%s%s%s", baseCmd, searchCmd, searchAllCMD)

	errored := 0
	allowedValues := []string{"create", "modify", "delete", "search", "group", "admin", "users", "groups", "admins"}
	parser := argparse.NewParser(constants.MyProgname, constants.MyDescription)
	configFile := parser.String("c", "configFile",
		&argparse.Options{
		Required:	false,
		Help:		"Path to the configuration file to be use",
	})

	ldapEnv := parser.String("e", "environment",
		&argparse.Options{
		Required:	false,
		Help:		"Server environment",
		Default:	"dev",
	})

	ldapCmd := parser.Selector("m", "mode", allowedValues,
		&argparse.Options{
		Required:	false,
		Help:		HelpMessage,
	})

	showInfo := parser.Flag("i", "info",
		&argparse.Options{
		Required:	false,
		Help:		"Show information",
	})

	showVersion := parser.Flag("v", "version",
		&argparse.Options{
		Required:	false,
		Help:		"Show version",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *showVersion {
		utils.ClearScreen()
		utils.PrintColor(constants.Yellow, constants.MyProgname + " version: " + constants.MyVersion + "\n")
		os.Exit(0)
	}

	if *showInfo {
		utils.ClearScreen()
		utils.PrintColor(constants.Yellow, constants.MyDescription + "\n")
		utils.PrintColor(constants.Cyan, constants.MyInfo)
		os.Exit(0)
	}

	if len(*configFile) == 0 {
		utils.PrintColor(constants.Red, "the flag -c/--config is required\n")
		errored = 1
	}

	if len(*ldapCmd) == 0 {
		utils.PrintColor(constants.Red, "the flag -m/--mode is required\n")
		errored = 1
	}

	if errored == 1 {
		utils.PrintColor(constants.Red, "Aborting..\n")
		os.Exit(1)
	}

	if ok, _ := utils.Exist(*configFile, true, false); !ok {
		utils.PrintColor(constants.Red, "Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.ConfigFile	= *configFile
	c.Env			= *ldapEnv
	c.Cmd			= *ldapCmd
}

// function to add the values to the Config object from the configuration file
func (c *Config) InitializeConfigs() {
	var configValues tomlConfig
	if _, err := toml.DecodeFile(c.ConfigFile, &configValues); err != nil {
		utils.PrintColor(constants.Red, "Error reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// from the configuration file
	c.LockFile			= configValues.Defaults.LockFile
	c.Shell				= configValues.Defaults.Shell
	c.ValidShells		= configValues.Defaults.ValidShells
	c.UserSearch		= configValues.Defaults.UserSearch
	c.GroupSearch		= configValues.Defaults.GroupSearch
	c.GroupName			= configValues.Defaults.GroupName
	c.GroupId			= configValues.Defaults.GroupId
	c.ShadowMax			= configValues.Defaults.ShadowMax
	c.ShadowWarning		= configValues.Defaults.ShadowWarning
	c.Wait				= configValues.Defaults.Wait
	c.PassLenght		= configValues.Defaults.PassLenght
	c.PassComplex		= configValues.Defaults.PassComplex
	c.LogsDir			= configValues.LogConfig.LogsDir
	c.LogFile			= configValues.LogConfig.LogFile
	c.LogMaxSize		= configValues.LogConfig.LogMaxSize
	c.LogMaxBackups		= configValues.LogConfig.LogMaxBackups
	c.LogMaxAge			= configValues.LogConfig.LogMaxAge
	c.ValidEnvs			= configValues.Envs.ValidEnvs
	c.Admins			= configValues.Groups.Admins
	c.VPNs				= configValues.Groups.VPNs
	c.Groups			= configValues.Groups.Groups
	c.GroupsMap			= configValues.Groups.GroupsMap
	c.Server			= configValues.Servers[c.Env].Server
	c.BaseDN			= configValues.Servers[c.Env].BaseDN
	c.Admin				= configValues.Servers[c.Env].Admin
	c.AdminPass			= configValues.Servers[c.Env].AdminPass
	c.UserDN			= configValues.Servers[c.Env].UserDN
	c.GroupDN			= configValues.Servers[c.Env].GroupDN
	c.EmailDomain		= configValues.Servers[c.Env].EmailDomain
	c.TLS				= configValues.Servers[c.Env].TLS
	c.Enabled			= configValues.Servers[c.Env].Enabled
}
