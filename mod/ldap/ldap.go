// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package ldap

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"badassops.ldap/vars"
	"badassops.ldap/utils"
	"badassops.ldap/configurator"

	ldapv3 "gopkg.in/ldap.v2"
)

type (
	Connection struct {
		Conn		*ldapv3.Conn
		User		vars.UserRecord
		Config		*configurator.Config
		LockFile	string
		LockPid		int
	}
)

// function to initialize a user record
func New(config *configurator.Config) *Connection {
	// set variable for the ldap connection
	var ppolicy *ldapv3.ControlBeheraPasswordPolicy

	// check if we can search the server, timeout set to 15 seconds
	timeout := 15 * time.Second
	dialConn, err := net.DialTimeout("tcp", net.JoinHostPort(config.ServerValues.Server, "389"), timeout)
	if err != nil {
		utils.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
		utils.ExitWithMesssage(err.Error())
	}
	dialConn.Close()

	ServerConn, err := ldapv3.Dial("tcp", fmt.Sprintf("%s:%d", config.ServerValues.Server ,389))
	if err != nil {
		utils.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
		utils.ExitWithMesssage(err.Error())
	}

	// now we need to reconnect with TLS
	// if config.TLS {
	// 	err := ServerConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// 	if err != nil {
	//		utils.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
	// 		utils.ExitIfError(err)
	// 	}
	// }

	// setup control
	controls := []ldapv3.Control{}
	controls = append(controls, ldapv3.NewControlBeheraPasswordPolicy())

	// bind to the ldap server
	bindRequest := ldapv3.NewSimpleBindRequest(config.ServerValues.Admin, config.ServerValues.AdminPass, controls)
	request, err := ServerConn.SimpleBind(bindRequest)
	ppolicyControl := ldapv3.FindControl(request.Controls, ldapv3.ControlTypeBeheraPasswordPolicy)
	if ppolicyControl != nil {
		ppolicy = ppolicyControl.(*ldapv3.ControlBeheraPasswordPolicy)
	 }
	if err != nil {
		errStr := "ERROR: Cannot bind: " + err.Error()
		if ppolicy != nil && ppolicy.Error >= 0 {
			errStr += ":" + ppolicy.ErrorString
		}
		utils.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
		utils.ExitWithMesssage(errStr)
	}

	// the rest of the values will be filled during the process
    return &Connection {
		Conn:		ServerConn,
		Config:		config,
		User:		vars.User,
		LockFile:	config.DefaultValues.LockFile,
        LockPid:	config.LockPID,
	}
}

func (c *Connection) GetUser(user string) int {
	attributes := []string{}
	searchBase :=fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", user)
	records, cnt := c.search(searchBase, attributes)
	if cnt == 1 {
		c.User.DN = records.Entries[0].DN
		for _, field := range vars.Fields{
			switch field {
				case "shadowWarning", "shadowMax", "uidNumber", "gidNumber":
					value, _ := strconv.Atoi(records.Entries[0].GetAttributeValue(field))
					c.User.Ints[field] = vars.IntRecord{Value: value  , Changed: false}
				default:
					c.User.Strings[field] =
						vars.StringRecord{Value: records.Entries[0].GetAttributeValue(field) , Changed: false}
			}
		}
		c.userGroups()
	}
	return cnt
}

func (c *Connection) GetGroup(group string) (int, string) {
	attributes := []string{}
	groupTypes := []string{"posix", "groupOfNames"}
	var cnt int
	var searchBase string
	for cnt, groupType := range groupTypes {
		switch groupType {
			case "posix":
				searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", group)
			case "groupOfNames":
				searchBase = fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", group)
		}
		_, cnt = c.search(searchBase, attributes)
		if cnt > 0 {
			return cnt, groupType
		}
	}
	return cnt, "nofound"
}

func (c *Connection) SearchUser(user string) int {
	attributes := []string{}
	searchBase := fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", user)
	records, cnt := c.search(searchBase, attributes)
	if cnt == 0 {
		return 0
	}
	for _, entry := range records.Entries {
		utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", entry.DN))
	}
	return cnt
}

func (c *Connection) SearchUsers(baseInfo bool) int {
	attributes := []string{}
	searchBase := fmt.Sprintf("(objectClass=person)")
	records, cnt := c.search(searchBase, attributes)
	if cnt == 0 {
		return 0
	}
	for idx, entry := range records.Entries {
		utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", entry.DN))
		if baseInfo {
			userBaseInfo := fmt.Sprintf("\tFull namae: %s %s\n\t\tdepartmentNumber %s",
				records.Entries[idx].GetAttributeValue("givenName"),
				records.Entries[idx].GetAttributeValue("sn"),
				records.Entries[idx].GetAttributeValue("departmentNumber"))
			utils.PrintColor(utils.Cyan, userBaseInfo)
			fmt.Printf("\n")
		}
	}
	utils.PrintColor(utils.Yellow, fmt.Sprintf("\n\tTotal records: %d \n", cnt))
	return cnt
}

func (c *Connection) SearchGroup(group string, baseInfo bool) int {
	var searchBase string
	var memberField string
	var cnt int
	var groupType string
	attributes := []string{}
	cnt, groupType = c.GetGroup(group)
	if cnt == 0 {
		return 0
	}
	switch groupType{
		case "posix":
			searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", group)
			memberField = "memberUid"
		case "groupOfNames":
			searchBase = fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", group)
			memberField = "member"
	}
	records, cnt := c.search(searchBase, attributes)
	fmt.Printf("\n")
	for idx, entry := range records.Entries {
		if baseInfo {
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", entry.DN))
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tcn: %s\n",
				records.Entries[idx].GetAttributeValue("cn")))
			if groupType == "posix" {
				utils.PrintColor(utils.Cyan,
					fmt.Sprintf("\tgidNumber: %s\n", entry.GetAttributeValue("gidNumber")))
			}
			for _, member := range entry.GetAttributeValues(memberField) {
				utils.PrintColor(utils.Cyan, fmt.Sprintf("\t%s: %s\n", memberField, member))
			}
		} else {
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tcn: %s\n",
				records.Entries[idx].GetAttributeValue("cn")))
		}
	}
	if ! baseInfo {
		utils.PrintColor(utils.Yellow, fmt.Sprintf("\n\tTotal records: %d \n", cnt))
	}
	return cnt
}

func (c *Connection) SearchGroups() {
	var searchBase string
	var memberField string
	attributes := []string{}
	groupTypes := []string{"posix", "groupOfNames"}
	group := "*"
	for _, groupType := range groupTypes {
		switch groupType{
			case "posix":
				searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", group)
				memberField = "memberUid"
			case "groupOfNames":
				searchBase = fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", group)
				memberField = "member"
		}
		records, cnt := c.search(searchBase, attributes)
		fmt.Printf("\n")
		for idx, entry := range records.Entries {
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", entry.DN))
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tcn: %s\n",
				records.Entries[idx].GetAttributeValue("cn")))
			if groupType == "posix" {
				utils.PrintColor(utils.Cyan,
				fmt.Sprintf("\tgidNumber: %s\n", entry.GetAttributeValue("gidNumber")))
			}
			for _, member := range entry.GetAttributeValues(memberField) {
				utils.PrintColor(utils.Cyan, fmt.Sprintf("\t%s: %s\n", memberField, member))
			}
			utils.PrintColor(utils.Yellow,
					fmt.Sprintf("\tTotal members: %d \n", len(entry.GetAttributeValues(memberField))))
			fmt.Printf("\n")
		}
		utils.PrintColor(utils.Yellow, fmt.Sprintf("\n\tTotal %s groups: %d \n", groupType, cnt))
	}
}

func (c *Connection) userGroups() {
	searchBase := fmt.Sprintf("(&(objectClass=groupOfNames)(member=%s))", c.User.DN)
	records, cnt := c.search(searchBase, []string{"dn"})
	if cnt > 0 {
		for _, entry := range records.Entries {
			if utils.InList(c.User.Groups, entry.DN) == false {
				c.User.Groups = append(c.User.Groups, entry.DN)
			}
		}
	}
}

func (c *Connection) GetNextUID() int {
	var highestUID = 0
	var uidValue int
	attributes := []string{"uidNumber"}
	searchBase := fmt.Sprintf("(objectClass=person)")
	records, _ := c.search(searchBase, attributes)
	for _, entry := range records.Entries {
			uidValue, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
			if uidValue > highestUID {
				highestUID = uidValue
			}
	}
	return highestUID + 1
}

func (c *Connection) search(searchBase string, searchAttribute []string) (*ldapv3.SearchResult, int) {
	searchRecords := ldapv3.NewSearchRequest(
		c.Config.ServerValues.BaseDN,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases, 0, 0, false,
		searchBase,
		searchAttribute,
		nil,
	)
	sr, err := c.Conn.Search(searchRecords)
	if err != nil {
		c.Conn.Close()
		utils.ReleaseIT(c.LockFile, c.LockPid)
		utils.ExitWithMesssage(err.Error())
	}
	if len(sr.Entries) > 0 {
		return sr, len(sr.Entries)
	}
	return sr, 0
}
