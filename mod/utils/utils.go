// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package utils

import (
  "bufio"
  "fmt"
  "math/rand"
  "os"
  "os/user"
  "os/signal"
  "path"
  "time"
  "runtime"
  "strings"
  "strconv"
  "syscall"

  l "badassops.ldap/logs"
  ps "github.com/mitchellh/go-ps"
)


const (
  Off      = "\x1b[0m"       // Text Reset
  Black    = "\x1b[1;30m"    // Black
  Red      = "\x1b[1;31m"    // Red
  Green    = "\x1b[1;32m"    // Green
  Yellow   = "\x1b[1;33m"    // Yellow
  Blue     = "\x1b[1;34m"    // Blue
  Purple   = "\x1b[1;35m"    // Purple
  Cyan     = "\x1b[1;36m"    // Cyan
  White    = "\x1b[1;37m"    // White

  RedBase    = "\x1b[0;31m"    // Red no highlighted
  Greenbase  = "\x1b[0;32m"    // Green no highlighted
  YellowBase = "\x1b[0;33m"    // Yellow no highlighted
  BlueBase   = "\x1b[0;34m"    // Blue no highlighted
  PurpleBase = "\x1b[0;35m"    // Purple no highlighted
  CyanBase   = "\x1b[0;36m"    // Cyan no highlighted
  WhiteBase  = "\x1b[0;37m"    // White no highlighted

  RedUnderline = "\x1b[4;31m"    // Red underline 
  OneLineUP    = "\x1b[A"

  clearLine    = "\x1b[0G\x1b[2K\x1b[0m\r"
  clearScreen  = "\x1b[H\x1b[2J"
  HEADER       = "---------------"
  LINE         = "_________________________________________________"
  alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
  specialChars = "@#$%^*(){}[]<>/\\"

  USformat  = "Jan 2, 2006 15:04:05"

)

var (
  prefixes = []string{"main."}

  removeLines = strings.NewReplacer(
    "\r\n", "\\r\\n",
    "\r", "\\r",
    "\n", "\\n")

  DangerZone = fmt.Sprintf("%sDanger Zone%s, be sure you understand the implication!",
    RedUnderline, Off)
)

type Wheel struct {
  StopChannel  chan int
  Speed        time.Duration
}

// function to exit if an error occured
func ExitIfError(err error) {
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: " + err.Error())
    l.Log(fmt.Sprint(err.Error()), "ERROR")
    os.Exit(1)
  }
}

// function to exit if pointer is nill
func ExitIfNill(ptr interface{}) {
  if ptr == nil {
    fmt.Fprintln(os.Stderr, "Error: got a nil pointer.")
    l.Log("got a nil pointer", "ERROR")
    os.Exit(1)
  }
}

// function to exit wihthe given error message
func ExitWithMesssage(messsage string) {
  PrintColor(Red, messsage)
  os.Exit(1)
}

// function to check if the process is run as root
func IsRoot() (bool) {
  if os.Geteuid() == 0 {
    return true
  }
  return false
}

// function to check if the process is run as the given user
func IsUser(runUser string) (bool) {
  user, err := user.Current()
  ExitIfError(err)
  if user.Username == runUser {
    return true
  }
  return false
}

// function to get the user name of the runming process
func RunningUser() (string) {
  user, err := user.Current()
  ExitIfError(err)
  return user.Username
}

// function check if file or directory exist
// isFile : true check for file
// isFile : false check for directory
func Exist(fullPath string, isFile bool, verbose bool) (bool, error) {
  fileInfo, errStat := os.Stat(fullPath)
  if errStat != nil {
    if !verbose {
      // check if the parent directoy exist
      dirPath, _ := path.Split(fullPath)
      if _, err := os.Stat(dirPath); err != nil {
        fmt.Fprintln(os.Stderr, "Parent directory " + dirPath + " does not exist\n")
      }
    }
    if verbose {
      fmt.Fprintln(os.Stderr, "Error: " + fmt.Sprint(errStat))
    }
    return false, errStat
  }

  object := fileInfo.IsDir()
  // check is not a directory if isFile is true
  if isFile && object {
    if verbose {
      fmt.Fprintln(os.Stderr, "Given " + fullPath + " is a directory")
    }
    return false, nil
  }
  // check is not a file if isFile is false
  if !isFile && !object {
    if verbose {
      fmt.Fprintln(os.Stderr, "Given " + fullPath + " is a file")
    }
    return false, nil
  }
  return true, nil
}

// function to capture reveived signal
func SignalHandler(lockFile string, pid int) {
  interrupt := make(chan os.Signal, 1)
  // we handle only these signal: SIGINT(2) - SIGTRAP(5) - SIGKILL(9) - SIGTERM(15), SIGHUP(1)
  signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTRAP, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP)
  go func() {
    sigId := <-interrupt
    fmt.Fprintln(os.Stderr, "Received signal " + fmt.Sprintf("%v %d", sigId, sigId))
    exit, _ := strconv.Atoi(fmt.Sprintf("%d", sigId))
    LockFile(lockFile, pid, 2)
    os.Exit(exit)
  }()
}

// function to create or remove a lock file
// mode 1 : create
// mode 2 : remove
// mode 3 : get pid value
func LockFile(lockFile string, pid int, mode int) (bool, error) {
  state, err := Exist(lockFile, true, false)
  switch mode {
    case 1:
      // file should not exist
      if err == nil && state {
        return false, err
      }
      fp, err := os.Create(lockFile)
      ExitIfError(err)
      defer fp.Close()
      _, err = fp.WriteString(strconv.Itoa(pid))
      ExitIfError(err)

    case 2:
      // file should exist
      if !state {
        return false, err
      }
      ExitIfError(os.Remove(lockFile))

  }
  return true, nil
}

// function to get pid value in given file
func PidFileInfo(pidFile string) (string) {
  var currLine string
  fp, err := os.Open(pidFile)
  ExitIfError(err)
  defer fp.Close()
  scanner := bufio.NewScanner(fp)
  for scanner.Scan() {
    currLine = scanner.Text()
    break
  }
  return currLine
}

// function to print message in the given color
func PrintColor(messageColor string, messsage string) {
  fmt.Printf("%s%s%s", messageColor, messsage, Off)
}

func CreateColorMsg(messageColor string, messsage string) string {
  msg := fmt.Sprintf("%s%s%s", messageColor, messsage, Off)
  return msg
}

// function print the header after doing a clear screen
func PrintHeader(messageColor string, messsage string, clear bool) {
  if clear {
    fmt.Printf("%s", clearScreen)
  }
  fmt.Printf("\n\t%s %s %s %s %s\n\n", messageColor, HEADER, messsage, HEADER, Off)
}

// function to print error messsage in blue
func PrintBlue(message string) {
  PrintColor(Blue, message )
}

// function to print error messsage in cyan
func PrintCyan(message string) {
  PrintColor(Cyan, message )
}

// function to print error messsage in green
func PrintGreen(message string) {
  PrintColor(Green, message )
}

// function to print error messsage in purple
func PrintPurple(message string) {
  PrintColor(Purple, message )
}

// function to print error messsage in red
func PrintRed(message string) {
  PrintColor(Red, message )
}

// function to print error messsage in yellow
func PrintYellow(message string) {
  PrintColor(Yellow, message )
}

// function spepartion line
func PrintLine(lineColor string) {
  fmt.Printf("\n\t%s %s %s\n\n", lineColor, LINE, Off)
}

// function spepartion line
func PrintLinePurple() {
  fmt.Printf("\n\t%s %s %s\n\n", Purple, LINE, Off)
}


// function to clear the screen
func ClearScreen() {
  fmt.Printf("%s", clearScreen)
}

// function to read input from keyboard with a message before reading the input
func ReadInput(messageColor string, messsage string) string {
  reader := bufio.NewReader(os.Stdin)
  fmt.Printf("%s%s%s", messageColor, messsage, Off)
    data, err := reader.ReadString('\n')
  ExitIfError(err)
  return data
}

// function to check if the system is the given OS
// example : darwin, linux
func IsOS(osName string) (string, bool) {
  if runtime.GOOS == osName {
    return runtime.GOOS, true
  }
  return runtime.GOOS, false
}

// function to check if the given pid is running and is running the given command
func IsRunning(process string, pid int) bool {
  procInfo, err := ps.FindProcess(pid)
  ExitIfError(err)
  if procInfo == nil {
    return false
  }
  if procInfo.Executable() == process && procInfo.Pid() == pid {
    return true
  }
  return false
}

// function to create the lock file and make sure there
// no same process currently running
func LockIT(lockFile string, pid int, progName string) {
  if status, _ := LockFile(lockFile, pid, 1); status == false {
    pidValue, _ := strconv.Atoi(PidFileInfo(lockFile))
    // check if there is a running process already
    if IsRunning(progName, pidValue) {
      // fmt.Printf("There is already a process %s running, aborting\n", progName)
      PrintColor(
        Red,
        fmt.Sprintf("Error: There is already a process %s running, aborting...\n", progName))
      os.Exit(1)
    }
    fmt.Printf("Lock file %s exist, but not process with the pid %d is running...\n", lockFile, pidValue)
    fmt.Printf("Removing the lock file %s\n", lockFile)
    LockFile(lockFile, pid, 2)
  }
  // install interrupt handler
  SignalHandler(lockFile, pid)
}

// function to release remove the lockfile
func ReleaseIT(lockFile string, pid int) {
    LockFile(lockFile, pid, 2)
}

// function to return file information
func FileInfo(file string) (string, string, error) {
  fileInfo, err := os.Stat(file)
  if err != nil {
    return "", "", err
  }
  filePermissions := fmt.Sprintf("%04o", fileInfo.Mode().Perm())

  fileState := fileInfo.Sys().(*syscall.Stat_t)
  userId := strconv.FormatUint(uint64(fileState.Uid), 10)
  userName, err := user.LookupId(userId)
  if err != nil {
    return "", filePermissions, err
  }
  return userName.Username, filePermissions, err
}

// function to make sure the given file has the proper owner and persmission
func CheckFileSettings(file string, owners, permissions []string) bool {
  var ok int = 0
  fileOwner, filePerm, err := FileInfo(file)
  if err != nil {
    return false
  }

  usersCount := len(owners) -1
  for ownerCount, owner := range owners {
    if fileOwner == owner {
      break
    }
    if ownerCount == usersCount {
      PrintRed(fmt.Sprintf("Error: The file %s is not own by %s, is own by %s.\n", file, owner, fileOwner))
      ok++
    }
  }

  modCount := len(permissions) - 1
  for perCount, filePermission := range permissions {
        if filePermission == filePerm {
            break
        }
    if perCount == modCount {
      PrintRed(fmt.Sprintf("Error: The file %s permission are to wide open %s.\n", file, filePerm))
      ok++
    }
  }
  if ok != 0 {
    return false
  }
  return true
}

// the spinner class
func Spinner(speed time.Duration) *Wheel {
  return &Wheel{
    StopChannel:  make(chan int),
    Speed:      speed,
  }
}

// run the spinner until a StopChannel was reveived
func (w *Wheel) Run() {
  var spinstr = "| / - \\"
  fmt.Fprintf(os.Stdout, "\x1b[?25l")
  for {
    select {
      case <-w.StopChannel:
        w.StopChannel <- 1
        return
        ;;
      default:
        break
        ;;
    }
    for _, c := range spinstr {
      fmt.Printf("[%c]", c)
      time.Sleep(w.Speed * time.Millisecond)
      fmt.Printf("\b\b\b\b\b\b\x1b[1;35m")
    }
  }
}

// stop the spinner by sending 1 to the StopChannel 
func (w *Wheel) Stop() {
  // shutdown the spinner
  w.StopChannel <- 1
  <-w.StopChannel
  close(w.StopChannel)
  // show cursor
  fmt.Fprint(os.Stdout, "\x1b[?25h")
  fmt.Fprint(os.Stdout, clearLine)
}

// function to get value since epoch time in given argument
func GetEpoch(base string) int64  {
  current_time := time.Now().Unix()

  switch base {
    case "weeks": return (current_time / (86400 * 7)) 
    case "days": return (current_time / 86400)
    case "hours": return ((current_time / 86400) * 24)
    case "minutes": return (current_time / 86400) * (60 * 24)
    case "seconds": return current_time
    default: return current_time
  }
}

// function to make epoch time to human readable value
func GetReadableEpoch(epoch int64) (time.Time, string) {
  value := time.Unix(epoch, 0)
  return value, value.Format(USformat)
}

// function to generate a radom string
func GenerateRandom(useSpecialChar bool, length int) string {
  charSet := alphaNumeric
  if useSpecialChar {
    charSet = charSet + specialChars
  }

  random := make([]byte, length)
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for cnt, _ := range random {
    random[cnt] = charSet[r.Intn(len(charSet))]
  }
  return string(random)
}

// function to check for given y/n
func GetYN(keyboardInput string, defaultReturn bool) bool {
  yesSelection := []string{"yes", "y"}

  if keyboardInput == "" {
    return defaultReturn
  }
  for _, selected := range yesSelection {
        if strings.EqualFold(string(selected), keyboardInput) {
            return true
        }
    }
    return defaultReturn
}

func InList(slice []string, val string) (bool) {
  for _, item := range slice {
    if item == val {
      return true
    }
  }
  return false
}

func TheEnd () {
  PrintColor(Green, "\tEnjoy a cuppa of hot coffee ☕️   / 🥃\n")
  PrintColor(Green, "\tThe End\n")
}

// function to log user's ldap record or password 
func RecordPassword(passFile string, data map[string]string) (bool, string) {
  fp, err := os.OpenFile(passFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
  if err != nil {
    return false, err.Error()
  }
  current_time := time.Now().Format("Date: Feb 02 2006, Time: 15:04:05")
  fmt.Fprintf(fp, "\t%s\n", LINE)
  fmt.Fprintf(fp, "\t%s\n", current_time)
  fmt.Fprintf(fp, "\tThe user %s password is set to %s\n", data["user"], data["password"])
  fmt.Fprintf(fp, "\tEnjoy a cuppa of hot coffee ☕️   / 🥃\n")
  fmt.Fprintf(fp, "\tThe End\n")
  fp.Close()
  return true, ""
}


func Tidy(s string) string {
  return strings.TrimSpace(removeLines.Replace(s))
}

func getFrame(skipFrames int) runtime.Frame {
    // We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
    targetFrameIndex := skipFrames + 2

    // Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
    programCounters := make([]uintptr, targetFrameIndex+2)
    n := runtime.Callers(0, programCounters)

    frame := runtime.Frame{Function: "unknown"}
    if n > 0 {
        frames := runtime.CallersFrames(programCounters[:n])
        for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
            var frameCandidate runtime.Frame
            frameCandidate, more = frames.Next()
            if frameIndex == targetFrameIndex {
                frame = frameCandidate
            }
        }
    }
    return frame
}

// funcution to returns the caller of the function that called it
func FuncName() string {
        // Skip GetCallerFunctionName and the function to get the caller of
    caller := getFrame(2).Function
    if len(caller) > 0 {
      endIndex := strings.LastIndex(caller, ".")
      return caller[endIndex+1:]
    }
    return ""
}
