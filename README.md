## WORK IN PROGRESS!

## ldap-tool
A simple Go script to manage OpenLDAP users

### Background
The script is based on a certain LDAP dn values

### History
Orignally it was written in bash using the ldap CLI's, the script was mean to be able 
to manage OpenLDAP user, such as add, modify, delete and severel search capabilities 
Using an UI interface such a phpLDAPadmin is not always possible, and so I decide 
to build this tools

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
			 search (user), group (group)
		     get group members commands:
			 group (base group), admin (admin group)
		     get all users and members of all groups commands:
			 users, groups, admins (admin groups)

  -i  --info         Show information
  -v  --version      Show version
```
