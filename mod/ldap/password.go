// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package ldap

import (
  "fmt"
  "os"

  l "badassops.ldap/logs"
  u "badassops.ldap/utils"

  ldapv3 "gopkg.in/ldap.v2"
)

func (c *Connection) modifyPassword(recordId, recordType string, request *ldapv3.PasswordModifyRequest) bool {
  if _, err := c.Conn.PasswordModify(request); err != nil {
    msg = fmt.Sprintf("Error modifying the %s for user %s, error %s", recordType, recordId, err.Error())
    l.Log(msg, "ERROR")
    return false
  }
  msg = fmt.Sprintf("The %s of the user %s has been modified", recordType, recordId)
  l.Log(msg, "INFO")
  return true
}

func (c *Connection) setPassword() bool {
  // save before it get encrypted
  userPassData := make(map[string]string)
  userPassData["user"] = c.User.Field["uid"]
  userPassData["password"] = c.User.Field["userPassword"]

  // once the record is create we need to hash the password
  passwordReq := ldapv3.NewPasswordModifyRequest(
    c.User.Field["dn"],
    c.User.Field["userPassword"],
    c.User.Field["userPassword"])

  if !c.modifyPassword(c.User.Field["uid"], "password", passwordReq) {
    return false
  }
  passDir := c.Config.LogValues.LogsDir + "/users/"
  passFile := passDir + c.User.Field["uid"] + ".info"
  os.MkdirAll(passDir, 0750)
  ok, errorMsg := u.RecordPassword(passFile, userPassData)
  if !ok {
    msg = fmt.Sprintf("Unable to save the user %s's password in the file %s, error %s",
      c.User.Field["uid"], passFile, errorMsg)
    l.Log(msg, "ERROR")
  }
  return true
}
