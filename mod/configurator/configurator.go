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

	"badassops.ldap/consts"
	"badassops.ldap/vars"
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
		Debug			bool
		// from the configuration file
		DefaultValues	Defaults
		LogValues		LogConfig
		ValidEnvs		[]string
		GroupValues		Groups
		ServerValues	Server
		// passed by main
		LockPID			int
	}

	// the entries structure in the toml file
	Defaults struct {
		LockFile		string
		Shell			string
		ValidShells		[]string
		UserSearch		string
		GroupSearch		string
		GroupName		string
		GroupId			int
		ShadowMin		int
		ShadowMax		int
		ShadowAge		int
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
		SpecialGroups	[]string
		Groups			[]string
		GroupsMap		[]GroupMap
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
				utils.CreateColorMsg(consts.Yellow, "create"),
				utils.CreateColorMsg(consts.Yellow, "modify"),
				utils.CreateColorMsg(consts.Yellow, "delete"),
	)
	//searchCmd := fmt.Sprintf("\t\t     search commands:\n\t\t\t (user) %s, (group) %s\n",
	//			utils.CreateColorMsg(consts.Green, "search"),
	//			utils.CreateColorMsg(consts.Green, "group"),
	//)
	//searchAllCMD := fmt.Sprintf("\t\t     get all the users or groups records commands:\n\t\t\t (user) %s, (group) %s\n",
	//			utils.CreateColorMsg(consts.Blue, "users"),
	//			utils.CreateColorMsg(consts.Blue, "groups"),
	//)
	searchCmd := fmt.Sprintf("\t\t     search: (%s)ser, (%s)ll Users, (%s)roup and All Group(%s)",
			utils.CreateColorMsg(consts.Green, "U"),
			utils.CreateColorMsg(consts.Green, "A"),
			utils.CreateColorMsg(consts.Green, "G"),
			utils.CreateColorMsg(consts.Green, "S"))

	HelpMessage := fmt.Sprintf("%s%s", baseCmd, searchCmd)

	errored := 0
	allowedValues := []string{"create", "modify", "delete", "search"}
	parser := argparse.NewParser(vars.MyProgname, vars.MyDescription)
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

	debug := parser.Flag("d", "debug",
        &argparse.Options{
        Required:   false,
        Help:       "Enable debug",
        Default:    false,
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
		utils.PrintColor(consts.Yellow, vars.MyProgname + " version: " + vars.MyVersion + "\n")
		os.Exit(0)
	}

	if *showInfo {
		utils.ClearScreen()
		utils.PrintColor(consts.Yellow, vars.MyDescription + "\n")
		utils.PrintColor(consts.Cyan, vars.MyInfo)
		os.Exit(0)
	}

	if len(*configFile) == 0 {
		utils.PrintColor(consts.Red, "the flag -c/--config is required\n")
		errored = 1
	}

	if len(*ldapCmd) == 0 {
		utils.PrintColor(consts.Red, "the flag -m/--mode is required\n")
		errored = 1
	}

	if errored == 1 {
		utils.PrintColor(consts.Red, "Aborting..\n")
		os.Exit(1)
	}

	if ok, _ := utils.Exist(*configFile, true, false); !ok {
		utils.PrintColor(consts.Red, "Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.ConfigFile	= *configFile
	c.Env			= *ldapEnv
	c.Cmd			= *ldapCmd
	c.Debug			= *debug
}

// function to add the values to the Config object from the configuration file
func (c *Config) InitializeConfigs() {
	var configValues tomlConfig
	if _, err := toml.DecodeFile(c.ConfigFile, &configValues); err != nil {
		utils.PrintColor(consts.Red, "Error reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// from the configuration file
	c.DefaultValues		= configValues.Defaults
	c.LogValues			= configValues.LogConfig
	c.ValidEnvs			= configValues.Envs.ValidEnvs
	c.GroupValues		= configValues.Groups
	c.ServerValues		= configValues.Servers[c.Env]
}
