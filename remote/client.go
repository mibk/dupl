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
	t, built := job.BuildTree(bchan)

	clientCnt := len(clients)
	cchan := make(chan *rpc.Call, clientCnt)
	done := make(chan bool)
	for addr, client := range clients {
		addr, client := addr, client
		go func() {
			for {
				var reply Response
				call := client.Go("Scanner.Next", true, &reply, cchan)
				<-call.Done
				if call.Error != nil {
					log.Fatal(call.Error)
				}
				if reply.Done {
					clientCnt--
					if clientCnt == 0 {
						close(cchan)
						done <- true
					}
					return
				}
				bchan <- job.NewBatch(addr, reply.Seq)
			}
		}()
	}
	<-done
	close(bchan)

	<-built
	return t, clients
}
