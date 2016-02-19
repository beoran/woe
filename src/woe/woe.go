package main

// import "fmt"
import "github.com/beoran/woe/server"
import "github.com/beoran/woe/monolog"
import "os"
import "os/exec"
import "flag"
import "fmt"

/* Command line flags. */
var server_mode  = flag.Bool("s", false, "Run in server mode");
var server_tcpip = flag.String("l", ":7000", "TCP/IP Address where the server will listen");
 
/* Need to restart the server or not? */
var server_restart = true

func runServer() (status int) {
    monolog.Setup("woe.log", true, false)
    defer monolog.Close()
    monolog.Info("Starting WOE!")
    monolog.Info("Server runs at %s!", *server_tcpip)
    woe, err := server.NewServer(*server_tcpip)
    if err != nil {
        monolog.Error(err.Error())
        panic(err)
    }
    defer woe.Close()
    status, err = woe.Serve()
    if err != nil {
        monolog.Error(err.Error())
        panic(err)
    }
    return status
}



func runSupervisor() (status int) {
    monolog.Setup("woe.log", true, false)
    defer monolog.Close()
    monolog.Info("Starting WOE supervisor!")
    for (server_restart) {
        // wd  , _ := os.Getwd()
        exe  := fmt.Sprintf("%s", os.Args[0]) 
        argp := fmt.Sprintf("-l=%s", *server_tcpip)
        cmd  := exec.Command(exe, "-s=true", argp)
        monolog.Info("Starting server %s at %s!", exe, *server_tcpip)
        cmd.Stderr = os.Stderr
        cmd.Stdout = os.Stdout
        err  := cmd.Run()
        monolog.Info("Server at %s shut down!", *server_tcpip)
        // monolog.Info("Server output: %s!", out);
        if (err != nil ) { 
            monolog.Error("Server shut down with error %s!", err)
            server_restart = false;
            return 1
        }
    }
    return 0
}



/* Woe can be run in supervisor mode (the default) or server mode (-s).
 * Server mode is the mode in which the real server is run. In supervisor mode, 
 * woe runs a single woe server in server mode using os/exec. This is used to 
 * be able to restart the server gracefully on recompile of the sources. 
 */
func main() {
    flag.Parse()
    if *server_mode {
        os.Exit(runServer())
    } else {
        os.Exit(runSupervisor())
    }
}

