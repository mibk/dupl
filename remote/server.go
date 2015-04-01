package remote

import (
	"io/ioutil"
	"log"
	"net"
	"net/rpc"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/syntax"
)

type Scanner struct {
	ch chan []*syntax.Node
}

type Response struct {
	Seq  []*syntax.Node
	Done bool
}

func (s *Scanner) Next(ignore bool, r *Response) error {
	seq, ok := <-s.ch
	r.Seq, r.Done = seq, !ok
	return nil
}

func (s *Scanner) ReadFile(filename string, content *[]byte) error {
	c, err := ioutil.ReadFile(filename)
	*content = c
	return err
}

func RunServer(port, dir string) {
	server := new(Scanner)
	rpc.Register(server)

	l, err := net.Listen("tcp", ":"+port)
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
