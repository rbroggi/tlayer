package acknowledge

type Acknowledge struct {
	SegIDX int32
}

var ClosingAck = Acknowledge{
	SegIDX: -1,
}
