package main

import "log"

var ServerIP = []string{
	"0.0.0.0:80",
}

type Client struct {
	clientId    string
	isConnected bool
	fd          int
}

func (c *Client) getClientId() string {
	return c.clientId
}

func (c *Client) connectedToServer(serverAddr string) bool {

}

func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) TryLock(LockKey string) bool {

}

func (c *Client) TryUnLock(LockKey string) bool {

}

func (c *Client) OwnTheLock(LockKey string) bool {

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

func DistributedLock(serverAddr string, clientId string) Client {
	c := Client{}
	c.clientId = clientId
	c.isConnected = false
	c.fd = getFD(serverAddr)
	if c.fd == -1 {
		log.Fatalf("ServerAddr is false: %v", serverAddr)
	}

	c.connectedToServer(serverAddr)
}
