// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package ldap

import (
	"strconv"
)

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
