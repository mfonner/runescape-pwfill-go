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
    "golang.org/x/term" 
	"gopkg.in/ini.v1"
)

func main() {

	// TODO: Add log retention, value should be configurable in config.ini
    // TODO: Add RuneScape type to config file (osrs vs RS3)
    // This way, I can merge this all into master and programmitcally call the correct version

	logger.InfoLogger.Println("Reading config file")
	cfg, err := ini.Load("config.ini")
	if err != nil {
		logger.ErrorLogger.Println("Failed to read config file: ", err)
		os.Exit(1)
	}

    sec, err := cfg.GetSection("config")
    if err != nil {
        logger.ErrorLogger.Println("Error setting config scope to ", sec)
        os.Exit(1)
    }

	logger.InfoLogger.Println("Initiating retrieval from KeePass")
	logger.InfoLogger.Println("Current databasePath value from config:", sec.Key("databasePath").String())

	rsPass := retrievePass(sec.Key("databasePath").String())
	logger.InfoLogger.Println("Data retrieved from KeePass")

	// Checking if our RuneScape launcher is running
	logger.InfoLogger.Println("Searching for RuneScape PID")

	// Handling if RuneScape isn't running
	needsLaunched := false

	// Searching for RuneScape PID in case it's already running
	fpid, err := robotgo.FindIds("RuneLite")
	if len(fpid) == 0 {
		logger.ErrorLogger.Println("PID not found, attempting to launch RuneScape")
		needsLaunched = true
	}

	// RuneScape is not running so launch it
	if needsLaunched == true {

		logger.InfoLogger.Println("Current installPath value from config:", sec.Key("installPath").String())

        // The _JAVA_OPTIONS env variable moves the files that the OSRS launcher places in /home/$USER
        // This isn't needed, but helps keeps the user's home directory nice and tidy
         
        if sec.HasKey("javaEnv") && sec.Key("javaEnv").String() != "" {

            logger.InfoLogger.Println("Java options found for RuneLite, setting those before continuing.")

            rsJavaEnv := os.Setenv(sec.Key("javaEnv").String(), sec.Key("javaVal").String()) 

            if rsJavaEnv != nil {
                logger.ErrorLogger.Println("Failed to set Java environment variables before launch. Error: ", rsJavaEnv)
            }

            logger.InfoLogger.Println("Java environment variable set, launching RuneLite.")
        }


		cmd := exec.Command(sec.Key("installPath").String())

		// Suppressing the output of launching RuneLite
		cmd.Stdout = nil
		cmd.Stderr = nil

		// cmd.Start() does not wait for the process to end (cmd.Run() does)
		rs := cmd.Start()
		if rs != nil {
			logger.ErrorLogger.Println("Failed to launch RuneScape, ", rs)
			os.Exit(1)
		}

		logger.InfoLogger.Println("RuneScape launched, waiting for the application to load")

		// Reading waitTime value and using that value to wait for RuneScape to launch
		waitTime := sec.Key("waitTime").MustInt64()
		logger.InfoLogger.Println("Current waitTime value from config:", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)

		robotgo.KeyTap("enter")
		time.Sleep( 1 * time.Second)
		robotgo.TypeStr(rsPass)

        logger.InfoLogger.Println("Process completed, exiting.")

		// Searching for pid again now that RuneScape is running
		fpid, err = robotgo.FindIds("RuneLite")
		if len(fpid) == 0 {
			logger.ErrorLogger.Println("PID not found after launching RuneScape, exiting with error", err)
			os.Exit(1)
		}

		//logger.InfoLogger.Println("Application loaded, found RuneScape PID:", fpid)
	}

	//logger.InfoLogger.Println("RuneScape PID found:", fpid)

	// Grabbing the PID of the launcher
	// This might not be needed?
	//pidExist, err := robotgo.PidExists(fpid[0])
	//if err != nil {
	//	logger.ErrorLogger.Println("Error retrieving PID from ", fpid)
	//}

	// Grabbing the password from our KeePass db
	// Setting the RuneScape launcher as our active window
	// Typing RS password into the window
	// TODO: Handle if RuneLite is already running
	// Window switching isn't working yet in my testing
	//if pidExist {
    
    //   switch active window here

	//}
}

func retrievePass(databasePath string) (passOut string) {

	// Prompting for user password and hiding Stdin
	fmt.Print("password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
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
