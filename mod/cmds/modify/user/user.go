// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package modify

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
	"badassops.ldap/cmds/common/user"
)

var (
	fields = []string{"givenName", "sn", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowMax", "shadowExpire", "sshPublicKey"}
)

func cleanUpData (conn *ldap.Connection, data string) string {
	// remove the userDN part
	newData := strings.Split(data, ",")[0]
	// remove the cn= part
	newData = strings.TrimPrefix(newData, "cn=")
	// remove the groupDN part
	groupDN := fmt.Sprintf(",%s", conn.Config.ServerValues.GroupDN)
	return strings.TrimPrefix(newData, groupDN )
}

func createNewUserRecord(conn *ldap.Connection) {
	var shadowMaxChanged bool = false
	var shells string
	var departments string
	var currShadowMax string
	var changeFields map[string]string
	var userGroupList []string
	var delList []string
    var addList []string
	changeFields = make(map[string]string)
	utils.PrintColor(consts.Purple,
		fmt.Sprintf("\tUsing user: %s\n", conn.User.Field["uid"]))
	utils.PrintColor(consts.Yellow,
		fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	utils.PrintLine(utils.Purple)
	for _, fieldName := range fields  {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
			case "givenName":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value: %s\n", conn.User.Field[fieldName]))

			case "sn":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value: %s\n", conn.User.Field[fieldName]))

			case "mail":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value: %s\n", conn.User.Field[fieldName]))

			case "departmentNumber":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value: %s\n", conn.User.Field[fieldName]))
				for _ , value := range conn.Config.GroupValues.Groups {
					departments = departments + " " + value
				}
				utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tValid departments:%s\n", departments))

			case "loginShell":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value: %s\n", conn.User.Field[fieldName]))
				for _ , value := range conn.Config.DefaultValues.ValidShells {
					shells = shells + " " + value
				}
				utils.PrintColor(consts.Purple, fmt.Sprintf("\t\tValid shells:%s\n", shells))

			case "userPassword":
				passWord := utils.GenerateRandom(
					conn.Config.DefaultValues.PassComplex,
					conn.Config.DefaultValues.PassLenght)
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value (encrypted!): %s\n", conn.User.Field[fieldName]))
				utils.PrintColor(consts.Yellow,
					fmt.Sprintf("\t\tsuggested password: %s\n", passWord))

			case "shadowMax":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent max password age: %s\n", conn.User.Field[fieldName]))
				utils.PrintColor(consts.Purple,
					fmt.Sprintf("\t\tMin %d days and max %d days\n",
					conn.Config.DefaultValues.ShadowMin,
					conn.Config.DefaultValues.ShadowMax))

			case "shadowExpire":
				value, _ := strconv.ParseInt(conn.User.Field["shadowExpire"], 10, 64)
				_, passExpired := utils.GetReadableEpoch(value * 86400)
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent password will expire on: %s\n", passExpired))

			case "sshPublicKey":
				utils.PrintColor(consts.Cyan,
					fmt.Sprintf("\tCurrent value: %s\n", conn.User.Field[fieldName]))

		}
		prefix := ""
		if conn.Config.Debug {
			prefix = fmt.Sprintf("(%s) - ", fieldName)
		}
		if fieldName == "shadowExpire" {
			fmt.Printf("\t%sReset password expired to (%s days from now) Y/N: ",
				prefix, changeFields["shadowMax"])
		} else {
			fmt.Printf("\t%s%s: ", prefix, vars.Template[fieldName].Prompt)
		}
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch fieldName {
			case "givenName", "sn":
				valueEntered = strings.Title(valueEntered)

			case "mail":
				valueEntered = strings.ToLower(valueEntered)

			case "shadowMax":
				if len(valueEntered) != 0 {
					shadowMax, _ := strconv.Atoi(valueEntered)
					if shadowMax < conn.Config.DefaultValues.ShadowMin ||
						shadowMax > conn.Config.DefaultValues.ShadowMax {
						utils.PrintColor(consts.Red,
						fmt.Sprintf("\t\tGiven value %d, is out or range, is set to %d\n",
							shadowMax, conn.Config.DefaultValues.ShadowAge))
						valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowAge)
					}
					shadowMaxChanged = true
				}

			case "shadowExpire":
				if len(valueEntered) == 0 {
					utils.PrintColor(consts.Cyan,
						fmt.Sprintf("\tPassword expiration date will not be changed\n"))
				} else {
					// calculate when it will be expired based on default value if shadowMax
					// otherwise it will be today + new shadowMax value
					currShadowMax = conn.User.Field["shadowMax"]
					if shadowMaxChanged == true {
						currShadowMax = changeFields["shadowMax"]
					}

					// set last changed to now
					changeFields["shadowLastChange"] = vars.Template["shadowLastChange"].Value
					// calculate the new shadowExpire
					shadowLastChange, _ := strconv.ParseInt(changeFields["shadowLastChange"], 10, 64)
					shadowMax, _ := strconv.ParseInt(currShadowMax, 10, 64)
					_, passExpired := utils.GetReadableEpoch((shadowLastChange + shadowMax) * 86400)
					utils.PrintColor(consts.Cyan,
						fmt.Sprintf("\tCurrent password will now expire on: %s\n", passExpired))
					// replace the 'Y' with the correct value
					valueEntered = strconv.FormatInt((shadowLastChange + shadowMax), 10)
				}
		}
		if len(valueEntered) != 0 {
				changeFields[fieldName] = valueEntered
		}
	}
	conn.ModifyUser(changeFields)
	// we only handle groupOfNames type of group
	fmt.Printf("\n")
	for _, group := range conn.User.Groups {
		userGroupList = append(userGroupList, cleanUpData(conn, group))
	}
	availableGroups := conn.GetGroupsName("groupOfNames")
	utils.PrintColor(utils.Purple, fmt.Sprintf("\t\tAvailable groups: %s\n",
		strings.Join(availableGroups[:], " ")))

	utils.PrintColor(utils.Purple, fmt.Sprintf("\t\tUser %s groups: %s\n",
		conn.User.Field["uid"], strings.Join(userGroupList, " ")))
	for _, leaveGroup := range userGroupList {
		utils.PrintColor(consts.Red,
			fmt.Sprintf("\tRemove the group %s? default to not remove group, [Y/N]:  ", leaveGroup))
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch valueEntered {
			case "y", "Y": delList = append(delList, leaveGroup)
		}
	}
	for _, joinGroup := range availableGroups {
		if utils.InList(userGroupList, joinGroup) == false {
			utils.PrintColor(consts.Green,
				fmt.Sprintf("\tJoin the group %s? default to not join group, [Y/N]:  ", joinGroup))
			reader := bufio.NewReader(os.Stdin)
			valueEntered, _ := reader.ReadString('\n')
			valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
			switch valueEntered {
				case "y", "Y": addList = append(addList, joinGroup)
			}
		}
	}
	conn.ModifyUserGroup(conn.User.Field["uid"], addList, delList)
}

func Modify(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Modify User", true)
	if common.User(conn, true, false) {
		createNewUserRecord(conn)
	}
	utils.PrintLine(utils.Purple)
	return
}
