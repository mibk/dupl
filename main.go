package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/output/text"
	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
)

var (
	dir        = "."
	threshold  = flag.Int("t", 15, "minimum token sequence as a clone")
	serverPort = flag.String("serve", "", "run server at port")
	addrs      AddrList
)

type AddrList []string

func (l *AddrList) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *AddrList) Set(val string) error {
	*l = append(*l, val)
	return nil
}

func init() {
	flag.Var(&addrs, "c", "connect to the given 'addr:port'")
}

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	if len(addrs) != 0 {
		runClient()
	} else if *serverPort != "" {
		runServer()
	} else {
		schan := job.CrawlDir(dir)
		bchan := make(chan *job.Batch)
		go func() {
			for seq := range schan {
				bchan <- job.NewBatch("", seq)
			}
			close(bchan)
		}()
		t, done := job.BuildTree(bchan)
		<-done
		printDupls(t, new(LocalFileReader))
	}
}

type Scanner struct {
	ch chan []*syntax.Node
}

type Response struct {
	Seq []*syntax.Node
	Ok  bool
}

func (s *Scanner) Next(ignore bool, r *Response) error {
	r.Seq, r.Ok = <-s.ch
	return nil
}

func (s *Scanner) ReadFile(filename string, content *[]byte) error {
	c, err := ioutil.ReadFile(filename)
	*content = c
	return err
}

func runServer() {
	server := new(Scanner)
	rpc.Register(server)

	l, err := net.Listen("tcp", ":"+*serverPort)
	if err != nil {
		log.Fatal("error:", err)
	}
	log.Println("server started")

	for {
		if conn, err := l.Accept(); err != nil {
			log.Fatal(err.Error())
		} else {
			log.Println("connection accepted")
			server.ch = make(chan []*syntax.Node)
			go func() {
				schan := job.CrawlDir(dir)
				for seq := range schan {
					server.ch <- seq
				}
				close(server.ch)
			}()
			rpc.ServeConn(conn)
			log.Println("done")
		}
	}
}

type RemoteFileReader struct {
	clients map[string]*rpc.Client
}

func (r *RemoteFileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	client, ok := r.clients[node.Addr]
	if !ok {
		panic("client '" + node.Addr + "' is not present")
	}
	var content []byte
	err := client.Call("Scanner.ReadFile", node.Filename, &content)
	return content, err
}

func runClient() {
	clients := make(map[string]*rpc.Client)
	for _, addr := range addrs {
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
			if !reply.Ok {
				delete(tempClients, addr)
			}
			bchan <- job.NewBatch(addr, reply.Seq)
		}
	}
	close(bchan)

	<-done
	printDupls(t, &RemoteFileReader{clients})
}

type char int

func (c char) Val() int {
	return int(c)
}

type LocalFileReader struct{}

func (r *LocalFileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	return ioutil.ReadFile(node.Filename)
}

func printDupls(t *suffixtree.STree, fr text.FileReader) {
	// finish stream
	t.Update(char(-1))

	// print clones
	printer := text.NewPrinter(os.Stdout, fr)
	mchan := t.FindDuplOver(*threshold)
	cnt := 0
	for m := range mchan {
		if dups := syntax.FindSyntaxUnits(t, m, *threshold); len(dups) != 0 {
			printer.Print(dups)
			cnt++
		}
	}
	fmt.Printf("\nFound total %d clone groups.\n", cnt)
}
