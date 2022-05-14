// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package create

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"

	"badassops.ldap/consts"
	"badassops.ldap/vars"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
)

func createRecord(conn *ldap.Connection) bool {
	var fieldName string
	var email string
	var nextUID int
	var passWord string
	var shells string
	var departments string

	for idx :=0 ; idx < len(vars.RecordFields) ; idx++ {
		fieldName = vars.RecordFields[idx].FieldName
		switch vars.RecordFields[idx].FieldName {
			case "mail":
				email = fmt.Sprintf("%s.%s@%s",
					strings.ToLower(conn.User.Strings["givenName"].Value),
					strings.ToLower(conn.User.Strings["sn"].Value),
					conn.Config.ServerValues.EmailDomain)
				utils.PrintColor(consts.Cyan, fmt.Sprintf("\tDefault email: %s\n", email))
			case "uidNumber":
				nextUID = conn.GetNextUID()
				utils.PrintColor(consts.Purple,
					fmt.Sprintf("\t\tOptional set user UID, press enter to use the next UID: %d\n", nextUID))
			case "departmentNumber":
				for _ , value := range conn.Config.GroupValues.Groups {
					departments = departments + " " + value
				}
				utils.PrintColor(consts.Yellow, fmt.Sprintf("\t\t** Default to: %s\n",
					conn.Config.DefaultValues.GroupName))
				utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tValid departments: %s\n", departments))
			case "loginShell":
				for _ , value := range conn.Config.DefaultValues.ValidShells {
					shells = shells + " " + value
				}
				utils.PrintColor(consts.Yellow, fmt.Sprintf("\t\t**Default to: %s\n",
					conn.Config.DefaultValues.Shell))
				utils.PrintColor(consts.Purple, fmt.Sprintf("\t\t**valid shells: %s\n", shells))
			case "userPassword":
				passWord = utils.GenerateRandom(
					conn.Config.DefaultValues.PassComplex,
					conn.Config.DefaultValues.PassLenght)
				utils.PrintColor(consts.Yellow,
					fmt.Sprintf("\t\tPress Enter to accept the given password\n\t\tSuggested password: %s\n", passWord))

		}
		if vars.RecordFields[idx].Default != "" {
			utils.PrintColor(consts.Cyan,
				fmt.Sprintf("\tDefault to: %s\n", vars.RecordFields[idx].Default))
		}
		fmt.Printf("\t%s: ", vars.RecordFields[idx].Prompt)
		reader := bufio.NewReader(os.Stdin)
		enterData, _ := reader.ReadString('\n')
		enterData = strings.ToLower(strings.TrimSuffix(enterData, "\n"))
		switch fieldName {
			case "uid":
				cnt := conn.GetUser(enterData)
				if cnt != 0 {
					utils.PrintColor(consts.Red,
						fmt.Sprintf("\tGiven user %s already exist, aborting...\n\n", enterData))
					return false
				}
				utils.PrintColor(consts.Purple, fmt.Sprintf("\tUsing user: %s\n", enterData))
			case "givenName", "sn": enterData = strings.Title(enterData)
			case "mail":
				if len(enterData) == 0 {
					enterData = email
				}
			case "uidNumber":
				if len(enterData) == 0 {
					enterData = strconv.Itoa(nextUID)
				}
			case "departmentNumber" :
				if len(enterData) == 0 {
					enterData = conn.Config.DefaultValues.GroupName
				}
			case "loginShell" :
				if len(enterData) == 0 {
					enterData = conn.Config.DefaultValues.Shell
				}
			case "userPassword" :
				if len(enterData) == 0 {
					enterData = passWord
				}
		}
		if len(enterData) == 0 && vars.RecordFields[idx].NoEmpty == true {
				utils.PrintColor(consts.Red, "\tNo value was entered aborting...\n\n")
				return false
		}

		conn.User.Strings[fieldName] = vars.StringRecord{Value: enterData}
		// fmt.Printf("<uid %d fieldName %s>\n", nextUID, fieldName)
		fmt.Printf("<%v>\n", conn.User.Strings[fieldName].Value)
	}
	return true
}

func Create(conn *ldap.Connection)  {
	utils.PrintHeader(consts.Purple, "create user")
	createRecord(conn)
}
