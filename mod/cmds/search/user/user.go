// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:  0.1
//

package user

import (
	"fmt"
	"strconv"

	"badassops.ldap/cmds/common"
	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printUser(c *l.Connection, records *ldapv3.SearchResult) {
	// the values are in days so we need to multiple by 86400
	value, _ := strconv.ParseInt(records.Entries[0].GetAttributeValue("shadowLastChange"), 10, 64)
	passChanged := e.ReadableEpoch(value * 86400)

	value, _ = strconv.ParseInt(records.Entries[0].GetAttributeValue("shadowExpire"), 10, 64)
	passExpired := e.ReadableEpoch(value * 86400)

	userName := records.Entries[0].GetAttributeValue("uid")
	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	p.PrintBlue(fmt.Sprintf("\tdn: %s\n", records.Entries[0].DN))
	userDN := records.Entries[0].DN
	for _, fieldName := range v.DisplayUserFields {
		p.PrintCyan(fmt.Sprintf("\t%s: %s\n", fieldName, records.Entries[0].GetAttributeValue(fieldName)))
	}

	c.SearchInfo.SearchBase =
		fmt.Sprintf("(|(&(objectClass=posixGroup)(memberUid=%s))(&(objectClass=groupOfNames)(member=%s)))",
			userName, userDN)
	c.SearchInfo.SearchAttribute = []string{"dn"}
	groupRecords, _ := c.Search()
	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	p.PrintPurple(fmt.Sprintf("\tUser %s groups:\n", userName))
	for _, entry := range groupRecords.Entries {
		p.PrintCyan(fmt.Sprintf("\tdn: %s\n", entry.DN))
	}

	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	p.PrintPurple(fmt.Sprintf("\tUser %s password information\n", userName))
	p.PrintCyan(fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	p.PrintRed(fmt.Sprintf("\tPassword will expired on %s\n\n", passExpired))
}

func User(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Search User", 18, true))
	v.SearchResultData.WildCardSearchBase = v.UserWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.UserWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.UserDisplayFieldID
	if common.GetObjectRecord(c, true, "user") {
		printUser(c, v.SearchResultData.SearchResult)
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
