//go:generate goversioninfo -icon=rspw.ico -platform-specific

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	logger "rspw/logger"

	"github.com/go-vgo/robotgo"
	"github.com/tobischo/gokeepasslib"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/ini.v1"
)

func main() {

	logger.InfoLogger.Println("Reading config file")
	cfg, err := ini.Load("config.ini")
	if err != nil {
		logger.ErrorLogger.Println("Failed to read config file: ", err)
		os.Exit(1)
	}

	// Checking if our RuneScape launcher is running
	logger.InfoLogger.Println("Searching for RuneScape PID")

	// Handling if RuneScape isn't running
	needsLaunched := false

	// Searching for RuneScape PID in case it's already running
	fpid, err := robotgo.FindIds("rs2client")
	if len(fpid) == 0 {
		logger.ErrorLogger.Println("PID not found, attempting to launch RuneScape")
		needsLaunched = true
	}

	if needsLaunched == true {

		// RuneScape is not running so launch it
		logger.InfoLogger.Println("Current installPath value from config:", cfg.Section("config").Key("installPath").String())
		cmd := exec.Command(cfg.Section("config").Key("installPath").String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// cmd.Start() does not wait for the process to end (cmd.Run() does)
		rs := cmd.Start()
		if rs != nil {
			logger.ErrorLogger.Println("Failed to launch RuneScape, ", rs)
			os.Exit(1)
		}

		logger.InfoLogger.Println("RuneScape launched, waiting for the application to load")

		// Reading waitTime value and using that value to wait for RuneScape to launch
		waitTime := cfg.Section("config").Key("waitTime").MustInt64()
		logger.InfoLogger.Println("Current waitTime value from config:", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)

		// Searching for pid again now that RuneScape is running
		fpid, err = robotgo.FindIds("rs2client")
		if len(fpid) == 0 {
			logger.ErrorLogger.Println("PID not found after launching RuneScape, exiting")
			os.Exit(1)
		}
	}

	// Grabbing the PID of the launcher
	// This might not be needed?
	logger.InfoLogger.Println("Application loaded, found RuneScape PID:", fpid)
	pidExist, err := robotgo.PidExists(fpid[0])
	if err != nil {
		logger.ErrorLogger.Println("Error retrieving PID from ", fpid)
	}

	// Grabbing the password from our KeePass db
	// Setting the RuneScape launcher as our active window
	// Typing RS password into the window
	if pidExist {
		logger.InfoLogger.Println("Initiating retrieval from KeePass")
		// TODO: Ensure rspw is the active window here
		// this way the user doesn't have to alt tab if RuneScape wasn't running when rspw was launched
		logger.InfoLogger.Println("Current databasePath value from config:", cfg.Section("config").Key("databasePath").String())
		rsPass := retrievePass(cfg.Section("config").Key("databasePath").String())
		logger.InfoLogger.Println("Data retrieved from KeePass")

		logger.InfoLogger.Println("Setting RuneScape as active window")
		robotgo.ActivePID(fpid[0])

		robotgo.TypeStr(rsPass)

		logger.InfoLogger.Println("Process completed, closing rspw")
	}
}

func retrievePass(databasePath string) (passOut string) {

	// Prompting for user password and hiding Stdin
	fmt.Print("password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		logger.ErrorLogger.Println("Error reading password from stdin")
	}

	password := string(bytePassword)

	file, err := os.Open(databasePath)
	if err != nil {
		logger.ErrorLogger.Println("Error opening database")
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	_ = gokeepasslib.NewDecoder(file).Decode(db)

	db.UnlockProtectedEntries()

	// Password entry should be the first folder under root
	entry := db.Content.Root.Groups[0].Groups[0].Entries[0]
	return entry.GetPassword()
}
