// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:  0.1
//

package sudo

import (
	"fmt"

	"badassops.ldap/cmds/common"
	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printSudo(records []*ldapv3.Entry) {
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
	for _, entry := range records {
		p.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		for _, attributes := range entry.Attributes {
			for _, value := range attributes.Values {
				if attributes.Name != "objectClass" {
					if attributes.Name == "cn" {
						p.PrintBlue(fmt.Sprintf("\t%s : %s \n", attributes.Name, value))
					} else {
						p.PrintCyan(fmt.Sprintf("\t%s : %s \n", attributes.Name, value))
					}
				}
			}
		}
	}
	fmt.Printf("\n")
}

func Sudo(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Search Sudo Rule", 16, true))
	v.SearchResultData.WildCardSearchBase = v.SudoWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.SudoWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.SudoDisplayFieldID
	if common.GetObjectRecord(c, true, "sudo rule") {
		printSudo(v.SearchResultData.SearchResult.Entries)
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
