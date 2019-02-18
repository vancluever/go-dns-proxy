package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

const forwardingAddress = "8.8.8.8:53"

func main() {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	server := &dns.Server{PacketConn: conn}
	dns.HandleFunc(".", forwarder)

	log.Printf("listening on %s, press CTRL-C or send SIGTERM to exit", conn.LocalAddr())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Printf("%s received, shutting down.", s.String())
		if err := server.Shutdown(); err != nil {
			log.Printf("Error closing listener: %s", err)
		}
		os.Exit(0)
	}()

	if err := server.ActivateAndServe(); err != nil {
		log.Fatal(err)
	}
}

func forwarder(r dns.ResponseWriter, msg *dns.Msg) {
	defer r.Close()
	log.Printf("query from %s: %s", r.RemoteAddr(), msg)

	c := new(dns.Client)
	resp, rtt, err := c.Exchange(msg, forwardingAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("reply from %s (RTT=%d): %s", forwardingAddress, rtt, resp)
	r.WriteMsg(resp)
}
