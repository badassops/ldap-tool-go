# WORK IN PROGRESS!

## Progress
- search completed (May 9, 2022)
- create user completed (May 15, 2020)
- create group completed (May 15, 2020)
- delete user completed (May 15, 2020)
- delete group completed (May 15, 2020)
- modify group completed (May 16, 2020)

TODO: (in order)
- modify (user)

### ldap-tool
A simple Go script to manage OpenLDAP users

### Background
The script is based on a certain LDAP settings
- OpenLDAP
- the memberOf ldap plugin
- the SSH schema
- the SUDO schema
- password length and use of special charachter is in the config file
- the config file is toml formatted
- Runs on OSX or Linux (Ubuntu 20.04 or newer)

### History
Using an UI interface such a phpLDAPadmin is not always possible, and so I decide 
to build this tools. Orignally it was written in **bash** using the ldap CLI's and then
in **Python**, but it was not always working on someone else's laptop do the specific command
(*bash*) or module (*Python*), so I decided to write it in Go

### Capabilities
The script is meant to be able to manage OpenLDAP user such:
- add user and group
- modify user and group
- delete user and group
- search user and group

## Usage
```
usage: ldap-tool [-h|--help] [-c|--configFile "<value>"]
                 [-e|--environment "<value>"] [-m|--mode
                 (create|modify|delete|search)] [-d|--debug]
                 [-i|--info] [-v|--version]

                 Simple script to manage LDAP users

Arguments:

  -h  --help         Print help information
  -c  --configFile   Path to the configuration file to be use
  -e  --environment  Server environment. Default: dev
  -m  --mode         base commands:
                       create, modify, delete
                     search: (U)ser, (A)ll Users, (G)roup
                        and All Group(S)
  -d  --debug        Enable debug. Default: false
  -i  --info         Show information
  -v  --version      Show version
```

### Build or run the code the code
To build the code into a single binary as simple as
```
go build ldap-tools.go
```
If everything is well, then this will produce a binary called **ldap-tools** 

To run the code
```
go run ldap-tools.go -c <your-config-file> -m <mode>
```

### TODO
 - redis for delete of user only
	- start : add key in redis
	- end   : delete key in redis
	- safe guard : run a crontab that will key(s)
		then is for evry key, delete the user in ldap
		so in case one loses conneciom, we have a safe guard
		that the user will be deleted
	- if redis it not available write the key on disk

#### IDEA:
 - an user delete steps, change password and then delete?

 and per request 👻

### The End
Your friendly BOFH 🦄 😈          
