package remote

import (
	"log"
	"net/rpc"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
)

type FileReader struct {
	clients map[string]*rpc.Client
}

func NewFileReader(clients map[string]*rpc.Client) *FileReader {
	return &FileReader{clients}
}

func (r *FileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	client, ok := r.clients[node.Addr]
	if !ok {
		panic("client '" + node.Addr + "' is not present")
	}
	var content []byte
	err := client.Call("Scanner.ReadFile", node.Filename, &content)
	return content, err
}

func RunClient(addrs []string) (*suffixtree.STree, map[string]*rpc.Client) {
	clients := make(map[string]*rpc.Client)
	for _, addr := range addrs {
		if _, present := clients[addr]; present {
			// ignore duplicate addr
			continue
		}
		client, err := rpc.Dial("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		clients[addr] = client
	}
	log.Println("connection established")

	bchan := make(chan *job.Batch)
	t, done := job.BuildTree(bchan)

	tempClients := make(map[string]*rpc.Client)
	for addr, client := range clients {
		tempClients[addr] = client
	}

	for len(tempClients) > 0 {
		var reply Response
		for addr, client := range tempClients {
			err := client.Call("Scanner.Next", true, &reply)
			if err != nil {
				log.Fatal(err)
			}
			if reply.Done {
				delete(tempClients, addr)
			}
			bchan <- job.NewBatch(addr, reply.Seq)
		}
	}
	close(bchan)

	<-done
	return t, clients
}
