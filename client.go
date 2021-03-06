package main

import (
	"DistributedLockProject/LockManager_rpc"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

var ServerIP = []string{
	"0.0.0.0:80",
	"0.0.0.0:81",
	"0.0.0.0:82",
}

type DLClient struct {
	clientId    string
	isConnected bool
	fd          int
	client      *rpc.Client
}

func (c *DLClient) getClientId() string {
	return c.clientId
}

func (c *DLClient) connectedToServer(serverAddr string) bool {
	c.client, _ = rpc.DialHTTP("tcp", serverAddr)
	request := LockManager_rpc.ClientConnectArgs{}
	request.ClientId = c.clientId
	reply := &LockManager_rpc.ClientConnectReply{}
	err := c.client.Call("Server.ClientConnect", request, reply)
	if err != nil {
		return false
	}
	if reply.Error != nil {
		fmt.Println(reply.Error.Error())
		return false
	}
	c.isConnected = true
	fmt.Printf("Client %v connected to server %v \n", c.clientId, serverAddr)
	return true
}

func (c *DLClient) IsConnected() bool {
	return c.isConnected
}

func (c *DLClient) TryLock(LockKey string) bool {
	request := LockManager_rpc.LockArgs{}
	reply := &LockManager_rpc.LockReply{}
	request.LockName = LockKey
	request.ClientId = c.getClientId()
	err := c.client.Call("Server.Lock", request, reply)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("Client %v locks %v failed!\n", c.clientId, LockKey)
		return false
	}
	if reply.Success == false {
		fmt.Printf("Client %v locks %v failed!\n", c.clientId, LockKey)
		return false
	}
	fmt.Printf("Client %v locks %v successfully!\n", c.clientId, LockKey)
	return true
}

func (c *DLClient) TryUnLock(LockKey string) bool {
	request := LockManager_rpc.LockArgs{}
	reply := &LockManager_rpc.LockReply{}
	request.LockName = LockKey
	request.ClientId = c.getClientId()
	err := c.client.Call("Server.UnLock", request, reply)
	if err != nil {
		return false
	}
	if reply.Success == false {
		fmt.Printf("Client %v try to unlocks %v failed!\n", c.clientId, LockKey)
		return false
	}

	return true
}

func (c *DLClient) OwnTheLock(LockKey string) bool {
	request := LockManager_rpc.LockArgs{}
	reply := &LockManager_rpc.LockReply{}
	request.LockName = LockKey
	request.ClientId = c.getClientId()
	err := c.client.Call("Server.OwnTheLock", request, reply)
	if err != nil {
		return false
	}
	if reply.Success == false {
		return false
	}

	return true
}

func getFD(serverAddr string) int {
	i := 0
	for ; i < len(ServerIP); i++ {
		if ServerIP[i] == serverAddr {
			return i
		}
	}

	return -1
}

func DistributedLock(serverAddr string, clientId string) DLClient {
	c := DLClient{}
	c.clientId = clientId
	c.isConnected = false
	c.fd = getFD(serverAddr)
	if c.fd == -1 {
		log.Fatalf("ServerAddr is false: %v", serverAddr)
	}

	c.connectedToServer(serverAddr)

	return c
}

func main() {
	serverAddr := os.Args[1]
	clientId := os.Args[2]
	DL := DistributedLock(serverAddr, clientId)

	//DL.TryLock("lock2")
	if DL.OwnTheLock("lock2") {
		fmt.Println("Client own the lock")
	} else {
		fmt.Println("Client not own the lock")
	}
	//DL.TryUnLock("lock1")
	//if DL.OwnTheLock("lock1") {
	//	fmt.Println("Client own the lock")
	//} else {
	//	fmt.Println("Client not own the lock")
	//}

}
