//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package create

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"

	"github.com/badassops/packages-go/is"
	"github.com/badassops/packages-go/print"
)

var (
	i = is.New()
	p = print.New()
)

func createSudoRecord(c *l.Connection) bool {
	sudoRules := c.GetAllSudoRules()

	for _, fieldName := range v.SudoFields {
		if v.Template[fieldName].Value != "" {
			fmt.Printf("\t%sDefault to:%s %s%s%s\n",
				v.Purple, v.Off, v.Cyan, v.Template[fieldName].Value, v.Off)
		}

		fmt.Printf("\t%s: ", v.Template[fieldName].Prompt)

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		// make sure any combination of `all` is made uppercase
		if strings.ToLower(valueEntered) == "all" {
			valueEntered = "ALL"
		}

		switch fieldName {
		case "cn":
			if i.IsInList(sudoRules, valueEntered) {
				p.PrintRed(fmt.Sprintf("\n\tGiven cn %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
			p.PrintPurple(fmt.Sprintf("\tUsing Sudo Rule: %s\n\n", valueEntered))
			v.WorkRecord.Fields["cn"] = valueEntered
			v.WorkRecord.Fields["dn"] = fmt.Sprintf("cn=%s,%s", valueEntered, c.Config.SudoValues.SudoersBase)
			v.WorkRecord.Fields["objectClass"] = "sudoRole"

		case "sudoCommand", "sudoHost", "sudoRunAsUser":
			if len(valueEntered) > 0 {
				v.WorkRecord.Fields[fieldName] = valueEntered
			} else {
				v.WorkRecord.Fields[fieldName] = v.Template[fieldName].Value
			}
		case "sudoOption":
			if len(valueEntered) > 0 {
				v.WorkRecord.Fields[fieldName] = valueEntered
			}
		case "sudoOrder":
			v.WorkRecord.Fields[fieldName] = v.Template[fieldName].Value
			if len(valueEntered) > 0 {
				value, _ := strconv.Atoi(valueEntered)
				if value < 3 || value > 10 {
					p.PrintRed(fmt.Sprintf("%s\tGiven order %s is not allowed, set to the default %s\n",
						v.OneLineUP, valueEntered, v.Template[fieldName].Value))
					valueEntered = v.Template[fieldName].Value
				} else {
					v.WorkRecord.Fields[fieldName] = valueEntered
				}
			}
		}
		if len(valueEntered) == 0 && v.Template[fieldName].NoEmpty == true {
			p.PrintRed("\tNo value was entered aborting...\n\n")
			return false
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
	return c.CreateSudoRule()
}

func Create(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Create Group", 18, true))
	if createSudoRecord(c) {
		p.PrintGreen(fmt.Sprintf("\tSudo rule %s created\n", v.WorkRecord.Fields["cn"]))
	} else {
		p.PrintRed(fmt.Sprintf("\tFailed to create the sudo rule %s, check the log file\n",
			v.WorkRecord.Fields["cn"]))
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
