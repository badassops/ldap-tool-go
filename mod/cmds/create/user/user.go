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

  "badassops.ldap/vars"
  u "badassops.ldap/utils"
  "badassops.ldap/ldap"
  //"badassops.ldap/logs"
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

  // given field value
  email       string
  passWord    string
  shells      string
  departments string
  nextUID     int
  shadowMax   int
)

func createUserRecord(c *ldap.Connection) bool {
  //var logRecord string

  for _, fieldName := range fields {
    // these will be valid once the field was filled since they depends
    // on some of the fields value
    switch fieldName {
      case "uid":
        u.PrintColor(u.Yellow,
          fmt.Sprintf("\tThe userid / login name is case sensitive, it will be made all lowercase\n"))

      case "mail":
        email = fmt.Sprintf("%s.%s@%s",
          strings.ToLower(c.User.Field["givenName"]),
          strings.ToLower(c.User.Field["sn"]),
          c.Config.ServerValues.EmailDomain)
        u.PrintCyan(fmt.Sprintf("\tDefault email: %s\n", email))

      case "uidNumber":
        nextUID = c.GetNextUID()
        u.PrintPurple(fmt.Sprintf("\t\tOptional set user UID, press enter to use the next UID: %d\n", nextUID))

      case "departmentNumber":
        for _ , value := range c.Config.GroupValues.Groups {
          departments = departments + " " + value
        }
        u.PrintPurple(fmt.Sprintf("\t\tValid departments:%s\n", departments))

      case "loginShell":
        for _ , value := range c.Config.DefaultValues.ValidShells {
          shells = shells + " " + value
        }
        u.PrintPurple(fmt.Sprintf("\t\tValid shells:%s\n", shells))

      case "userPassword":
        passWord = u.GenerateRandom(
          c.Config.DefaultValues.PassComplex,
          c.Config.DefaultValues.PassLenght)
        u.PrintPurple("\t\tPress Enter to accept the suggested password\n")
        u.PrintYellow(fmt.Sprintf("\tSuggested password: %s\n", passWord))

      case "shadowMax":
        u.PrintPurple(fmt.Sprintf("\t\tMin %d days and max %d days\n",
          c.Config.DefaultValues.ShadowMin,
          c.Config.DefaultValues.ShadowMax))
    }

    if vars.Template[fieldName].Value != "" {
      u.PrintYellow(fmt.Sprintf("\t ** Default to: %s **\n", vars.Template[fieldName].Value))
    }

    if c.Config.Debug {
      fmt.Printf("\t(%s) - %s: ", fieldName, vars.Template[fieldName].Prompt)
    } else {
      fmt.Printf("\t%s: ", vars.Template[fieldName].Prompt)
    }

    reader := bufio.NewReader(os.Stdin)
    valueEntered, _ := reader.ReadString('\n')
    valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))

    switch fieldName {
      case "uid":
        if cnt := c.CheckUser(valueEntered); cnt != 0  {
          u.PrintRed(fmt.Sprintf("\n\tGiven user %s already exist, aborting...\n\n", valueEntered))
          return false
        }
        u.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", valueEntered))

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
          valueEntered = strings.ToUpper(c.Config.DefaultValues.GroupName)
          c.User.Field["gidNumber"] = strconv.Itoa(c.Config.DefaultValues.GroupId)
        } else {
          if cnt := c.CheckGroup(valueEntered); cnt == 0 {
            u.PrintRed(fmt.Sprintf("\n\tGiven departmentNumber %s is not valid, aborting...\n\n",
              valueEntered))
            return false
          }
          for _, mapValues := range c.Config.GroupValues.GroupsMap {
            if mapValues.Name == valueEntered {
              c.User.Field["gidNumber"] = strconv.Itoa(mapValues.Gid)
            }
          }
          valueEntered = strings.ToUpper(valueEntered)
        }

      case "loginShell" :
        if len(valueEntered) == 0 {
          valueEntered = "/bin/" + c.Config.DefaultValues.Shell
        } else {
          if u.InList(c.Config.DefaultValues.ValidShells, valueEntered) {
            u.PrintRed(fmt.Sprintf("\n\tGiven shell %s is not valid, aborting...\n\n", valueEntered))
            return false
          }
          valueEntered = "/bin/" + valueEntered
        }

      case "userPassword" :
        if len(valueEntered) == 0 {
          valueEntered = passWord
        }

      case "shadowMax":
        if len(valueEntered) == 0 {
          valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowMax)
        } else {
          shadowMax, _ = strconv.Atoi(valueEntered)
          if shadowMax < c.Config.DefaultValues.ShadowMin ||
            shadowMax > c.Config.DefaultValues.ShadowMax {
            u.PrintYellow(fmt.Sprintf("\tGiven value %d, is out or range, is set to %d\n",
                shadowMax, c.Config.DefaultValues.ShadowAge))
            valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowAge)
          }
        }

      case "shadowWarning":
        if len(valueEntered) == 0 {
          valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowWarning)
        }

      default:
        if len(valueEntered) == 0 {
          valueEntered = vars.Template[fieldName].Value
        }
    }

    if len(valueEntered) == 0 && vars.Template[fieldName].NoEmpty == true {
        u.PrintRed("\tNo value was entered aborting...\n\n")
        return false
    }

    // update the user record so it can be submitted
    c.User.Field[fieldName] = valueEntered
  }

  // setup the groups for the user
  u.PrintPurple(fmt.Sprintf("\t\tSpecial Groups: %v\n", c.Config.GroupValues.SpecialGroups))
  u.PrintPurple(fmt.Sprintf("\t\tEnter 'add' or press enter to skip\n"))
  for _, userGroup := range c.Config.GroupValues.SpecialGroups {
    u.PrintYellow(fmt.Sprintf("\tGroup %s (add)? : ", userGroup))
    reader := bufio.NewReader(os.Stdin)
    valueEntered, _ := reader.ReadString('\n')
    valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
    if valueEntered == "add" {
      c.User.Groups = append(c.User.Groups , userGroup)
    }
  }

  // this are always firstName lastName
  for _, userFullnameFields := range userFullname {
    c.User.Field[userFullnameFields] = c.User.Field["givenName"] + " " + c.User.Field["sn"]
  }

  // dn is create base on given uid and user DN
  c.User.Field["dn"] = fmt.Sprintf("uid=%s,%s", c.User.Field["uid"], c.Config.ServerValues.UserDN)

  // this is always /home + userlogin
  c.User.Field["homeDirectory"] = "/home/" + c.User.Field["uid"]

  // initialized to be today's epoch days
  c.User.Field["shadowExpire"] = vars.Template["shadowExpire"].Value
  c.User.Field["shadowLastChange"] = vars.Template["shadowLastChange"].Value
  return true
}

func Create(c *ldap.Connection) {
  u.PrintHeader(u.Purple, "Create User", true)
  if createUserRecord(c) {
    u.PrintLine(u.Purple)
    if !c.AddUser() {
      u.PrintRed(fmt.Sprintf("\n\tFailed adding the user %s, check the log file\n", c.User.Field["uid"]))
    } else{
      u.PrintGreen(fmt.Sprintf("\n\tUser %s added successfully\n", c.User.Field["uid"]))
    }
  }
  u.PrintLine(u.Purple)
}
