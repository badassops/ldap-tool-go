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
  "strings"
  "strconv"

  v "badassops.ldap/vars"
  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
)

var (
)

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
  fmt.Printf(" %v \n", v.ModRecord)
  return true
}

func Create(c *l.Connection) {
  u.PrintHeader(u.Purple, "Create Sudo Rule", true)
  createSudoRecord(c)
  u.PrintLine(u.Purple)
}
