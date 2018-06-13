/*
 * Copyright 2018 Idealnaya rabota LLC
 * Licensed under Multy.io license.
 * See LICENSE for details
 */

package main

import (
	"fmt"
	"log"
	//"net"
	"os"

	"github.com/Appscrunch/Multy-Back-Bitshares/bitshares"
	//pb "github.com/Appscrunch/Multy-Back-Steemit/proto"
	"github.com/urfave/cli"
	//"google.golang.org/grpc"
	"time"
)

var (
	commit    string
	branch    string
	buildtime string
)

const (
	VERSION = "v0.2"
)

type Server struct {
	*bitshares.Server
}

// NewServer constructs new server handler
func NewServer(ws, rpc, account, key string) (*Server, error) {
	a, err := bitshares.NewServer(ws, rpc, account, key)
	return &Server{a}, err
}

func run(c *cli.Context) error {
	//// check net arguement
	//network := c.String("net")
	//if network != "test" && network != "steem" {
	//	return cli.NewExitError(fmt.Sprintf("net must be \"steem\" or \"test\": %s", network), 1)
	//}

	server, err := NewServer(
		c.String("ws"),
		c.String("rpc"),
		c.String("account"),
		c.String("key"),
	)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("cannot init server: %s", err), 2)
	}
	log.Println("new server")

	//DEBUG
	server.TrackedAddresses["bit-vovapi"] = true

	for i := 0; i < 1000; i++ {
		time.Sleep(time.Second)
	}

	//addr := fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))

	// init gRPC server
	//lis, err := net.Listen("tcp", addr)
	//if err != nil {
	//	return cli.NewExitError(fmt.Sprintf("failed to listen: %s", err), 2)
	//}
	// Creates a new gRPC server
	//s := grpc.NewServer()
	//pb.RegisterNodeCommunicationsServer(s, server)
	//log.Printf("listening on %s", addr)
	//return cli.NewExitError(s.Serve(lis), 3)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "multy-bitshares"
	app.Usage = `bitshares node gRPC API for Multy backend`
	app.Version = fmt.Sprintf("%s (commit: %s, branch: %s, buildtime: %s)", VERSION, commit, branch, buildtime)
	app.Author = "vovapi"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Usage:  "hostname to bind to",
			EnvVar: "MULTY_BITSHARES_HOST",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "port",
			Usage:  "port to bind to",
			EnvVar: "MULTY_BITSHARES_PORT",
			Value:  "8080",
		},
		cli.StringFlag{
			Name:   "ws",
			Usage:  "node websocket address",
			EnvVar: "MULTY_BITSHARES_WS",
		},
		cli.StringFlag{
			Name: "rpc",
			Usage: "node rpc address",
			EnvVar: "MULTY_BITSHRES_RPC",
		},
		//cli.StringFlag{
		//	Name:   "net",
		//	Usage:  `network: "steem" for mainnet or "test" for testnet`,
		//	EnvVar: "MULTY_STEEM_NET",
		//	Value:  "test",
		//},
		cli.StringFlag{
			Name:   "account",
			Usage:  "bitshares account for user registration",
			EnvVar: "MULTY_BITSHARES_ACCOUNT",
		},
		cli.StringFlag{
			Name:   "key",
			Usage:  "active key for specified user for user registration",
			EnvVar: "MULTY_BITSHARES_KEY",
		},
	}
	app.Action = run
	app.Run(os.Args)
}
