//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"fmt"
	"strconv"
)

// get the next user UID from the ldap database
func (c *Connection) GetNextUID() int {
	var uidValue int
	startUID := c.Config.DefaultValues.UidStart
	c.SearchInfo.SearchBase = "(objectClass=person)"
	c.SearchInfo.SearchAttribute = []string{"uidNumber"}
	records, _ := c.Search()
	for _, entry := range records.Entries {
		uidValue, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if uidValue > startUID {
			startUID = uidValue
		}
	}
	return startUID + 1
}

// get the next group GID from the ldap database
func (c *Connection) GetNextGID() int {
	var uidValue int
	startGID := c.Config.DefaultValues.GidStart
	c.SearchInfo.SearchBase = "(objectClass=posixGroup)"
	c.SearchInfo.SearchAttribute = []string{"gidNumber"}
	records, _ := c.Search()
	for _, entry := range records.Entries {
		uidValue, _ = strconv.Atoi(entry.GetAttributeValue("gidNumber"))
		if uidValue == c.Config.DefaultValues.GroupId {
			// we skip this special gid
			continue
		}
		if uidValue > startGID {
			startGID = uidValue
		}
	}
	return startGID + 1
}

// get the groups an user belong to
func (c *Connection) GetUserGroups(userID, userDN string) int {
	c.SearchInfo.SearchBase =
		fmt.Sprintf("(|(&(objectClass=posixGroup)(memberUid=%s))(&(objectClass=groupOfNames)(member=%s)))",
			userID, userDN)
	c.SearchInfo.SearchAttribute = []string{"dn"}
	records, recordsCount := c.Search()
	for _, entry := range records.Entries {
		c.Record.UserGroups = append(c.Record.UserGroups, entry.DN)
	}
	return recordsCount
}

// get the group of which a user does not belong to
func (c *Connection) GetAvailableGroups(userID, userDN string) int {
	c.SearchInfo.SearchBase =
		fmt.Sprintf("(|(&(objectClass=posixGroup)(!memberUid=%s))(&(objectClass=groupOfNames)(!member=%s)))",
			userID, userDN)
	c.SearchInfo.SearchAttribute = []string{"dn"}
	records, recordsCount := c.Search()
	for _, entry := range records.Entries {
		c.Record.AvailableGroups = append(c.Record.AvailableGroups, entry.DN)
	}
	return recordsCount
}

// get all group and their type: posix or groupOfNames
func (c *Connection) GetGroupType() map[string][]string {
	result := make(map[string][]string)
	c.SearchInfo.SearchBase = "(&(objectClass=posixGroup))"
	c.SearchInfo.SearchAttribute = []string{"dn"}
	records, _ := c.Search()
	for _, posix := range records.Entries {
		result["posixGroup"] = append(result["posixGroup"], posix.DN)
	}
	c.SearchInfo.SearchBase = "(&(objectClass=groupOfNames))"
	c.SearchInfo.SearchAttribute = []string{"dn"}
	records, _ = c.Search()
	for _, groupOfNames := range records.Entries {
		result["groupOfNames"] = append(result["groupOfNames"], groupOfNames.DN)
	}
	return result
}

// get all group in the ldap database
func (c *Connection) GetAllGroups() []string {
	groups := c.GetGroupType()
	return append(groups["posixGroup"], groups["groupOfNames"]...)
}

// get all the posixGroup's group GID
func (c *Connection) GetAlGroupsGID() map[string]string {
	gitNumberList := make(map[string]string)
	c.SearchInfo.SearchBase = "(&(objectClass=posixGroup))"
	c.SearchInfo.SearchAttribute = []string{"gidNumber", "cn"}
	records, _ := c.Search()
	for _, gidNumber := range records.Entries {
		gitNumberList[gidNumber.GetAttributeValue("gidNumber")] = gidNumber.GetAttributeValue("cn")
	}
	return gitNumberList
}

// get all sudo rule in the ldap database
func (c *Connection) GetAllSudoRules() []string {
	var sudoRuleList []string
	c.SearchInfo.SearchBase = "(&(objectClass=sudoRole))"
	c.SearchInfo.SearchAttribute = []string{"gidNumber", "cn"}
	records, _ := c.Search()
	for _, sudoRule := range records.Entries  {
		sudoRuleList = append(sudoRuleList, sudoRule.GetAttributeValue("cn"))
	}
	return sudoRuleList
}
