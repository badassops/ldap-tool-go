## WORK IN PROGRESS!

## ldap-tool
A simple Go script to manage OpenLDAP users

### Background
The script is based on a certain LDAP settings
- OpenLDAP
- the use memberOf
- password length and use of special charachter is in the config file
- the config file is toml formatted

### History
Using an UI interface such a phpLDAPadmin is not always possible, and so I decide 
to build this tools.. 
Orignally it was written in bash using the ldap CLI's, the script is meant to be able 
to manage OpenLDAP user, such as add, modify, delete and several search capabilities 

### 

## Usage
```
usage: ldap-tool [-h|--help] [-c|--configFile "<value>"] [-e|--environment
                 "<value>"] [-m|--mode
                 (create|modify|delete|search|group|admin|users|groups|admins)]
                 [-i|--info] [-v|--version]

                 Simple script to manage LDAP users

Arguments:

  -h  --help         Print help information
  -c  --configFile   Path to the configuration file to be use
  -e  --environment  Server environment. Default: dev
  -m  --mode         base commands:
			 create, modify, delete
		     search commands:
			 (user) search, (group) group, admin
		     get all records users and groups commands:
			 (user) users, (group) groups, admins

  -i  --info         Show information
  -v  --version      Show version

```
