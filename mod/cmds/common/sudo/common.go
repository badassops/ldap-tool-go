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
  //"strconv"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  v "badassops.ldap/vars"
  ldapv3 "gopkg.in/ldap.v2"
)

var (
  sudoRecord []*ldapv3.Entry
  recordCount int
)

func printSudoRecord(record []*ldapv3.Entry) {
  for _, entry := range record {
    u.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
    for _, attributes := range entry.Attributes {
      for _, value := range attributes.Values {
        if attributes.Name != "objectClass" {
          if attributes.Name == "cn" {
            u.PrintBlue(fmt.Sprintf("\t%s : %s \n", attributes.Name, value ))
          } else {
            u.PrintCyan(fmt.Sprintf("\t%s : %s \n", attributes.Name, value ))
          }
        }
      }
    }
  }
}

func Sudo(c *l.Connection, firstTime bool, showRecord bool) bool {
  reader := bufio.NewReader(os.Stdin)
  fmt.Printf("\tEnter the sudo CN to be use: ")
  enterData, _ := reader.ReadString('\n')
  enterData = strings.TrimSuffix(enterData, "\n")

  if enterData == "" {
    u.PrintRed(fmt.Sprintf("\n\tNo users was given aborting...\n"))
    return false
  }

  if firstTime {
    fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
    wildCard, _ := reader.ReadString('\n')
    wildCard = strings.TrimSuffix(wildCard, "\n")
    if u.GetYN(wildCard, false) == true {
      enterData = "*" + enterData + "*"
      c.SearchSudoCN(enterData, true)
      fmt.Printf("\n\tSelect the cn from the above list:\n")
      Sudo(c, false, showRecord)
      return true
    }
  }

  if sudoRecord, recordCount = c.GetSudoCN(enterData); recordCount == 0 {
    u.PrintColor(u.Red, fmt.Sprintf("\n\tCN %s was not found, aborting...\n", enterData))
    return false
  }
  if showRecord {
    printSudoRecord(sudoRecord)
  }
  v.ModRecord.Field["cn"] = enterData
  return true
}
