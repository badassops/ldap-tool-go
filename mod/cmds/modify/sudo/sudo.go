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
	allowedField     = []string{"sudoCommand", "sudoHost", "sudoOption",
		"sudoOrder", "sudoRunAsUser"}
)

func modifySudo(c *l.Connection, records []*ldapv3.Entry) bool {
	v.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", v.WorkRecord.ID, c.Config.SudoValues.SudoersBase)
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
	reader := bufio.NewReader(os.Stdin)
	for _, entry := range records {
		p.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		for _, attributes := range entry.Attributes {
			for _, value := range attributes.Values {
				if (attributes.Name != "objectClass") && (attributes.Name != "cn") {
					p.PrintCyan(fmt.Sprintf("\tField: %s%s%s\n", v.Red, attributes.Name, v.Off))
					switch attributes.Name {
					case "sudoCommand":
						p.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoHost":
						p.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoOption":
						p.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoOrder":
						p.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoRunAsUser":
						p.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))
					}
					fmt.Printf("\tEnter %sdelete%s to delete or press enter to keep: ", v.RedUnderline, v.Off)
					valueEntered, _ = reader.ReadString('\n')
					valueEntered = strings.TrimSuffix(valueEntered, "\n")
					if valueEntered == "delete" {
						v.WorkRecord.SudoDelList[attributes.Name] =
							append(v.WorkRecord.SudoDelList[attributes.Name], value)
					} else {
						fmt.Printf("\n")
					}
				}
			}
		}
	}

	p.PrintCyan(fmt.Sprintf("\n\tEach field can have multiple entries\n"))
	p.PrintCyan(fmt.Sprintf("\tPress enter to skip, or enter value for field\n"))
	for _, fieldname := range allowedField {
		for true {
			p.PrintGreen(fmt.Sprintf("\tField %s%s%s: enter value: ", v.Purple, fieldname, v.Off))
			reader := bufio.NewReader(os.Stdin)
			valueEntered, _ := reader.ReadString('\n')
			valueEntered = strings.TrimSuffix(valueEntered, "\n")
			if len(valueEntered) != 0 {
				v.WorkRecord.SudoAddList[fieldname] =
					append(v.WorkRecord.SudoAddList[fieldname], valueEntered)
			} else {
				fmt.Printf("\n")
				break
			}
		}
	}

	modCount = len(v.WorkRecord.SudoDelList) + len(v.WorkRecord.SudoAddList)
	if modCount == 0 {
		p.PrintBlue(fmt.Sprintf("\n\tNo change, no modification made to sudo rule %s\n", v.WorkRecord.ID))
		return false
	}

	if len(v.WorkRecord.SudoDelList) > 0 {
		c.DeleteSudoRule()
	}
	if len(v.WorkRecord.SudoAddList) > 0 {
		c.AddSudoRule()
	}
	return true
}

func Modify(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Sudo Rules", 18, true))
	v.SearchResultData.WildCardSearchBase = v.SudoWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.SudoWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.SudoDisplayFieldID
	if common.GetObjectRecord(c, true, "sudo rule") {
		modifySudo(c, v.SearchResultData.SearchResult.Entries)
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
