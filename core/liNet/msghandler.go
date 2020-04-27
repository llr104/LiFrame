package liNet

import (
	"container/list"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/utils"
	"reflect"
	"strings"
)

type MsgHandle struct {
	Apis           map[string] *list.List //存放每个msgName 所对应的处理方法的map属性
	WorkerPoolSize uint32                 //业务工作Worker池的数量
	TaskQueue      []chan liFace.IRequest //Worker负责取任务的消息队列
	TaskExit       []chan bool
}

func NewMsgHandle(workerSize uint32) *MsgHandle {
	return &MsgHandle{
		Apis: make(map[string]*list.List),
		WorkerPoolSize:workerSize,
		//一个worker对应一个queue
		TaskQueue:make([]chan liFace.IRequest, utils.GlobalObject.ServerWorkerSize),
		TaskExit:make([]chan bool, utils.GlobalObject.ServerWorkerSize),
	}
}

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request liFace.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgName(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}


//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request liFace.IRequest, respond liFace.IMessage) {

	//执行对应处理方法
	rpcType := request.GetMessage().GetType()
	if rpcType == liFace.RpcAck {
		m := request.GetMessage()
		request.GetConnection().CheckRpc(request.GetMessage().GetSeq(), m)
	}else{
		msgName := request.GetMessage().GetMsgName()
		arr := strings.Split(msgName,".")
		isFound := false

		if len(arr) == 2{
			//绝对匹配
			nameSpace := arr[0]
			funcName := arr[1]

			if l, ok := mh.Apis[nameSpace]; ok {
				for handler := l.Front(); nil != handler; handler = handler.Next() {
					//t := reflect.TypeOf(handler)
					//fmt.Println("DoMsgHandler reflect type",t)
					v := reflect.ValueOf(handler.Value)
					//fmt.Println("DoMsgHandler reflect value",v)
					method := v.MethodByName(funcName)
					if method.IsValid() == false {
						utils.Log.Warn("DoMsgHandler warning %s successFun not found",funcName)
					}else{
						isFound = true
						router := handler.Value.(liFace.IRouter)
						ret := router.PreHandle(request, respond)
						if ret {
							in := make([]reflect.Value, 2)
							in[0] = reflect.ValueOf(request)
							in[1] = reflect.ValueOf(respond)
							method.Call(in)
							router.PostHandle(request, respond)
							//回复client
							request.GetConnection().RpcReply(request.GetMessage().GetMsgName(), request.GetMessage().GetSeq(), respond.GetBody())

						}else{
							utils.Log.Warn("DoMsgHandler skip: %s",msgName)
							//回复client
							request.GetConnection().RpcReply(request.GetMessage().GetMsgName(), request.GetMessage().GetSeq(), respond.GetBody())
						}
					}
				}
			}
		}

		//查看是否有通配匹配
		if l, ok := mh.Apis["*.*"]; ok{
			isFound = true
			respond := Message{}
			respond.SetSeq(request.GetMessage().GetSeq())

			for handler := l.Front(); nil != handler; handler = handler.Next() {
				router := handler.Value.(liFace.IRouter)
				router.EveryThingHandle(request, &respond)
				//回复client
				request.GetConnection().RpcReply(request.GetMessage().GetMsgName(), request.GetMessage().GetSeq(), respond.GetBody())
			}
		}

		if isFound == false{
			utils.Log.Warning("DoMsgHandler not found: %s handler",msgName)
		}
	}

}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(router liFace.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	name := router.NameSpace()
	if _, ok := mh.Apis[name]; ok {
		panic("repeated api, NameSpace = " + name)
	}

	//2 添加msg与api的绑定关系
	l, ok := mh.Apis[name]
	if ok == false {
		 l = list.New()
		 mh.Apis[name] = l
	}
	l.PushBack(router)
	mh.Apis[name] = l

	utils.Log.Info("Add api NameSpace = %s", name)
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan liFace.IRequest,taskExit chan bool) {
	utils.Log.Info("Worker ID = %d is started.", workerID)
	//不断的等待队列中的消息
	for {
		select {
			//有消息则取出队列的Request，并执行绑定的业务方法
			case request := <-taskQueue:
				rsp := Message{}
				mh.DoMsgHandler(request, &rsp)
			case isExit := <-taskExit:
				if isExit {
					utils.Log.Info("Worker ID = %d is stop.", workerID)
					return
				}
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i:= 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan liFace.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		mh.TaskExit[i] = make(chan bool, 1)

		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i], mh.TaskExit[i])
	}
}

func (mh *MsgHandle) StopWorkerPool() {
	for i:= 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskExit[i] <- true
		close(mh.TaskExit[i])
	}
}