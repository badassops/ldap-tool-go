// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"badassops.ldap/utils"
)

func EnterValue(dataID string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tEnter %s to be use: ", dataID)
	enterData, _ := reader.ReadString('\n')
	enterData = strings.TrimSuffix(enterData, "\n")

	if enterData == "" {
		utils.PrintColor(utils.Red, fmt.Sprintf("\n\tNo %s was given aborting...\n", dataID))
		return "", false
	}
	fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
	wildCard, _ := reader.ReadString('\n')
	wildCard = strings.TrimSuffix(wildCard, "\n")
	if utils.GetYN(wildCard, false) == false {
		return enterData, false
	}
	return enterData, true
}
