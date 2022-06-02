// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:  0.1
//

package user

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	r "github.com/badassops/packages-go/readinput"
	ldapv3 "gopkg.in/ldap.v2"
)

func printUsers(records *ldapv3.SearchResult, recordCount int) {
	baseInfo := false
	fmt.Printf("\tPrint full name and department (default to N)? [y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	valueEntered, _ := reader.ReadString('\n')
	valueEntered = strings.TrimSuffix(valueEntered, "\n")
	if r.ReadYN(valueEntered, false) == true {
		baseInfo = true
	}

	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 55))
	for idx, entry := range records.Entries {
		p.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		if baseInfo {
			userBaseInfo := fmt.Sprintf("\tFull namae: %s %s\t\tdepartmentNumber %s\n\n",
				records.Entries[idx].GetAttributeValue("givenName"),
				records.Entries[idx].GetAttributeValue("sn"),
				records.Entries[idx].GetAttributeValue("departmentNumber"))
			p.PrintCyan(userBaseInfo)
		}
	}
	p.PrintYellow(fmt.Sprintf("\n\tTotal records: %d \n", recordCount))
}

func Users(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Search Users", 20, true))
	c.SearchInfo.SearchBase = v.UserSearchBase
	c.SearchInfo.SearchAttribute = []string{}
	printUsers(c.Search())
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
