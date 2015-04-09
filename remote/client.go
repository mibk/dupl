package remote

import (
	"log"
	"math"
	"net/rpc"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/syntax"
)

type router struct {
	Slots map[int][]*rpc.Client
}

func newRouter() *router {
	return &router{make(map[int][]*rpc.Client)}
}

func (r *router) AddClient(slotId int, client *rpc.Client) {
	if _, ok := r.Slots[slotId]; !ok {
		r.Slots[slotId] = make([]*rpc.Client, 0, 1)
	}
	r.Slots[slotId] = append(r.Slots[slotId], client)
}

func RunClient(addrs []string, threshold int, dir string) <-chan [][]*syntax.Node {
	addrClients := make(map[string]*rpc.Client)
	clients := make([]*rpc.Client, 0, len(addrs))
	for _, addr := range addrs {
		if _, present := addrClients[addr]; present {
			// ignore duplicate addr
			continue
		}
		client, err := rpc.Dial("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		clients = append(clients, client)
		addrClients[addr] = client
	}
	log.Println("connection established")

	// count number of batches for the given number of clients
	batchCnt := int(math.Ceil(math.Sqrt(2*float64(len(clients))+0.25) + 0.5))

	router := newRouter()
	id := 0
	for i := 0; i < batchCnt-1; i++ {
		for j := i + 1; j < batchCnt; j++ {
			router.AddClient(i, clients[id])
			router.AddClient(j, clients[id])
			id = (id + 1) % len(clients)
		}
	}

	schan := job.CrawlDir(dir)
	nodesChan := make(chan [][]*syntax.Node)
	go func() {
		batch := 0
		for seq := range schan {
			for _, client := range router.Slots[batch] {
				err := client.Call("Dupl.UpdateTree", seq, nil)
				if err != nil {
					log.Fatal(err)
				}
			}
			batch = (batch + 1) % batchCnt
		}

		for len(addrClients) > 0 {
			for addr, client := range addrClients {
				var reply Response
				err := client.Call("Dupl.NextMatch", threshold, &reply)
				if err != nil {
					log.Fatal(err)
				}
				if reply.Done {
					delete(addrClients, addr)
					continue
				}
				nodesChan <- reply.Match
			}
		}
		close(nodesChan)
	}()
	return nodesChan
}
