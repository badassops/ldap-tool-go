module ldap

go 1.18

require (
  badassops.ldap/utils v0.0.0-00010101000000-000000000000
  badassops.ldap/configurator v0.0.0-00010101000000-000000000000
)

replace badassops.ldap/utils => ../utils
replace badassops.ldap/configurator => ../configurator
