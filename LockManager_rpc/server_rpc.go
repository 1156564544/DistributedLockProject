package LockManager_rpc

type LockArgs struct {
	ClientId string
	LockName string
}

type LockReply struct {
	Error error
}

type LockManageArgs struct {
	Method   int // 0 for preempt; 1 for release
	ClientId string
	LockName string
}

type LockManageReply struct {
	Success bool
}

type LockModifyArgs struct {
	LockName string
	ClientId string
}

type LockModifyReply struct {
	Success bool
}

type ClientConnectArgs struct {
	ClientId string
}
type ClientConnectReply struct {
	Error error
}
