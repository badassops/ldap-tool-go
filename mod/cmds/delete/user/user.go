//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package delete

import (
	"fmt"

	"badassops.ldap/cmds/common"
	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
)

var (
	p = print.New()
)

// once an user has been deleted, we need to make sure
// it are removed from all the group its belong too
func removeUserFromGroups(c *l.Connection) {
	var groupsList []string
	userUID := v.WorkRecord.ID
	c.SearchInfo.SearchBase = "(&(objectClass=posixGroup))"
	c.SearchInfo.SearchAttribute = []string{"cn", "memberUid"}

	records, _ := c.Search()
	for idx, entry := range records.Entries {
		for _, member := range entry.GetAttributeValues("memberUid") {
			if member == userUID {
				groupsList = append(groupsList, records.Entries[idx].GetAttributeValue("cn"))
			}
		}
	}
	if len(groupsList) > 0 {
		for _, groupName := range groupsList {
			v.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
			if !c.RemoveFromGroups() {
				p.PrintRed(fmt.Sprintf("User % was not remobe from the group %s, check the log...\n",
					userUID, groupName))
			}
		}
	}
}

func Delete(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Delete User", 18, true))
	v.SearchResultData.WildCardSearchBase = v.UserWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.UserWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.UserDisplayFieldID
	// we only handle posix group
	v.WorkRecord.GroupType = "posix"
	v.WorkRecord.MemberType = "memberUid"
	if common.GetObjectRecord(c, true, "user") {
		common.DeleteObjectRecord(c, v.SearchResultData.SearchResult, "user")
		removeUserFromGroups(c)
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
