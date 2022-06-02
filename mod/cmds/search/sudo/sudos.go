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

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printSudos(records []*ldapv3.Entry, recordCount int, protectedSudoRules []string) {
	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	for _, entry := range records {
		p.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		for _, attributes := range entry.Attributes {
			for _, value := range attributes.Values {
				if attributes.Name != "objectClass" {
					if attributes.Name == "cn" {
						if i.IsInList(protectedSudoRules, value) {
							p.PrintYellow("\t\tThis entry can not be modified\n")
						}
					}
					p.PrintCyan(fmt.Sprintf("\t%s : %s \n", attributes.Name, value))
				}
			}
		}
		fmt.Printf("\n")
	}
}

func Sudos(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Search Sudo Rules", 15, true))
	c.SearchInfo.SearchBase = v.SudoRuleSearchBase
	c.SearchInfo.SearchAttribute = []string{}
	records, recordCount := c.Search()
	printSudos(records.Entries, recordCount, c.Config.SudoValues.ExcludeSudo)
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
