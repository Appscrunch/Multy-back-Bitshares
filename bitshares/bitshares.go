/*
 * Copyright 2018 Idealnaya rabota LLC
 * Licensed under Multy.io license.
 * See LICENSE for details
 */

package bitshares

import (
	"log"

	"github.com/denkhaus/bitshares/api"
	pb "github.com/Appscrunch/Multy-Back-Bitshares/proto"
	"github.com/denkhaus/bitshares/types"
	"time"
)

const (
	SubscriberId  = 5
	// BlockInterval is hardcoded
	// because get_global_properties not implemented yet
	BlockInterval =  time.Second * 3
)

// Server is a struct for interaction with bitshares chain
type Server struct {
	client           api.BitsharesAPI
	account          string
	activeKey        string
	TrackedAddresses map[string]bool
	BalanceChangedCh chan *pb.Balance
	NewBlockCh chan *pb.Block
}

// NewServer initializes and validates new api struct
// and runs chain monitoring loop
func NewServer(ws, rpc, account, key string) (*Server, error) {
	cli := api.New(ws, rpc)
	err := cli.Connect()
	if err != nil {
		log.Printf("connect: %s", err)
		return nil, err
	}

	cli.OnError(func(e error) {
		log.Printf("crap: %s", e)
	})

	log.Println("new client")
	s := &Server{
		client:           cli,
		account:          account,
		activeKey:        key,
		TrackedAddresses: make(map[string]bool),
	}

	//cli.SetDebug(true)
	err = cli.SetSubscribeCallback(SubscriberId, false)
	if err != nil {
		log.Printf("subscribe: %s", err)
		return nil, err
	}
	//cli.SetDebug(true)
	log.Printf("SpaceTypeImplementation %d.%d", types.SpaceTypeImplementation, types.ObjectTypeDynamicGlobalProperty)
	obj, err := s.client.GetObjects(types.NewGrapheneID("2.1.0"))
	log.Println(obj)
	if err != nil {
		log.Printf("get object: %s", err)
		return nil, err
	}
	cli.SetDebug(true)
	err = cli.OnNotify(SubscriberId, newBlockHandler)
	if err != nil {
		log.Printf("notify: %s", err)
		return nil, err
	}

	//cli.SetKeys(&client.Keys{AKey: []string{key}})

	s.BalanceChangedCh = make(chan *pb.Balance)
	s.NewBlockCh = make(chan *pb.Block)

	//go s.ProcessLoop(0)

	return s, nil
}

func newBlockHandler(msg interface{}) error{
	log.Printf("new block %v", msg)
	return nil
}

// NewBlockLoop checks for new blocks and send them to chans
// start is a number of starting block for iteration
// if start is 0, using head_block_number
func (s *Server) ProcessLoop(start uint32) {
	blockNum := start

	for {
		props, err := s.GetDynamicGlobalProperties()
		if err != nil {
			log.Printf("get global properties: %s", err)
			time.Sleep(time.Duration(BlockInterval) * time.Second)
			continue
		}
		if blockNum == 0 {
			blockNum = uint32(props.HeadBlockNumber)
		}
		// maybe LastIrreversibleBlockNum, cause possible microforks
		if uint32(props.HeadBlockNumber)-blockNum > 0 {
			log.Printf("new block: %d", blockNum + 1)
			block, err := s.client.GetBlock(uint64(blockNum + 1))
			if err != nil {
				log.Printf("get block: %s", err)
				time.Sleep(time.Duration(BlockInterval) * time.Second)
				continue
			}
			msg := makeBlock(block)
			select {
			case s.NewBlockCh <- &msg:
				// process block, now its only balance change check
				go s.processBalance(block)
			}
			blockNum++
		} else {
			time.Sleep(time.Duration(BlockInterval) * time.Second)
		}
	}
}
//
//
//
