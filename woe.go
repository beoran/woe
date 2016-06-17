package main

// import "fmt"
import "github.com/beoran/woe/server"
import "github.com/beoran/woe/monolog"
import "github.com/beoran/woe/raku"
import "os"
import "os/exec"
import "flag"
import "fmt"

type serverLogLevels []string

var server_loglevels serverLogLevels = serverLogLevels{
	"FATAL", "ERROR", "WARNING", "INFO",
}

/* Command line flags. */
var server_mode = flag.Bool("s", false, "Run in server mode")
var raku_mode = flag.Bool("r", false, "Run in Raku interpreter mode")
var server_tcpip = flag.String("l", ":7000", "TCP/IP Address where the server will listen")
var enable_logs = flag.String("el", "FATAL,ERROR,WARNING,INFO", "Log levels to enable")
var disable_logs = flag.String("dl", "", "Log levels to disable")

func enableDisableLogs() {
	monolog.EnableLevels(*enable_logs)
	monolog.EnableLevels(*disable_logs)
}

/* Need to restart the server or not? */
var server_restart = true

func runServer() (status int) {
	monolog.Setup("woe.log", true, false)
	defer monolog.Close()
	enableDisableLogs()
	monolog.Info("Starting WOE server...")
	monolog.Info("Server will run at %s.", *server_tcpip)
	woe, err := server.NewServer(*server_tcpip)
	if err != nil {
		monolog.Error("Could not initialize server!")
		monolog.Error(err.Error())
		panic(err)
	}
	monolog.Info("Server at %s init ok.", *server_tcpip)
	defer woe.Close()
	status, err = woe.Serve()
	if err != nil {
		monolog.Error("Error while running WOE server!")
		monolog.Error(err.Error())
		panic(err)
	}
	monolog.Info("Server shut down without error indication.", *server_tcpip)
	return status
}

func runSupervisor() (status int) {
	monolog.Setup("woe.log", true, false)
	defer monolog.Close()
	enableDisableLogs()
	monolog.Info("Starting WOE supervisor.")
	for server_restart {
		// wd  , _ := os.Getwd()
		exe := fmt.Sprintf("%s", os.Args[0])
		argp := fmt.Sprintf("-l=%s", *server_tcpip)
		argel := fmt.Sprintf("-el=%s", *enable_logs)
		argdl := fmt.Sprintf("-dl=%s", *disable_logs)
		cmd := exec.Command(exe, "-s=true", argp, argel, argdl)
		monolog.Info("Starting server %s at %s.", exe, *server_tcpip)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		monolog.Debug("Server command line: %s.", cmd.Args)
		err := cmd.Run()
		monolog.Info("Server at %s shut down.", *server_tcpip)
		if err != nil {
			monolog.Error("Server shut down with error %s!", err)
			server_restart = false
			return 1
		}
	}
	return 0
}

func runRaku() (status int) {
	lexer := raku.OpenLexer(os.Stdin)
	_ = lexer
	return 0
}

/* Woe can be run in supervisor mode (the default) or server mode (-s).
 * Server mode is the mode in which the real server is run. In supervisor mode,
 * woe runs a single woe server in server mode using os/exec. This is used to
 * be able to restart the server gracefully on recompile of the sources.
 */
func main() {
	defer func() {
		pani := recover()
		if pani != nil {
			monolog.Fatal("Panic: %s", pani)
			os.Exit(255)
		}
	}()

	flag.Parse()
	if *server_mode {
		os.Exit(runServer())
	} else if *raku_mode {
		os.Exit(runRaku())
	} else {
		os.Exit(runSupervisor())
	}
}
