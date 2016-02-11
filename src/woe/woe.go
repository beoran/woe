package main

// import "fmt"
import "github.com/beoran/woe/server"
import "github.com/beoran/woe/monolog"



func main() {
    monolog.Setup("woe.log", true, false)
    defer monolog.Close()
    monolog.Info("Starting WOE!")
    monolog.Info("Server runs at port %d!", 7000)
    woe, err := server.NewServer(":7000")
    if err != nil {
        monolog.Error(err.Error())
        panic(err)
    }
    defer woe.Close()
    woe.Serve()
}

