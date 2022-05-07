module ldap

go 1.18

require (
	badassops.ldap/configurator v0.0.0-00010101000000-000000000000
	badassops.ldap/constants v0.0.0-00010101000000-000000000000
	badassops.ldap/initializer v0.0.0-00010101000000-000000000000
	badassops.ldap/ldap v0.0.0-00010101000000-000000000000
	badassops.ldap/logs v0.0.0-00010101000000-000000000000
	badassops.ldap/utils v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/akamensky/argparse v1.3.1 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace badassops.ldap/configurator => ./mod/configurator

replace badassops.ldap/initializer => ./mod/initializer

replace badassops.ldap/logs => ./mod/logs

replace badassops.ldap/utils => ./mod/utils

replace badassops.ldap/constants => ./mod/constants

replace badassops.ldap/ldap => ./mod/ldap

replace badassops.ldap/cmds => ./mod/cmds
