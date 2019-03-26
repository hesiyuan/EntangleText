package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

// args in insert(args)
type InsertArgs struct {
	Char       uint8  // character to insert
	Identifier uint8  // position identifier of the char TODO:
	Clock      uint64 // value of logical clock at the issuing client
	Clientid   uint8
}

// args in put(args)
type DeleteArgs struct {
	Char       uint8  // character to delete, could be omitted
	Identifier uint8  // position identifier of the char to delete TODO:
	Clock      uint64 // value of logical clock at the issuing client
	Clientid   uint8
}

// args in disconnect(args)
type DisconnectArgs struct {
	Clientid uint8 // client id who voluntarilly quit the editor
}

// Reply from service for all the API calls above.
// This is actually not gonna be used.
type ValReply struct {
	Val string // value; depends on the call
}

type EntangleClient int

// Command line arg.
var numPeers uint8

//a slice holding peer ip addresses
var peerAddresses []string

// a slice hoding rpc service of peers
var peerServices []*rpc.Client

// a insert char message from a peer
func (ec *EntangleClient) Insert(args *InsertArgs, reply *ValReply) error {
	// TODO

	//testing
	fmt.Printf(string(args.Char))
	fmt.Println("remote insert")

	return nil
}

// a delete char message from a peer
func (ec *EntangleClient) Delete(args *DeleteArgs, reply *ValReply) error {
	// TODO

	return nil
}

// DISCONNECT from a peer.
func (ec *EntangleClient) Disconnect(args *DisconnectArgs, reply *ValReply) error {
	// TODO

	return nil
}

// Entangle client main loop.
func main() {
	// Parse args.
	usage := fmt.Sprintf("Usage: %s [ip:port] [N-clients] [ip1:port] ... [ipN:port]\n", os.Args[0])
	if len(os.Args) < 4 {
		fmt.Printf(usage)
		os.Exit(1)
	}

	ip_port := os.Args[1]
	arg, err := strconv.ParseUint(os.Args[2], 10, 8)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if arg == 0 {
		fmt.Printf(usage)
		fmt.Printf("\tnum-clients arg must be non-zero\n")
		os.Exit(1)
	}
	numPeers = uint8(arg)

	// Setup key-value store and register service.
	entangleClient := new(EntangleClient)
	rpc.Register(entangleClient)

	// listen first
	l, e := net.Listen("tcp", ip_port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// then dial
	peerAddresses = make([]string, len(os.Args)-3)
	peerServices := make([]*rpc.Client, len(os.Args)-3)
	for i := range peerAddresses {
		peerAddresses[i] = os.Args[i+3]
		// Connect to other peers via RPC.
		peerServices[i], err = rpc.Dial("tcp", peerAddresses[i])

		checkError(err)
	}

	var kvVal ValReply
	InsertArgs := InsertArgs{Char: 'c', Identifier: 12, Clientid: 1, Clock: 123}
	ticker := time.NewTicker(time.Duration(2) * time.Second)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("insert")
				val := peerServices[0].Call("EntangleClient.Insert", InsertArgs, &kvVal)
				fmt.Println(val)
			case <-quit:
				ticker.Stop()
				close(quit)
				fmt.Println("clock stopped")
				return
			}
		}
	}()

	// TODO: Enter servicing loop, like:

	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}
