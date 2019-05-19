package sender

import (
	"fmt"
	"github.com/rbroggi/tlayer/connection"
	"github.com/rbroggi/tlayer/pkg"
)

func NewSender(connection connection.Connection, connectionParam connection.ConParam) Sender {
	return &sender{
		conn: connection,
		connParam: connectionParam,
	}
}

type Sender interface {
	SendData(data []byte)
}

type sender struct {
	conn connection.Connection
	connParam connection.ConParam
}

func (s *sender) SendData(data []byte) {

	fmt.Printf("Sender: splitting data in segments\n")
	packs := segmentData(data, s.connParam)
	fmt.Printf("Data splited into %v packets\n", len(packs))

	fmt.Printf("Sender: Starting sending packets\n")
	for i, pkgToSend := range packs {
		s.conn.SendData(&pkgToSend)
		if (int32(i+1) % s.connParam.WinSize) == 0 {
			ack, err := s.conn.ReceiveAck(s.connParam.RetransmissionTimeout)
			for err != nil {
				fmt.Printf("Sender: no acknowdedge received. Retransmitting\n")
				ack, err = s.conn.ReceiveAck(s.connParam.RetransmissionTimeout)
			}
			fmt.Printf("Sender: acknowledge received %v\n", ack.SegIDX)
		}
	}
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func segmentData(data []byte, conParam connection.ConParam) []pkg.Pkg {
	pkgs := make([]pkg.Pkg, conParam.SegNum)

	for i := 0; i < int(conParam.SegNum); i++ {
		currPkg := pkg.Pkg{
			SegNum: int32(i),
			Data: data[i * int(conParam.SegSize): min(len(data)-1, (i+1)*int(conParam.SegSize))],
		}
		pkgs[i] = currPkg
	}
	return pkgs

}


