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

	"crypto/tls"

	"badassops.ldap/consts"
	"badassops.ldap/vars"
	"badassops.ldap/utils"
	"badassops.ldap/configurator"
	"badassops.ldap/logs"

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

var (
	// these are the objectClasses needed for a user record
	userObjectClasses = []string{"top", "person",
		"organizationalPerson", "inetOrgPerson",
		"posixAccount", "shadowAccount", "ldapPublicKey"}

	groupObjectClasses = []string{"groupOfNames"}
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
	if config.ServerValues.TLS {
		err := ServerConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			utils.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
			utils.ExitIfError(err)
		}
	}

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
		c.User.Field["dn"] = records.Entries[0].DN
		for _, field := range vars.Fields {
			c.User.Field[field] = records.Entries[0].GetAttributeValue(field)
		}
	}
	c.userGroups()
	return cnt
}

func (c *Connection) CheckUser(user string) int {
	attributes := []string{}
	searchBase :=fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", user)
	_, cnt := c.search(searchBase, attributes)
	return cnt
}

func (c *Connection) GetGroup(group string) (int, string) {
	attributes := []string{}
	groupTypes := []string{"posix", "groupOfNames"}
	var cnt int
	var searchBase string
	for _, groupType := range groupTypes {
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
	for idx, _ := range records.Entries {
		// utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", entry.DN))
		utils.PrintColor(utils.Blue, fmt.Sprintf("\tuid: %s\n",
			records.Entries[idx].GetAttributeValue("uid")))
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
	searchBase := fmt.Sprintf("(&(objectClass=groupOfNames)(member=%s))", c.User.Field["dn"])
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

func (c *Connection) AddUser() bool {
	// adding a user record
	newUserRecord := ldapv3.NewAddRequest(c.User.Field["dn"])
	newUserRecord.Attribute("objectClass" ,userObjectClasses)
	for _, field := range vars.Fields {
		if field != "groups" {
			newUserRecord.Attribute(field, []string{c.User.Field[field]})
		}
	}
	err := c.Conn.Add(newUserRecord)
	if err != nil {
		msg := fmt.Sprintf("Error adding user %s, Error: %s",
				c.User.Field["uid"], err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg := fmt.Sprintf("User %s has been added", c.User.Field["uid"])
	logs.Log(msg, "INFO")
	// once the record is create we need to hash the password
	passwordModifyRequest := ldapv3.NewPasswordModifyRequest(
		c.User.Field["dn"],
		c.User.Field["userPassword"],
		c.User.Field["userPassword"])
	_, err = c.Conn.PasswordModify(passwordModifyRequest)
	if err != nil {
		msg := fmt.Sprintf("Error set the password for the user %s, Error: %s",
				c.User.Field["uid"], err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	return true
}

func (c *Connection) AddMember() bool {
	// adding the user to the groups
	var errorCnt int = 0
	for _, group := range c.User.Groups {
		groupCN := fmt.Sprintf("cn=%s,%s", group, c.Config.ServerValues.GroupDN)
		addNewMember := ldapv3.NewModifyRequest(groupCN)
		addNewMember.Add("member", []string{c.User.Field["dn"]})
		err := c.Conn.Modify(addNewMember)
		if err != nil {
			msg := fmt.Sprintf("Error adding user %s to group %s, Error: %s",
					c.User.Field["dn"], group, err.Error())
			logs.Log(msg, "ERROR")
			errorCnt++
		}
	}
	if errorCnt != 0 {
		return false
	}
	return true
}

func (c *Connection) AddRecord() bool {
	state := c.AddUser()
	if state == false {
		return false
	}
	return c.AddMember()
}

func (c *Connection) SearchUsersGroups(user string) []string {
	var groupsList []string
	searchBase := fmt.Sprintf("(&(objectClass=posixGroup))")
	attributes := []string{"cn", "memberUid"}
	records, _ := c.search(searchBase, attributes)
	for idx, entry := range records.Entries {
		for _, member := range entry.GetAttributeValues("memberUid") {
			if member == "user" {
				groupsList = append(groupsList, records.Entries[idx].GetAttributeValue("cn"))
			}
		}
	}
	return groupsList
}

func (c *Connection) RemoveFromGroups(groups []string, user string) {
	// errors are logged and the process continues
	for _, group := range groups {
		groupCN := fmt.Sprintf("cn=%s,%s", group, c.Config.ServerValues.GroupDN)
		removeMember := ldapv3.NewModifyRequest(groupCN)
		removeMember.Delete("memberUid", []string{user})
		err := c.Conn.Modify(removeMember)
		if err != nil {
			msg := fmt.Sprintf("Error remove user %s from the group %s, Error: %s",
				user, group, err.Error())
			logs.Log(msg, "ERROR")
		} else {
			msg := fmt.Sprintf("User %s has been removed from the group %s", user, group)
			logs.Log(msg, "INFO")
		}
	}
}

func (c *Connection) DeleteRecord() bool {
	delReq := ldapv3.NewDelRequest(c.User.Field["dn"], []ldapv3.Control{})
	if err := c.Conn.Del(delReq); err != nil {
		msg := fmt.Sprintf("Error deleting the user %s, error %s",
			c.User.Field["uid"], err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	inGroups := c.SearchUsersGroups(c.User.Field["uid"])
	c.RemoveFromGroups(inGroups, c.User.Field["uid"])
	msg := fmt.Sprintf("The user %s has been deleted", c.User.Field["uid"])
	logs.Log(msg, "INFO")
	return true
}

func (c *Connection) AddGroup(groupName string, groupType string, groupID int) bool {
	groupCN := fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
	newGroupRecord := ldapv3.NewAddRequest(groupCN)
	switch groupType {
		case "posix", "p", "":
			newGroupRecord.Attribute("objectClass", []string{"posixGroup"})
			newGroupRecord.Attribute("cn", []string{groupName})
			newGroupRecord.Attribute("gidNumber", []string{strconv.Itoa(groupID)})
		case "groupOfNames", "g":
			newGroupRecord.Attribute("objectClass", []string{"groupOfNames"})
			newGroupRecord.Attribute("cn", []string{groupName})
			newGroupRecord.Attribute("member", []string{"uid=initial-member,ou=users,dc=co,dc=badassops,dc=com"})
	}
	err := c.Conn.Add(newGroupRecord)
	if err != nil {
		msg := fmt.Sprintf("Error creating new group %s, Error: %s",
			groupName, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg := fmt.Sprintf("New group %s created", groupName)
	logs.Log(msg, "INFO")
	return true
}

func (c *Connection) DeleteGroup(groupName string) bool {
	groupCN := fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
	delReq := ldapv3.NewDelRequest(groupCN, []ldapv3.Control{})
	if err := c.Conn.Del(delReq); err != nil {
		msg := fmt.Sprintf("Error deleting the group %s, error %s",
			groupName, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg := fmt.Sprintf("The group %s has been deleted", groupName)
	logs.Log(msg, "INFO")
	return true
}

func (c *Connection) GetGroupType(groupName string) (string, bool) {
	cnt, typeGroup := c.GetGroup(groupName)
	if cnt != 0 {
		return typeGroup, true
	}
	return "errored", false
}

func (c *Connection) ModifyGroup(groupName string, addUsers []string, delUsers []string) (bool, int) {
	var memberField string
	var changes int = 0
	groupType, state := c.GetGroupType(groupName)
	if !state {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tUnable to determinated the group type, aboring...\n"))
		return false, changes
	}
	groupCN := fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
	modifyMember := ldapv3.NewModifyRequest(groupCN)
	for _, addUser:= range addUsers {
		if groupType == "groupOfNames" {
			addUser = fmt.Sprintf("uid=%s,%s", addUser, c.Config.ServerValues.UserDN)
			memberField = "member"
		}
		if groupType == "posix" {
			memberField = "memberUid"
		}
		modifyMember.Add(memberField, []string{addUser})
		changes++
	}
	for _, delUser := range delUsers {
		if groupType == "groupOfNames" {
			delUser = fmt.Sprintf("uid=%s,%s", delUser, c.Config.ServerValues.UserDN)
			memberField = "member"
		}
		if groupType == "posix" {
			memberField = "memberUid"
		}
		modifyMember.Delete(memberField, []string{delUser})
		changes++
	}

	err := c.Conn.Modify(modifyMember)
	if err != nil {
		msg := fmt.Sprintf("Error modifying the group group %s, Error: %s",
				groupName, err.Error())
		logs.Log(msg, "ERROR")
		return false, 0
	}
	return true, changes
}
