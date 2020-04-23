package liNet

import (
	"errors"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/utils"
	"sync"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint32]liFace.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

/*
	创建一个链接管理
 */
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections:make(map[uint32]liFace.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn liFace.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn

	utils.Log.Info("connection add to ConnManager successfully: conn num = %d", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn liFace.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	utils.Log.Info("connection Remove ConnID = %d, successfully: conn num = %d", conn.GetConnID(), connMgr.Len())
}

//利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (liFace.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) BroadcastMsg(msgName string, data []byte){
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for _, conn := range connMgr.connections {
		conn.RpcCall(msgName, data, nil)
	}
}

//清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	needClear := false
	//停止并删除全部的连接信息
	for _, conn := range connMgr.connections {
		//停止
		conn.Stop()
		needClear = true
	}

	if needClear{
		connMgr.connections = make(map[uint32]liFace.IConnection)
	}

	utils.Log.Info("Clear All Connections successfully: conn num = %d", connMgr.Len())
}
