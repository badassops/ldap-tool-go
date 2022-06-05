module main

go 1.18

require (
	badassops.ldap/cmds/common v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/create/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/create/menu v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/create/sudo v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/create/user v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/delete/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/delete/menu v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/delete/sudo v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/delete/user v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/limit v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/modify/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/modify/menu v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/modify/sudo v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/modify/user v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/search/group v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/search/menu v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/search/sudo v0.0.0-00010101000000-000000000000
	badassops.ldap/cmds/search/user v0.0.0-00010101000000-000000000000
	badassops.ldap/configurator v0.0.0-00010101000000-000000000000
	badassops.ldap/initializer v0.0.0-00010101000000-000000000000
	badassops.ldap/ldap v0.0.0-00010101000000-000000000000
	badassops.ldap/logs v0.0.0-00010101000000-000000000000
	badassops.ldap/vars v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/akamensky/argparse v1.3.1 // indirect
	github.com/badassops/packages-go/caller v0.0.0-20220530220720-227ab8c06333 // indirect
	github.com/badassops/packages-go/epoch v0.0.0-20220530190021-17555612d52b // indirect
	github.com/badassops/packages-go/exit v0.0.0-20220530220720-227ab8c06333 // indirect
	github.com/badassops/packages-go/is v0.0.0-20220530213221-2e3686fab2d7 // indirect
	github.com/badassops/packages-go/lock v0.0.0-20220530213221-2e3686fab2d7 // indirect
	github.com/badassops/packages-go/print v0.0.0-20220530213221-2e3686fab2d7 // indirect
	github.com/badassops/packages-go/random v0.0.0-20220530220720-227ab8c06333 // indirect
	github.com/badassops/packages-go/readinput v0.0.0-20220530220720-227ab8c06333 // indirect
	github.com/badassops/packages-go/spinner v0.0.0-20220530213221-2e3686fab2d7 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace badassops.ldap/configurator => ./mod/configurator

replace badassops.ldap/initializer => ./mod/initializer

replace badassops.ldap/ldap => ./mod/ldap

replace badassops.ldap/logs => ./mod/logs

replace badassops.ldap/vars => ./mod/vars

replace badassops.ldap/cmds/create/menu => ./mod/cmds/create/menu

replace badassops.ldap/cmds/create/group => ./mod/cmds/create/group

replace badassops.ldap/cmds/create/sudo => ./mod/cmds/create/sudo

replace badassops.ldap/cmds/create/user => ./mod/cmds/create/user

replace badassops.ldap/cmds/delete/menu => ./mod/cmds/delete/menu

replace badassops.ldap/cmds/delete/group => ./mod/cmds/delete/group

replace badassops.ldap/cmds/delete/sudo => ./mod/cmds/delete/sudo

replace badassops.ldap/cmds/delete/user => ./mod/cmds/delete/user

replace badassops.ldap/cmds/modify/menu => ./mod/cmds/modify/menu

replace badassops.ldap/cmds/modify/group => ./mod/cmds/modify/group

replace badassops.ldap/cmds/modify/sudo => ./mod/cmds/modify/sudo

replace badassops.ldap/cmds/modify/user => ./mod/cmds/modify/user

replace badassops.ldap/cmds/search/menu => ./mod/cmds/search/menu

replace badassops.ldap/cmds/search/group => ./mod/cmds/search/group

replace badassops.ldap/cmds/search/sudo => ./mod/cmds/search/sudo

replace badassops.ldap/cmds/search/user => ./mod/cmds/search/user

replace badassops.ldap/cmds/limit => ./mod/cmds/limit

replace badassops.ldap/cmds/common => ./mod/cmds/common
