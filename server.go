package main

import (
	"DistributedLockProject/LockManager_rpc"
	"sync"
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

	return nil
}

func (s *Server) UnLock(request LockManager_rpc.UnLockArgs, reply *LockManager_rpc.UnLockReply) error {

	return nil
}

func (s *Server) OwnTheLock(request LockManager_rpc.OwnTheLockArgs, reply *LockManager_rpc.OwnTheLockReply) error {

	return nil
}

func (s *Server) LockManage(request LockManager_rpc.LockManageArgs, reply *LockManager_rpc.LockManageReply) error {

	return nil
}

func (s *Server)
