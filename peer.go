package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				fmt.Println("this command is what we got =>", value)
				switch value.String() {
				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of variables for get command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}
				case CommandSET:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of variables for set command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}
				case CommandHELLO:
					cmd := HelloCommand{
						value: v.Array()[1].String(),
					}
					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}
				default:
					fmt.Println("got unknown command", v.Array())
				}
			}
		}
	}
	return nil
}
