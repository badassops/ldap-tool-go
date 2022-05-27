// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package modify

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  v "badassops.ldap/vars"
  cs "badassops.ldap/cmds/common/sudo"
)

var (
  valueEntered string

  allowedField = []string{"sudoCommand", "sudoHost", "sudoOption",
    "sudoOrder", "sudoRunAsUser" }
)

func deleteFields(c *l.Connection) int {
  recordChanged := 0
  reader := bufio.NewReader(os.Stdin)
  records, _ := c.GetSudoCN(v.ModRecord.Field["cn"])
  for _, entry := range records {
    u.PrintBlue(fmt.Sprintf("\tDN: %s\n", entry.DN))
    for _, attributes := range entry.Attributes {
      for _, value := range attributes.Values {
        if (attributes.Name != "objectClass") && (attributes.Name != "cn") {
          u.PrintCyan(fmt.Sprintf("\tField: %s\n", attributes.Name))
          switch attributes.Name {
            case "sudoCommand":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))

            case "sudoHost":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))

            case "sudoOption":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))

            case "sudoOrder":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))

            case "sudoRunAsUser":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))
          }
          fmt.Printf("\tEnter %sdelete%s to delete or press enter to keep: ", u.RedUnderline, u.Off)
          valueEntered, _ = reader.ReadString('\n')
          valueEntered = strings.TrimSuffix(valueEntered, "\n")
          if valueEntered == "delete" {
            v.ModSudo.DelList[attributes.Name] = append(v.ModSudo.DelList[attributes.Name], value)
            recordChanged++
          }
        }
      }
    }
  }
  return recordChanged
}

func addFields(c *l.Connection) int {
  newRecord := 0
  u.PrintCyan(fmt.Sprintf("\n\tEach field can have multiple entries\n"))
  u.PrintCyan(fmt.Sprintf("\tPress enter to skip, or enter value for field\n"))
  for _, fieldname := range allowedField {
    for true {
      u.PrintGreen(fmt.Sprintf("\tField %s: enter value: ", fieldname))
      reader := bufio.NewReader(os.Stdin)
      valueEntered, _ := reader.ReadString('\n')
      valueEntered = strings.TrimSuffix(valueEntered, "\n")
      if len(valueEntered) != 0 {
        v.ModSudo.AddList[fieldname] = append(v.ModSudo.AddList[fieldname], valueEntered)
        newRecord++
      } else {
        fmt.Printf("\n")
        break
      }
    }
  }
  return newRecord
}

func createSudoModRecord(c *l.Connection) int {
  changes := 0
  fmt.Printf("\n\t%sDeleting%s fields %sfirst%s...\n", u.Red, u.Off, u.Cyan, u.Off)
  if deleteFields(c) > 0 {
    changes++
  } else {
    u.PrintBlue("\tNo record will be deleted \n")
  }

  fmt.Printf("\n\t%sAdding%s fields %ssecond%s...\n", u.Red, u.Off, u.Cyan, u.Off)
  if addFields(c) > 0 {
    changes++ 
  } else {
    u.PrintBlue("\tNo record will be added \n")
  }
  return changes
}

func Modify(c *l.Connection) {
  u.PrintHeader(u.Purple, "Modify Sudo rule", true)
  if cs.Sudo(c, true, false) {
    if u.InList(c.Config.SudoValues.ExcludeSudo, v.ModRecord.Field["cn"]) {
      u.PrintRed(fmt.Sprintf("\t\tThe sudo rule %s can not be modified\n", v.ModRecord.Field["cn"]))
    } else {
      if createSudoModRecord(c) > 0 {
        v.ModSudo.DN = fmt.Sprintf("cn=%s,ou=%s", v.ModRecord.Field["cn"], c.Config.SudoValues.SudoersBase)
        fmt.Printf("\tDN  %v\n", v.ModSudo.DN)
        fmt.Printf("\tADD %v\n", v.ModSudo.AddList)
        fmt.Printf("\tDEL %v\n", v.ModSudo.DelList)
      }
    }
  }
  u.PrintLine(u.Purple)
}
