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
)

func main() {

	// TODO: Add more verbose logs
	// TODO: Launch RuneScape if not already running
	// TODO: Create a ini file to get paths from, this will help support non-default installs

	// Checking if our RuneScape launcher is running
	logger.InfoLogger.Println("Searching for RuneScape PID")

	// Handling if RuneScape isn't running
	needsLaunched := false

	fpid, err := robotgo.FindIds("rs2client")
	if len(fpid) == 0 {
		logger.ErrorLogger.Println("PID not found, attempting to launch RuneScape")
		needsLaunched = true
	}

	if needsLaunched == true {

		// RuneScape is not running so launch it
		// TODO: put this path in the ini file as it might not be default on every installation
		logger.InfoLogger.Println("Launching RuneScape")
		cmd := exec.Command("C:\\Program Files\\Jagex\\RuneScape Launcher\\RuneScape.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// cmd.Start() does not wait for the process to end (cmd.Run() does)
		rs := cmd.Start()
		if rs != nil {
			logger.ErrorLogger.Println("Failed to launch RuneScape, ", rs)
			os.Exit(1)
		}

		// Waiting for RuneScape to launch
		time.Sleep(3 * time.Second)

		// Searching for pid again now that RuneScape is running
		fpid, err = robotgo.FindIds("rs2client")
		if len(fpid) == 0 {
			logger.ErrorLogger.Println("PID not found after launching RuneScape, exiting")
			os.Exit(1)
		}
	}

	// Grabbing the PID of the launcher
	// This might not be needed?
	logger.InfoLogger.Println("Found RuneScape PID:", fpid)
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
		rsPass := retrievePass()

		logger.InfoLogger.Println("Setting RuneScape as active window")
		robotgo.ActivePID(fpid[0])

		robotgo.TypeStr(rsPass)
	}
}

func retrievePass() (passOut string) {

	// Prompting for user password and hiding Stdin
	fmt.Print("password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		logger.ErrorLogger.Println("Error reading password from stdin")
	}

	password := string(bytePassword)

	// TODO: Add this path to the ini file
	file, err := os.Open("C:\\Users\\matth\\Documents\\rs.kdbx")
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
