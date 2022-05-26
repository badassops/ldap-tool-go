// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package delete

import (
  // "fmt"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
)

func Delete(c *l.Connection) {
  u.PrintHeader(u.Purple, "Delete Sudo rule", true)
  u.PrintLine(u.Purple)
}
