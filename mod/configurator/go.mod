module configurator

go 1.18

require (
	badassops.ldap/utils v0.0.0-00010101000000-000000000000

	github.com/akamensky/argparse v1.3.1
	github.com/BurntSushi/toml v1.1.0
)

replace badassops.ldap/utils => ../utils
