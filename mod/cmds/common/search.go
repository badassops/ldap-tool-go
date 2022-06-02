// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
	"github.com/badassops/packages-go/readinput"
	ldapv3 "gopkg.in/ldap.v2"
)

func GetObjectRecord(c *l.Connection, firstTime bool, displayName string) bool {
	var records *ldapv3.SearchResult
	var recordCount int
	var displayFieldID = v.SearchResultData.DisplayFieldID
	var wildCardSearchBase = v.SearchResultData.WildCardSearchBase
	var recordSearchbase = v.SearchResultData.RecordSearchbase

	p := print.New()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tEnter %s name to be use: ", displayName)
	enterData, _ := reader.ReadString('\n')
	enterData = strings.TrimSuffix(enterData, "\n")

	if enterData == "" {
		p.PrintRed(fmt.Sprintf("\n\tNo %s name was given aborting...\n", displayName))
		return false
	}

	if firstTime {
		fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
		wildCard, _ := reader.ReadString('\n')
		wildCard = strings.TrimSuffix(wildCard, "\n")
		if readinput.ReadYN(wildCard, false) == true {
			enterData = "*" + enterData + "*"
			wildCardSearchBase = strings.ReplaceAll(wildCardSearchBase, "VALUE", enterData)
			c.SearchInfo.SearchBase = wildCardSearchBase
			c.SearchInfo.SearchAttribute = []string{displayFieldID}
			records, _ = c.Search()
			for idx, _ := range records.Entries {
				p.PrintBlue(fmt.Sprintf("\t%s: %s\n",
					displayFieldID,
					records.Entries[idx].GetAttributeValue(displayFieldID)))
			}
			fmt.Printf("\n\tSelect the %s from the above list:\n", displayName)
			return GetObjectRecord(c, false, displayName)
		}
	}

	recordSearchbase = strings.ReplaceAll(recordSearchbase, "VALUE", enterData)
	c.SearchInfo.SearchBase = recordSearchbase
	c.SearchInfo.SearchAttribute = []string{}
	records, recordCount = c.Search()

	if recordCount == 0 {
		p.PrintRed(fmt.Sprintf("\n\t%s %s was not found, aborting...\n", strings.Title(displayName), enterData))
		return false
	}
	v.SearchResultData.RecordCount = recordCount
	v.SearchResultData.SearchResult = records
	v.WorkRecord.ID = enterData
	return true
}
