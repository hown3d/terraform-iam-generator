package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	messageChan chan CsmMessage
	conn        *net.UDPConn
}

type CsmMessage struct {
	API     string `json:"Api"`
	Service string `json:"Service"`
}

func listen(port int) (*net.UDPConn, error) {
	return net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	})
}

func NewServerAndListen(messageChan chan CsmMessage) (Server, error) {
	conn, err := listen(31000)
	if err != nil {
		return Server{}, fmt.Errorf("listening with udp protocol on port 31000: %w", err)
	}
	return Server{
		conn:        conn,
		messageChan: messageChan,
	}, nil
}

func (s Server) Read() {
	b := make([]byte, 2048)
	for {
		n, _, err := s.conn.ReadFrom(b)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				close(s.messageChan)
				return
			}
			log.Println(fmt.Errorf("reading from udp connection: %w", err))
			return
		}

		var msg CsmMessage
		if err = json.Unmarshal(b[:n], &msg); err != nil {
			log.Println(fmt.Errorf("unmarshaling metrics message: %w", err))
			continue
		}
		s.messageChan <- msg
	}
}

func (s Server) Stop() error {
	return s.conn.Close()
}
