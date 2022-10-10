package log

import (
	"sync"
	"time"
)

const (
	StateRunnable = iota
	StateRunning
)

// Node 队列节点
type Node struct {
	task    func()        //任务
	expTime int64         //用来入队比较时间
	rotate  time.Duration //延时时间
	next    *Node
	prev    *Node
	t       *time.Timer
}

func NewNode(task func(), duration time.Duration) *Node {
	return &Node{
		task:    task,
		expTime: time.Now().Add(duration).UnixMilli(),
		rotate:  duration,
		t:       time.NewTimer(duration),
	}
}

// DelayRotateQueue 简易延时轮转队列
type DelayRotateQueue struct {
	head      *Node
	tail      *Node
	queueSize int //队列数量
	state     int //队列状态 0 可运行 1 运行中
	queueLock sync.Mutex
}

func NewDelayRotateQueue() *DelayRotateQueue {
	queue := &DelayRotateQueue{
		head:      nil,
		tail:      nil,
		queueSize: 0,
		state:     StateRunnable,
	}
	return queue
}

// Add 加入队列
func (queue *DelayRotateQueue) Add(node *Node) {
	if queue.state != StateRunnable {
		//已经开始运行不允许再添加新的节点
		return
	}
	queue.enq(node)
}

func (queue *DelayRotateQueue) enq(node *Node) {
	if node == nil || node.rotate < 0 {
		return
	}
	queue.queueLock.Lock()
	defer queue.queueLock.Unlock()
	//队列数量增加
	queue.queueSize++

	if queue.head == nil {
		queue.head = node
		queue.tail = node
		return
	}
	//根据时间入队
	for pNode := queue.tail; pNode != nil; pNode = pNode.prev {
		if pNode.expTime < node.expTime {
			if queue.tail == pNode {
				queue.tail = node
			} else {
				pNode.next.prev = node
				node.next = pNode.next
			}
			pNode.next = node
			node.prev = pNode
			return
		}
	}
	node.next = queue.head
	queue.head.prev = node
	queue.head = node
}

// Execute 执行队列延时任务
func (queue *DelayRotateQueue) Execute() {
	for {
		if queue.head == nil {
			//队列没有节点不允许执行
			queue.state = StateRunnable
			return
		}
		queue.state = StateRunning
		<-queue.head.t.C
		queue.queueLock.Lock()
		//重置节点
		executeNode := queue.head
		executeNode.expTime = time.Now().Add(executeNode.rotate).UnixMilli()
		executeNode.t.Reset(executeNode.rotate)
		//移除队首节点
		if executeNode.next == nil || executeNode.expTime <= executeNode.next.expTime {
			//只有一个节点或者重置后的节点依然能排在队首就不用出队了
			queue.queueLock.Unlock()
			executeNode.task()
			continue
		}
		queue.head = executeNode.next
		queue.head.prev = nil
		executeNode.next = nil
		queue.queueLock.Unlock()
		//重新归队
		queue.enq(executeNode)
		//执行队列任务
		executeNode.task()
	}
}

func (queue *DelayRotateQueue) Size() int {
	queue.queueLock.Lock()
	defer queue.queueLock.Unlock()
	return queue.queueSize
}
