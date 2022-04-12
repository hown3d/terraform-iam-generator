package metrics

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Read(t *testing.T) {
	type fields struct {
		conn udpConn
	}
	tests := []struct {
		name   string
		fields fields
		want   []CsmMessage
	}{
		{
			name: "5 messages",
			fields: fields{
				conn: &mockUdpConn{
					msgToMarshal: CsmMessage{
						API:     "testApiCall",
						Service: "s3",
					},
					totalMsgs: 5,
				},
			},
			want: []CsmMessage{
				{
					API:     "testApiCall",
					Service: "s3",
				},
				{
					API:     "testApiCall",
					Service: "s3",
				},
				{
					API:     "testApiCall",
					Service: "s3",
				},
				{
					API:     "testApiCall",
					Service: "s3",
				},
				{
					API:     "testApiCall",
					Service: "s3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// buffer channel because we only read after all messages are sent
			messageChan := make(chan CsmMessage, len(tt.want))
			s := Server{
				messageChan: messageChan,
				conn:        tt.fields.conn,
			}
			// will only return after the channel is closed, so safe to retrieve everything from the channel without blocking
			s.Read()
			for msg := range messageChan {
				assert.Contains(t, tt.want, msg)
			}
		})
	}
}

type mockUdpConn struct {
	msgToMarshal CsmMessage
	totalMsgs    int
}

func (u *mockUdpConn) ReadFrom(b []byte) (int, net.Addr, error) {
	var addr *net.UDPAddr
	// close network connection when the maximum number of msgs where send
	if u.totalMsgs == 0 {
		return 0, addr, net.ErrClosed
	}
	data, err := json.Marshal(u.msgToMarshal)
	if err != nil {
		return 0, addr, err
	}
	n := copy(b, data)
	u.totalMsgs--
	return n, addr, nil
}

func (u *mockUdpConn) Close() error {
	return nil
}
