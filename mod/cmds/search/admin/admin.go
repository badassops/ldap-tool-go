// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package admin

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"badassops.ldap/constants"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/search/common"
)

func Search(conn *ldap.Connection, mode string) {
	utils.PrintHeader(constants.Purple, "Search " +  mode)
	reader := bufio.NewReader(os.Stdin)
	givenValue := ""
	wildCard := false
	cnt := 0
	if strings.HasSuffix(mode, "s") == false {
		givenValue, wildCard = common.EnterValue(mode)
		if givenValue == "" {
			return
		}
	}
	switch mode {
		case "admin":
			utils.PrintLine(utils.Purple)
			if !wildCard {
				cnt = conn.SearchAdminGroup(givenValue, false, true, false)
			} else {
				cnt = conn.SearchAdminGroup(givenValue, false, false, true)
				if cnt == 0 {
					utils.PrintColor(constants.Red, fmt.Sprintf("\tNo group match %s\n", givenValue))
					return
				}
				fmt.Printf("\n\tSelect the group from the above list: ")
				givenValue, _ = reader.ReadString('\n')
				givenValue = strings.TrimSuffix(givenValue, "\n")
				if givenValue == "" {
					utils.PrintColor(utils.Red, fmt.Sprintf("\tNo group was given aborting...\n"))
					return
				}
				utils.PrintLine(utils.Purple)
				cnt = conn.SearchAdminGroup(givenValue, false, true, false)
			}

		case "admins":
			cnt = conn.SearchAdminGroup("", true, true, false)
	}
	if cnt == 0 {
		utils.PrintColor(constants.Red, fmt.Sprintf("\tNo group match %s\n", givenValue))
		return
	}
	utils.PrintLine(utils.Purple)
}
