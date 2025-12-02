package websocket

import (
	"bytes"
	"encoding/binary"
)

// EVENT CODES
const (
	EventShoot  = 1 // người bắn
	EventHit    = 2 // hành động bắn
	EventMove   = 3 // hành đọng di chuyển
	EventDead   = 4 // chết
	EventWinner = 5 // nguồi chiến thắng
)

// Packet hit: 1 byte event + 6 byte số
// [1][shooter][target][damage]
type HitPacket struct {
	Event   uint8
	Shooter uint16
	Target  uint16
	Damage  uint16
}

func EncodeHit(shooter, target, damage uint16) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(EventHit))
	binary.Write(buf, binary.BigEndian, shooter)
	binary.Write(buf, binary.BigEndian, target)
	binary.Write(buf, binary.BigEndian, damage)
	return buf.Bytes()
}
