package models

import (
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"github.com/fatih/set"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// Message 消息
type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //发送类型 1私聊 2群聊 3广播
	Media    int    //消息类型 1文字 2表情包 3图片 4音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clienMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// Chat 需要:发送者ID,接收者ID,消息类型,发送的内容,发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并检验token 等合法性
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//token := query.Get("token")
	//targetId := query.Get("targetId")
	//msgtype := query.Get("type")
	//context := query.Get("context")
	isvalida := true //checkToke() 待......
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3.用户关系
	//4.userid 跟node 绑定 并加锁
	rwLocker.Lock()
	clienMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go sendproc(node)
	//6.完成接收逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendproc >>>> msg:", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		broadMsg(data) //todo 将消息广播到局域网
		fmt.Println("[ws] recvProc <<<<<<", string(data))
	}
}

var udpsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine....")
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 1),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc data :", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc data :", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch data :", string(data))
		sendMsg(msg.TargetId, data)
		//case 2://群发
		//	sendGroupMsg()
		//case 3://广播
		//	sendAllMsg()
		//case 4:
		//
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg >>>>> userId :", userId, "  msg:", string(msg))
	rwLocker.RLock()
	node, ok := clienMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	community := Community{}

	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id= ? and target_id=? and type =2", userId,
		community.ID).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact = Contact{}
		contact.OwnerId = userId
		contact.TargetId = community.ID
		contact.Type = 2
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}
