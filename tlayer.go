package main

import (
	"fmt"
	"github.com/rbroggi/tlayer/connection"
	"github.com/rbroggi/tlayer/sender"
	"time"
)

func main() {

	fmt.Printf("Starting transport layer simulator\n")

	strToSend := "avoasidapsoivjasdofaseprnasvo[ispoiAWERO[IASDNVOAISF-EIAJWPIRNFVIUFHPIAUDSFDHG[QONRIOUHV[IADJGPOIAJERIONV[OADJF[OIAJ[OIASJERGO[IFNVO[ADNFVO[IASJVBO[IADJBH[OIADSJFOAI[SJTROI[RGB[IDJJKCNA,/CVMZXM,CNJLKGAHUOIQHEWURPYAWPIUTHY89Y408740TY8073T608GUEVHUIFDAKHGQ089374301984YRHKAHFGK;AJHDNFKJ;ANDF"
	fmt.Printf("%v\n", strToSend)

	data := []byte(strToSend)
	byteLength := len(data)
	segSize := 1
	segNum := byteLength / segSize
	connParam := connection.ConParam{WinSize: 3, SegSize: int32(segSize), SegNum: int32(segNum), RetransmissionTimeout: 100* time.Microsecond}
	connection := connection.NewMockConnection(connParam)
	fmt.Printf("Creating sender\n")
	sender := sender.NewSender(connection, connParam)
	fmt.Printf("Sending data\n")
	sender.SendData(data)
	fmt.Printf("Code finished\n")

}
