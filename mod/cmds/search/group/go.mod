module group

go 1.18

require (
	badassops.ldap/consts v0.0.0-00010101000000-000000000000
    badassops.ldap/cmds/search/common v0.0.1
	badassops.ldap/ldap v0.0.0-00010101000000-000000000000
)

replace badassops.ldap/consts => ../../../consts

replace badassops.ldap/cmds/search/common => ../common

replace badassops.ldap/ldap => ../../../ldap
