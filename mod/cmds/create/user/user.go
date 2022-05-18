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
	"badassops.ldap/logs"
)

var (
// not required for create a new user : cn, gidNumber, displayName, gecos
// homeDirectory, shadowLastChange, shadowLastChange
// groups is handled seperat;y

	fields = []string{"uid", "givenName", "sn",
		"uidNumber", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowWarning", "shadowMax",
		"sshPublicKey"}

	// construct base on FirstName + LastName
	 userFullname = []string{"cn", "displayName", "gecos"}
)

func createUserRecord(conn *ldap.Connection) bool {
	var email string
	var nextUID int
	var passWord string
	var shells string
	var departments string
	var shadowMax int
	var logRecord string

	for _, fieldName := range fields {

		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
			case "uid":
				utils.PrintColor(consts.Yellow,
					fmt.Sprintf("\tThe userid / login name is case sensitive, it will be made all lowercase\n"))
			case "mail":
				email = fmt.Sprintf("%s.%s@%s",
					strings.ToLower(conn.User.Field["givenName"]),
					strings.ToLower(conn.User.Field["sn"]),
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
				utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tValid departments:%s\n", departments))
			case "loginShell":
				for _ , value := range conn.Config.DefaultValues.ValidShells {
					shells = shells + " " + value
				}
				utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tValid shells:%s\n", shells))
			case "userPassword":
				passWord = utils.GenerateRandom(
					conn.Config.DefaultValues.PassComplex,
					conn.Config.DefaultValues.PassLenght)
				utils.PrintColor(consts.Purple, "\t\tPress Enter to accept the suggested password\n")
				utils.PrintColor(consts.Yellow, fmt.Sprintf("\tSuggested password: %s\n", passWord))
			case "shadowMax":
				utils.PrintColor(consts.Purple,
					fmt.Sprintf("\t\tMin %d days and max %d days\n",
						conn.Config.DefaultValues.ShadowMin,
						conn.Config.DefaultValues.ShadowMax))
		}
		if vars.Template[fieldName].Value != "" {
			utils.PrintColor(consts.Yellow,
				fmt.Sprintf("\t ** Default to: %s **\n", vars.Template[fieldName].Value))
		}
		if conn.Config.Debug {
			fmt.Printf("\t(%s) - %s: ", fieldName, vars.Template[fieldName].Prompt)
		} else {
			fmt.Printf("\t%s: ", vars.Template[fieldName].Prompt)
		}
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch fieldName {
			case "uid":
				cnt := conn.CheckUser(valueEntered)
				if cnt != 0 {
					utils.PrintColor(consts.Red,
						fmt.Sprintf("\n\tGiven user %s already exist, aborting...\n\n", valueEntered))
					return false
				}
				utils.PrintColor(consts.Purple, fmt.Sprintf("\tUsing user: %s\n", valueEntered))
			case "givenName", "sn": valueEntered = strings.Title(valueEntered)
			case "mail":
				if len(valueEntered) == 0 {
					valueEntered = email
				}
			case "uidNumber":
				if len(valueEntered) == 0 {
					valueEntered = strconv.Itoa(nextUID)
				}
			case "departmentNumber" :
				if len(valueEntered) == 0 {
					valueEntered = conn.Config.DefaultValues.GroupName
					conn.User.Field["gidNumber"] = strconv.Itoa(conn.Config.DefaultValues.GroupId)
				} else {
					for _, mapValues := range conn.Config.GroupValues.GroupsMap {
						if mapValues.Name == valueEntered {
							conn.User.Field["gidNumber"] = strconv.Itoa(mapValues.Gid)
							}
					}
				}
			case "loginShell" :
				if len(valueEntered) == 0 {
					valueEntered = "/bin/" + conn.Config.DefaultValues.Shell
				} else {
					valueEntered = "/bin/" + valueEntered
				}
			case "userPassword" :
				if len(valueEntered) == 0 {
					valueEntered = passWord
				}
			case "shadowMax":
				if len(valueEntered) == 0 {
					valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowMax)
				} else {
					shadowMax, _ = strconv.Atoi(valueEntered)
					if shadowMax < conn.Config.DefaultValues.ShadowMin ||
						shadowMax > conn.Config.DefaultValues.ShadowMax {
						utils.PrintColor(consts.Yellow,
							fmt.Sprintf("\tGiven value %d, is out or range, is set to %d\n",
								shadowMax, conn.Config.DefaultValues.ShadowAge))
						valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowAge)
					}
				}
			case "shadowWarning":
				if len(valueEntered) == 0 {
					valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowWarning)
				}
			default:
				if len(valueEntered) == 0 {
					valueEntered = vars.Template[fieldName].Value
				}
		}
		if len(valueEntered) == 0 && vars.Template[fieldName].NoEmpty == true {
				utils.PrintColor(consts.Red, "\tNo value was entered aborting...\n\n")
				return false
		}

		// update the user record so it can be submitted
		conn.User.Field[fieldName] = valueEntered
	}
	// setup the groups for the user
	utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tSpecial Groups: %v\n", conn.Config.GroupValues.SpecialGroups))
	utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tEnter 'add' or press enter to skip\n"))
	for _, userGroup := range conn.Config.GroupValues.SpecialGroups {
		utils.PrintColor(consts.Yellow, fmt.Sprintf("\tGroup %s (add)? : ", userGroup))
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		if valueEntered == "add" {
			conn.User.Groups = append(conn.User.Groups , userGroup)
		}
	}
	// this are always firstName lastName
	for _, userFullnameFields := range userFullname {
		conn.User.Field[userFullnameFields] = conn.User.Field["givenName"] + " " + conn.User.Field["sn"]
	}
	// dn is create base on given uid and user DN
	conn.User.Field["dn"] = fmt.Sprintf("uid=%s,%s", conn.User.Field["uid"], conn.Config.ServerValues.UserDN)
	// this is always /home + userlogin
	conn.User.Field["homeDirectory"] = "/home/" + conn.User.Field["uid"]
	// initialized to be today's epoch days
	conn.User.Field["shadowExpire"] = vars.Template["shadowExpire"].Value
	conn.User.Field["shadowLastChange"] = vars.Template["shadowLastChange"].Value
	// debug
	if conn.Config.Debug {
		for recordName, recordValue := range conn.User.Field {
			logRecord = fmt.Sprintf(" Field Name: %s - Field Value: %s", recordName, recordValue)
			logs.Log(logRecord, "DEBUG")
		}
		logs.Log(fmt.Sprintf("User's Groups: %v", conn.User.Groups), "DEBUG")
	}
	return true
}

func Create(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Create User", true)
	if createUserRecord(conn) {
		utils.PrintLine(utils.Purple)
		if !conn.AddRecord() {
			utils.PrintColor(consts.Red,
				fmt.Sprintf("\n\tFailed adding the user %s, check the log file\n", conn.User.Field["uid"]))
		} else{
			utils.PrintColor(consts.Green, fmt.Sprintf("\n\tUser %s added successfully\n", conn.User.Field["uid"]))
		}
	}
	utils.PrintLine(utils.Purple)
	return
}
