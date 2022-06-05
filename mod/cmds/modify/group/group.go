//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package modify

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"badassops.ldap/cmds/common"
	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
	ldapv3 "gopkg.in/ldap.v2"
)

var (
	valueEntered string
	modCount     int = 0
	p                = print.New()
)

// remove the user from groups its belong to if this was set during the template input
func deleteGroupEntries(c *l.Connection, groupName string) {
	for _, userID := range v.WorkRecord.GroupDelList {
		v.WorkRecord.ID = userID
		if !c.RemoveFromGroups() {
			p.PrintRed(fmt.Sprintf("\tFailed to remove the user %s from the group %s, check the log file\n",
				userID, v.WorkRecord.DN))
		} else {
			p.PrintGreen(fmt.Sprintf("\tUser %s removed from group %s\n", userID, groupName))
		}
	}
}

// add the user to groups if this was set during the template input
func addGroupEntries(c *l.Connection, groupName string) {
	for _, userID := range v.WorkRecord.GroupAddList {
		v.WorkRecord.ID = userID
		if !c.AddToGroup() {
			p.PrintRed(fmt.Sprintf("\n\tFailed to add the user %s to the group %s, check the log file\n",
				userID, v.WorkRecord.DN))
		} else {
			p.PrintGreen(fmt.Sprintf("\tUser %s added to group %s\n", userID, groupName))
		}
	}
}

// group modigy template
func modifyGroup(c *l.Connection, records *ldapv3.SearchResult) bool {
	orgGroup := v.WorkRecord.ID
	p.PrintPurple(fmt.Sprintf("\tUsing group: %s\n", orgGroup))
	p.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
	v.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", v.WorkRecord.ID, c.Config.ServerValues.GroupDN)
	for idx, entry := range records.Entries {
		if len(records.Entries[idx].GetAttributeValue("gidNumber")) != 0 {
			v.WorkRecord.MemberType = "memberUid"
			for _, member := range entry.GetAttributeValues("memberUid") {
				p.PrintCyan(fmt.Sprintf("\tmember: %s\n", member))
			}
		} else {
			v.WorkRecord.MemberType = "member"
			for _, member := range entry.GetAttributeValues("member") {
				p.PrintCyan(fmt.Sprintf("\tmember: %s\n", member))
			}
		}
	}

	p.PrintRed("\n\tEnter the user(s) to be deleted, select from the list above, (default to skip)\n")
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("\tUser : ")
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if valueEntered != "" {
			if v.WorkRecord.MemberType == "member" {
				valueEntered = fmt.Sprintf("uid=%s,%s", valueEntered, c.Config.ServerValues.UserDN)
			}
			v.WorkRecord.GroupDelList = append(v.WorkRecord.GroupDelList, valueEntered)
		} else {
			break
		}
	}

	p.PrintGreen("\n\tEnter the user(s) to be added, (default to skip)\n")
	for true {
		fmt.Printf("\tUser : ")
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if valueEntered != "" {
			if v.WorkRecord.MemberType == "member" {
				valueEntered = fmt.Sprintf("uid=%s,%s", valueEntered, c.Config.ServerValues.UserDN)
			}
			v.WorkRecord.GroupAddList = append(v.WorkRecord.GroupAddList, valueEntered)
		} else {
			break
		}
	}

	modCount = len(v.WorkRecord.GroupDelList) + len(v.WorkRecord.GroupAddList)
	if modCount == 0 {
		p.PrintBlue(fmt.Sprintf("\n\tNo change, no modification made to group %s\n", orgGroup))
		return false
	}

	if len(v.WorkRecord.GroupDelList) > 0 {
		deleteGroupEntries(c, orgGroup)
	}
	if len(v.WorkRecord.GroupAddList) > 0 {
		addGroupEntries(c, orgGroup)
	}
	return true
}

func Modify(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Modify Group", 18, true))
	v.SearchResultData.WildCardSearchBase = v.GroupWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.GroupWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.GroupDisplayFieldID
	if common.GetObjectRecord(c, true, "group") {
		modifyGroup(c, v.SearchResultData.SearchResult)
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
