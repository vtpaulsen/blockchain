package connection

import (
	"encoding/gob"
	"net"
	"sync"
)

type Connection struct {
	connections map[net.Conn]EncoderTuple
	mutex       sync.Mutex
}

type EncoderTuple struct {
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

func LockConnectionMutex(connection *Connection) {
	connection.mutex.Lock()
}

func UnlockConnectionMutex(connection *Connection) {
	connection.mutex.Unlock()
}

func GetConnections(connection *Connection) map[net.Conn]EncoderTuple {
	return connection.connections
}

func InitializeConnections(connection *Connection) {
	connection.connections = make(map[net.Conn]EncoderTuple)
}

func AddToConnections(connection *Connection, conn net.Conn) {
	LockConnectionMutex(connection)
	defer UnlockConnectionMutex(connection)

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	connection.connections[conn] = EncoderTuple{Encoder: enc, Decoder: dec}
}

func GetEncoder(connection *Connection, conn net.Conn) *gob.Encoder {
	return GetConnections(connection)[conn].Encoder
}

func GetDecoder(connection *Connection, conn net.Conn) *gob.Decoder {
	return GetConnections(connection)[conn].Decoder
}
