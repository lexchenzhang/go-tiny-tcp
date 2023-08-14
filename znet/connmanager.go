package znet

import (
	"errors"
	"fmt"
	"sync"
)

type IConnManager interface {
	Add(IConnection)
	Remove(IConnection)
	Get(uint32) (IConnection, error)
	Len() int
	ClearConn()
}

type ConnManager struct {
	connections map[uint32]IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]IConnection),
	}
}

func (c *ConnManager) Add(conn IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connManager add conn with ID = ", conn.GetConnID(), ", conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, conn.GetConnID())
	fmt.Println("connManager remove conn with ID = ", conn.GetConnID(), ", conn num = ", c.Len())
}

func (c *ConnManager) Get(connID uint32) (IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("can't find conn")
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for cid, conn := range c.connections {
		conn.Stop()
		delete(c.connections, cid)
	}
	c.connections = make(map[uint32]IConnection)
	fmt.Println("connManager clear all conn, conn num = ", c.Len())
}
