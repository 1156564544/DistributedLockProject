package LockManager_rpc

type LockArgs struct {
	ClientId string
	LockName string
}

type LockReply struct {
	Error error
}

type UnLockArgs struct {
	ClientId string
	LockName string
}

type UnLockReply struct {
	Error error
}

type OwnTheLockArgs struct {
	ClientId string
	LockName string
}

type OwnTheLockReply struct {
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
