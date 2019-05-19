package connection

import (
	"errors"
	"fmt"
	"github.com/rbroggi/tlayer/acknowledge"
	"github.com/rbroggi/tlayer/pkg"
	"time"
)

//NewMockConnection returns a new mockConnection with
// a receiver binded and in listening state
func NewMockConnection(conParam ConParam) Connection {
	c := &mockConnection{
		segmentTopic: make(chan *pkg.Pkg),
		ackTopic: make(chan acknowledge.Acknowledge),
	}

	//create receiver and start listening in separate routine
	rec := NewDumpReceiver(conParam)
	go rec.Listen(c)

	return c
}

type Connection interface {
	// To be used by sender

	SendData(data *pkg.Pkg)
	// Waits for timeout to receive an acknowledge else return error
	ReceiveAck(timeout time.Duration) (*acknowledge.Acknowledge, error)

	// To be used by receiver

	// Blocking method
	Receive() *pkg.Pkg
	SendAck(ack acknowledge.Acknowledge)
	Close()
}

type mockConnection struct {
	segmentTopic chan *pkg.Pkg
	ackTopic     chan acknowledge.Acknowledge
}


func (c *mockConnection) SendData(data *pkg.Pkg) {
	c.segmentTopic <- data
}

func (c *mockConnection) ReceiveAck(timeout time.Duration) (*acknowledge.Acknowledge, error) {
	select {
	case ret := <-c.ackTopic:
		return &ret, nil
	case <-time.After(timeout):
		return nil, errors.New("no acknowledge within timeout, retransmitting")
	}
}

func (c *mockConnection) Receive() *pkg.Pkg {
	return <-c.segmentTopic
}

func (c *mockConnection) SendAck(ack acknowledge.Acknowledge) {
	c.ackTopic <- ack
}

func (c *mockConnection) Close() {
	closingAck := acknowledge.Acknowledge{
		SegIDX: -1,
	}
	c.ackTopic <- closingAck
}

type ConParam struct {
	// #segments by window
	WinSize int32
	// #bytes by segment
	SegSize int32
	// #total segments
	SegNum int32
	// retransmission timeout
	RetransmissionTimeout time.Duration
}

type Receiver interface {
	// # Blocking method
	Listen(conn Connection)
}

func NewDumpReceiver(conParam ConParam) Receiver {
	return &receiver{
		data: make([]byte, conParam.SegNum*conParam.SegSize),
		segSize: conParam.SegSize,
		totalSegNum: conParam.SegNum,
		windowSize: conParam.WinSize,
	}
}

type receiver struct {
	data               []byte
	segSize int32
	totalSegNum        int32
	windowSize         int32
}

func (r *receiver) Listen(conn Connection) {
	fmt.Printf("Receiver: starting listening\n")
	var lastReceivedSegIDX int32 = 0

	//Listen to segment channel for incomming msg
	for lastReceivedSegIDX < r.totalSegNum {

		var windowIDX int32 = 0
		for (windowIDX < r.windowSize) && (lastReceivedSegIDX < r.totalSegNum) {

			//Read pkg from connection
			segment := conn.Receive()

			windowIDX++
			lastReceivedSegIDX++

			fmt.Printf("Receiver: Received segment number %v, window idx %v\n", segment.SegNum, windowIDX)

			//Trascribe segment to destination
			startByteIDX := segment.SegNum * r.segSize
			for i := 0; i < len(segment.Data); i++ {
				r.data[startByteIDX+int32(i)] = segment.Data[i]
			}

			//Increment the windowSize
		}

		//Send acknowledge into ack connection segment with next expected segmentIDX
		ack := acknowledge.Acknowledge{
			SegIDX: lastReceivedSegIDX,
		}
		fmt.Printf("Receiver: sending acknowledge %v\n", ack.SegIDX)
		conn.SendAck(ack)
	}
	// complete transmission
	conn.SendAck(acknowledge.ClosingAck)

	fmt.Printf("Received str: \n")
	fmt.Printf("%v\n", string(r.data))

}
