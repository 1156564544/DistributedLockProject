package main

import (
	"DistributedLockProject/LockManager_rpc"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type Server struct {
	LeaderIp  string            // Ip of the leader
	Ip        string            // Ip of the current node
	ServersIp []string          // Ips of server nodes
	Clients   []string          // ClientIds of clients
	Locks     map[string]string // Locks in the distributed system

	mu sync.Mutex
}

func (s *Server) Lock(request LockManager_rpc.LockArgs, reply *LockManager_rpc.LockReply) error {
	lockKey := request.LockName
	clienyId := request.ClientId
	if s.Ip == s.LeaderIp {
		_, ok := s.Locks[lockKey]
		if ok == true {
			reply.Error = errors.New("The Lock already exists!")
			return nil
		} else {
			s.Locks[lockKey] = clienyId
			reply.Error = nil
			return nil
		}
	}
	return nil
}

func (s *Server) UnLock(request LockManager_rpc.LockArgs, reply *LockManager_rpc.LockReply) error {
	lockKey := request.LockName
	clienyId := request.ClientId
	if s.Ip == s.LeaderIp {
		id, ok := s.Locks[lockKey]
		if ok == false {
			reply.Error = errors.New("The Lock doesnot exists!")
			return nil
		}
		if id != clienyId {
			reply.Error = errors.New("You donot own the lock!")
			return nil
		}
		delete(s.Locks, lockKey)
		reply.Error = nil
		return nil
	}
	return nil
}

func (s *Server) OwnTheLock(request LockManager_rpc.LockArgs, reply *LockManager_rpc.LockReply) error {
	lockKey := request.LockName
	clienyId := request.ClientId
	if s.Ip == s.LeaderIp {
		id, ok := s.Locks[lockKey]
		if ok == false {
			reply.Error = errors.New("The Lock doesnot exists!")
			return nil
		}
		if id != clienyId {
			reply.Error = errors.New("You donot own the lock!")
			return nil
		}
		reply.Error = nil
		return nil
	}
	return nil
}

func (s *Server) LockManage(request LockManager_rpc.LockManageArgs, reply *LockManager_rpc.LockManageReply) error {

	return nil
}

func (s *Server) LockModify(request LockManager_rpc.LockModifyArgs, reply *LockManager_rpc.LockModifyReply) error {

	return nil
}

func (s *Server) ClientConnect(request LockManager_rpc.ClientConnectArgs, reply *LockManager_rpc.ClientConnectReply) error {
	clientId := request.ClientId
	s.Clients = append(s.Clients, clientId)
	fmt.Printf("Client %v connected to Server %v\n", clientId, s.Ip)
	reply.Error = nil
	return nil
}

func ConstructServer(ip string, leaderIp string, serversIp []string) *Server {
	s := &Server{}
	s.Ip = ip
	s.LeaderIp = leaderIp
	s.ServersIp = serversIp
	s.Locks = make(map[string]string)
	s.Clients = make([]string, 0)

	go func() {
		err := rpc.Register(s)
		if err != nil {
			panic(err.Error())
		}
		// HTTP注册
		rpc.HandleHTTP()
		// 端口监听
		listen, err := net.Listen("tcp", ip)
		if err != nil {
			panic(err.Error())
		}
		// 启动服务
		_ = http.Serve(listen, nil)
	}()

	return s
}

func main() {
	var ServerIP = []string{
		"0.0.0.0:80",
	}
	ip := "0.0.0.0:80"
	leaderIp := "0.0.0.0:80"

	_ = ConstructServer(ip, leaderIp, ServerIP)

	for {
		time.Sleep(10 * time.Second)
		fmt.Printf("Server is online...\n")
	}
}
