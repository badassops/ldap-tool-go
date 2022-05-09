# WORK IN PROGRESS!

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
- add (user)
- modify (user)
- delete (user)
- search (user and group)

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
- once the script has been completed, the capabilities to
 create, modify and delete a group
- create binaries for OSX and Linux

### The End
Your friendly BOFH 🦄 😈          
