//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package sudo

import (
	"fmt"

	"badassops.ldap/cmds/common"
	"badassops.ldap/ldap"
	"badassops.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

func printSudo(records []*ldapv3.Entry, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	for _, entry := range records {
		funcs.P.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		for _, attributes := range entry.Attributes {
			for _, value := range attributes.Values {
				if attributes.Name != "objectClass" {
					if attributes.Name == "cn" {
						funcs.P.PrintBlue(fmt.Sprintf("\t%s : %s \n", attributes.Name, value))
					} else {
						funcs.P.PrintCyan(fmt.Sprintf("\t%s : %s \n", attributes.Name, value))
					}
				}
			}
		}
	}
	fmt.Printf("\n")
}

func Sudo(c *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search Sudo Rule", 16, true))
	vars.SearchResultData.WildCardSearchBase = vars.SudoWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.SudoWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.SudoDisplayFieldID
	if common.GetObjectRecord(c, true, "sudo rule", funcs) {
		printSudo(vars.SearchResultData.SearchResult.Entries, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
