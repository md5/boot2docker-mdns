package main

import (
    "os"
    "fmt"
    "net"
    "time"

    "github.com/armon/mdns"
)

const (
    mdnsQuietInterval = 100 * time.Millisecond
)

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "error in run: %v\n", err)
        os.Exit(1)
    }
}

func serviceIP() (net.IP, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return nil, err
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            return ipnet.IP, nil
        }
    }

    return nil, fmt.Errorf("Unable to determine service address")
}

func run() error {
    host, err := os.Hostname()
    if err != nil {
        return err
    }

    ip, err := serviceIP()
    if err != nil {
        return err
    }

    service := &mdns.MDNSService {
        Instance: host,
        Service:  "_docker._tcp",
        Addr:     ip,
        Port:     2375,
    }
    if err := service.Init(); err != nil {
        return err
    }

    server, err := mdns.NewServer(&mdns.Config { Zone: service })
    if err != nil {
        return err
    }
    defer server.Shutdown()

    var quiet <-chan time.Time

    for {
        select {
        case <-quiet:
            quiet = time.After(mdnsQuietInterval)
        }
    }

    return nil
}
