package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

// args in insert(args)
type InsertArgs struct {
	char       uint8  // character to insert
	identifier uint8  // position identifier of the char TODO:
	clock      uint64 // value of logical clock at the issuing client
}

// args in put(args)
type DeleteArgs struct {
	char       uint8  // character to delete, could be omitted
	identifier uint8  // position identifier of the char to delete TODO:
	clock      uint64 // value of logical clock at the issuing client
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

// a insert char message from a peer
func (ec *EntangleClient) Insert(args *InsertArgs, reply *ValReply) error {
	// TODO

	//testing
	fmt.Printf(string(args.char))

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
	usage := fmt.Sprintf("Usage: %s [ip:port] [num-clients]\n", os.Args[0])
	if len(os.Args) != 3 {
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
	l, e := net.Listen("tcp", ip_port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// TODO: Enter servicing loop, like:

	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}
