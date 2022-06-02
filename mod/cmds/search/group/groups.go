// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package group

import (
	"fmt"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printGroups(records *ldapv3.SearchResult, recordCount int) {
	var memberCount = 0
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 55))
	for idx, entry := range records.Entries {
		p.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		p.PrintBlue(fmt.Sprintf("\tcn: %s\n",
			records.Entries[idx].GetAttributeValue("cn")))
		if len(records.Entries[idx].GetAttributeValue("gidNumber")) != 0 {
			p.PrintCyan(fmt.Sprintf("\tgidNumber: %s\n",
				records.Entries[idx].GetAttributeValue("gidNumber")))
			memberCount = 0
			for _, member := range entry.GetAttributeValues("memberUid") {
				p.PrintCyan(fmt.Sprintf("\tmemberUid: %s\n", member))
				memberCount++
			}
			p.PrintYellow(fmt.Sprintf("\tTotal members: %d : posix group\n\n", memberCount))
		} else {
			memberCount = 0
			for _, member := range entry.GetAttributeValues("member") {
				p.PrintCyan(fmt.Sprintf("\tmember: %s\n", member))
				memberCount++
			}
			p.PrintYellow(fmt.Sprintf("\tTotal members: %d : groupOfNames group\n\n", memberCount))
		}
	}
}

func Groups(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Search Groups", 20, true))
	c.SearchInfo.SearchBase = v.GroupSearchBase
	c.SearchInfo.SearchAttribute = []string{}
	printGroups(c.Search())
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
