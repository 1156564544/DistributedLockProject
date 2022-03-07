package main

import (
	"DistributedLockProject/LockManager_rpc"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Server struct {
	LeaderIp  string            // Ip of the leader
	Ip        string            // Ip of the current node
	ServersIp []string          // Ips of server nodes
	Clients   []string          // ClientIds of clients
	Locks     map[string]string // Locks in the distributed system
	Servers   []*rpc.Client
	Leader    *rpc.Client

	mu sync.Mutex
}

func (s *Server) Lock(request LockManager_rpc.LockArgs, reply *LockManager_rpc.LockReply) error {
	lockKey := request.LockName
	clientId := request.ClientId
	if s.Ip == s.LeaderIp {
		// Leader node
		_, ok := s.Locks[lockKey]
		if ok == true {
			fmt.Printf("Lock is already exists\n")
			reply.Success = false
			return nil
		} else {
			s.Locks[lockKey] = clientId
			reply.Success = true
			fmt.Printf("Leader lock %v from %v successfully!\n", lockKey, clientId)
			return nil
		}
	} else {
		// Follower node
		manageRequest := LockManager_rpc.LockManageArgs{}
		manageRequest.Method = 0
		manageRequest.ClientId = clientId
		manageRequest.LockName = lockKey
		manageReply := &LockManager_rpc.LockManageReply{}
		err := s.Leader.Call("Server.LockManage", manageRequest, manageReply)
		if err != nil || manageReply.Success == false {
			fmt.Printf("Follower %v request to leader failed!\n", s.Ip)
			reply.Success = false
			return nil
		}
		reply.Success = true
		fmt.Printf("Follower %v lock %v from %v successfully!\n", s.Ip, lockKey, clientId)
		return nil
	}
}

func (s *Server) UnLock(request LockManager_rpc.LockArgs, reply *LockManager_rpc.LockReply) error {
	lockKey := request.LockName
	clientId := request.ClientId
	if s.Ip == s.LeaderIp {
		// Leader node
		id, ok := s.Locks[lockKey]
		if ok == false {
			reply.Success = false
			return nil
		}
		if id != clientId {
			reply.Success = false
			return nil
		}
		delete(s.Locks, lockKey)
		fmt.Printf("Leader unlock %v from %v successfully!\n", lockKey, clientId)
		reply.Success = true
		return nil
	} else {
		// Follower node
		manageRequest := LockManager_rpc.LockManageArgs{}
		manageRequest.Method = 1
		manageRequest.ClientId = clientId
		manageRequest.LockName = lockKey
		manageReply := &LockManager_rpc.LockManageReply{}
		err := s.Leader.Call("Server.LockManage", manageRequest, manageReply)
		if err != nil || manageReply.Success == false {
			fmt.Printf("Follower %v request to leader failed!", s.Ip)
			reply.Success = false
			return nil
		}
		fmt.Printf("Follower %v unlock %v from %v successfully!\n", s.Ip, lockKey, clientId)
		reply.Success = true
		return nil
	}
}

func (s *Server) OwnTheLock(request LockManager_rpc.LockArgs, reply *LockManager_rpc.LockReply) error {
	lockKey := request.LockName
	clientId := request.ClientId
	id, ok := s.Locks[lockKey]
	if ok == false {
		reply.Success = false
		return nil
	}
	if id != clientId {
		reply.Success = false
		return nil
	}
	reply.Success = true
	return nil
}

func (s *Server) LockManage(request LockManager_rpc.LockManageArgs, reply *LockManager_rpc.LockManageReply) error {
	method := request.Method
	lockKey := request.LockName
	clientId := request.ClientId
	if method == 0 {
		// Preempt request
		_, ok := s.Locks[lockKey]
		if ok == true {
			// Lock is already
			fmt.Printf("Lock is already exists\n")
			reply.Success = false
			return nil
		} else {
			// Lock
			s.Locks[lockKey] = clientId
			for i := 0; i < len(s.ServersIp); i++ {
				if s.ServersIp[i] == s.LeaderIp {
					continue
				}
				Modifyrequest := LockManager_rpc.LockModifyArgs{}
				Modifyrequest.ClientId = clientId
				Modifyrequest.LockName = lockKey
				Modifyreply := &LockManager_rpc.LockModifyReply{}
				err := s.Servers[i].Call("Server.LockModify", Modifyrequest, Modifyreply)
				if err != nil || Modifyreply.Success == false {
					fmt.Printf("Modify follower %v failed!\n", s.ServersIp[i])
					return nil
				} else {
					fmt.Printf("Modify follower %v successfully!\n", s.ServersIp[i])
				}
			}
			reply.Success = true
			return nil
		}
	} else {
		// Release request
		id, ok := s.Locks[lockKey]
		if ok == false {
			// Lock is not exists!
			reply.Success = false
			return nil
		}
		if id == clientId {
			// UnLock
			for i := 0; i < len(s.ServersIp); i++ {
				if s.ServersIp[i] == s.LeaderIp {
					continue
				}
				Modifyrequest := LockManager_rpc.LockModifyArgs{}
				Modifyrequest.ClientId = ""
				Modifyrequest.LockName = lockKey
				Modifyreply := &LockManager_rpc.LockModifyReply{}
				err := s.Servers[i].Call("Server.LockModify", Modifyrequest, Modifyreply)
				if err != nil || Modifyreply.Success == false {
					fmt.Printf("Modify Client %v failed!\n", s.ServersIp[i])
					return nil
				} else {
					fmt.Printf("Modify follower %v successfully!\n", s.ServersIp[i])
				}
			}
			delete(s.Locks, lockKey)
			reply.Success = true
			return nil
		} else {
			// Lock is not owned by client with clientId
			fmt.Printf("Client is not owned by client %v!\n", clientId)
			reply.Success = false
			return nil
		}
	}
}

func (s *Server) LockModify(request LockManager_rpc.LockModifyArgs, reply *LockManager_rpc.LockModifyReply) error {
	lockKey := request.LockName
	clientId := request.ClientId
	if clientId == "" {
		fmt.Printf("Follower %v release lock %v.\n", s.Ip, lockKey)
		delete(s.Locks, lockKey)
	} else {
		fmt.Printf("Follower %v add lock %v of client %v.\n", s.Ip, lockKey, clientId)
		s.Locks[lockKey] = clientId
	}
	reply.Success = true
	return nil
}

func (s *Server) ClientConnect(request LockManager_rpc.ClientConnectArgs, reply *LockManager_rpc.ClientConnectReply) error {
	clientId := request.ClientId
	for i := 0; i < len(s.Clients); i++ {
		if s.Clients[i] == clientId {
			reply.Error = nil
			fmt.Printf("Client %v was already connected to Server %v\n", clientId, s.Ip)
			return nil
		}
	}
	s.Clients = append(s.Clients, clientId)
	fmt.Printf("Client %v connected to Server %v\n", clientId, s.Ip)
	//reply.Error = nil
	return nil
}

func (s *Server) FollowerConnect(request LockManager_rpc.FollowerConnectArgs, reply *LockManager_rpc.FollowerConnectReply) error {
	followerip := request.FollowerIp
	index := 0
	for ; index < len(s.ServersIp); index++ {
		if followerip == s.ServersIp[index] {
			break
		}
	}
	s.Servers[index], _ = rpc.DialHTTP("tcp", followerip)
	fmt.Printf("Client %v connect to leader.\n", followerip)
	reply.Success = true
	return nil
}

func ConstructServer(ip string, leaderIp string, serversIp []string) *Server {
	s := &Server{}
	s.Ip = ip
	s.LeaderIp = leaderIp
	s.ServersIp = serversIp
	s.Locks = make(map[string]string)
	s.Clients = make([]string, 0)
	s.Servers = make([]*rpc.Client, len(s.ServersIp))

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
	// Construct rpc client with leader node
	if ip != leaderIp {
		time.Sleep(2 * time.Second)
		s.Leader, _ = rpc.DialHTTP("tcp", leaderIp)
		request := LockManager_rpc.FollowerConnectArgs{}
		request.FollowerIp = ip
		reply := &LockManager_rpc.FollowerConnectReply{}
		err := s.Leader.Call("Server.FollowerConnect", request, reply)
		if err != nil || reply.Success == false {
			fmt.Printf("Client %v connect to leader node failed!\n", ip)
			return nil
		}
	}

	return s
}

func main() {
	var ServerIP = []string{
		"0.0.0.0:80",
		"0.0.0.0:81",
		"0.0.0.0:82",
	}
	leaderIp := "0.0.0.0:80"
	ip := os.Args[1]

	_ = ConstructServer(ip, leaderIp, ServerIP)

	for {
		time.Sleep(10 * time.Second)
		fmt.Printf("Server is online...\n")

	}
}
