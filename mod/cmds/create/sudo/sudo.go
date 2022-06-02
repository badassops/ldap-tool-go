// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package create

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	l "badassops.ldap/ldap"
	u "badassops.ldap/utils"
	v "badassops.ldap/vars"
)

var ()

func createSudoRecord(c *l.Connection) bool {
	for _, fieldName := range v.Sudoers {
		fmt.Printf("\t%s: ", v.SudoTemplate[fieldName].Prompt)

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		if len(valueEntered) == 0 && v.SudoTemplate[fieldName].NoEmpty == true {
			u.PrintRed("\tNo value was entered aborting...\n\n")
			return false
		}
		if len(valueEntered) == 0 && v.SudoTemplate[fieldName].UseValue == true {
			valueEntered = v.SudoTemplate[fieldName].Value
		}
		fmt.Printf("\n")
		switch fieldName {
		case "cn":
			if c.SearchSudoCN(valueEntered, false) == 1 {
				u.PrintRed(fmt.Sprintf("\n\tGiven cn %s already exist, aborting...\n\n", valueEntered))
				return false
			}
		case "sudoCommand":
			// make sure any combination of all is made uppercase
			if strings.ToLower(valueEntered) == "all" {
				valueEntered = "ALL"
			}
		case "sudoHost":
		case "sudoOption":
		case "sudoOrder":
			value, _ := strconv.Atoi(valueEntered)
			if value < 3 || value > 10 {
				u.PrintRed(fmt.Sprintf("%s\tGiven order %s is not allowed, set to the default %s\n\n",
					u.OneLineUP, valueEntered, v.SudoTemplate[fieldName].Value))
			}
		case "sudoRunAsUser":
		}
		if len(valueEntered) != 0 {
			v.ModRecord.Field[fieldName] = valueEntered
		}
	}
	v.ModRecord.Field["dn"] = fmt.Sprintf("cn=%s,ou=%s",
		v.ModRecord.Field["cn"], c.Config.SudoValues.SudoersBase)
	return true
}

func Create(c *l.Connection) {
	u.PrintHeader(u.Purple, "Create Sudo Rule", true)
	if createSudoRecord(c) {
		if !c.AddSudoRule() {
			u.PrintRed(fmt.Sprintf("\n\tFailed adding the sudo rule %s, check the log file\n", v.ModRecord.Field["cn"]))
		} else {
			u.PrintGreen(fmt.Sprintf("\n\tSudo rule %s added successfully\n", v.ModRecord.Field["cn"]))
		}
	}
	u.PrintLine(u.Purple)
}
