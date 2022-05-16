module ldap

go 1.18

require (
	badassops.ldap/configurator v0.0.0-00010101000000-000000000000
	badassops.ldap/consts v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/create/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/create/user v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/delete/user v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/delete/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/modify/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/modify/user v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/search/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/search/user v0.0.0-00010101000000-000000000000
	badassops.ldap/initializer v0.0.0-00010101000000-000000000000
	badassops.ldap/ldap v0.0.0-00010101000000-000000000000
	badassops.ldap/logs v0.0.0-00010101000000-000000000000
	badassops.ldap/utils v0.0.0-00010101000000-000000000000

)

require (
	badassops.ldap/vars v0.0.0-00010101000000-000000000000 // indirect
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/akamensky/argparse v1.3.1 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace badassops.ldap/configurator => ./mod/configurator

replace badassops.ldap/consts => ./mod/consts

replace badassops.ldap/cmds/common/group => ./mod/cmds/common/group

replace badassops.ldap/cmds/common/user => ./mod/cmds/common/user

replace badassops.ldap/cmds/create/group => ./mod/cmds/create/group

replace badassops.ldap/cmds/create/user => ./mod/cmds/create/user

replace badassops.ldap/cmds/delete/user => ./mod/cmds/delete/user

replace badassops.ldap/cmds/delete/group => ./mod/cmds/delete/group

replace badassops.ldap/cmds/modify/user => ./mod/cmds/modify/user

replace badassops.ldap/cmds/modify/group => ./mod/cmds/modify/group

replace badassops.ldap/cmds/search/user => ./mod/cmds/search/user

replace badassops.ldap/cmds/search/group => ./mod/cmds/search/group

replace badassops.ldap/initializer => ./mod/initializer

replace badassops.ldap/ldap => ./mod/ldap

replace badassops.ldap/logs => ./mod/logs

replace badassops.ldap/vars => ./mod/vars

replace badassops.ldap/utils => ./mod/utils
