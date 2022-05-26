module sudo

go 1.18

require (
    badassops.ldap/cmds/common/sudo v0.0.1
)

replace badassops.ldap/cmds/common/sudo => ../../common/sudo
