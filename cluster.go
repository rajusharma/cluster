package cluster
import (
    "fmt"
	zmq "github.com/pebbe/zmq4"
    "os"
	//"time"
    "encoding/json"
    "encoding/gob"
    "bytes"
    "io/ioutil"
)
const (BROADCAST = -1)
type Envelope struct {
    // On the sender side, Pid identifies the receiving peer. If instead, Pid is
    // set to cluster.BROADCAST, the message is sent to all peers. On the receiver side, the
    // Id is always set to the original sender. If the Id is not found, the message is silently dropped
    Pid int
    //message string
   
    // An id that globally and uniquely identifies the message, meant for duplicate detection at
    // higher levels. It is opaque to this package.
    MsgId int64
   
    // the actual message.
    Msg interface{}
}
   
type Server interface {
    // Id of this server
    Pid() int
   
    // array of other servers' ids in the same cluster
    Peers() []int

    // the channel to use to send messages to other peers
    // Note that there are no guarantees of message delivery, and messages
    // are silently dropped
    Outbox() chan *Envelope
   
    // the channel to receive messages from other peers.
    Inbox() chan *Envelope
}

type Node struct{
   p_id int
   //my_add string
   peers []int
   adds []string
   sendids []int		//ids of peer whom to send
   rec_soc *zmq.Socket //socket for recieving messages
   send_soc []*zmq.Socket //socket for sending messages
   chout chan *Envelope
   chin chan *Envelope
}

func (n Node) Pid() int {
   return n.p_id
} 

func (n Node) Peers() []int {
   return n.peers
} 

func (n Node) Outbox() chan *Envelope{
	return n.chout
}

func (n Node) Inbox() chan *Envelope{
	return n.chin
}

//json structs
type jsonobject struct {
    Peers []points
}
 
type points struct {
    P_id int
    Host   string
}
//////////

//constructs a new server node 
func New(id int,infile string) *Node{
	//reading json file and storing in arrays
	file, e := ioutil.ReadFile(infile)
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    } 
    var jsontype jsonobject
    json.Unmarshal(file,&jsontype)
    
    //parsing into array from json object
    var id_arr []int
    var adds_arr []string
	var remote string
    for key :=range jsontype.Peers{
		id1 := jsontype.Peers[key].P_id
		addr := jsontype.Peers[key].Host
		//check the address of current pid to open a socket for recieving
		if id1==id {
			remote=addr
		}
    	id_arr=append(id_arr,id1)
    	adds_arr=append(adds_arr,addr)
    }
    //creating a socket for recieving
	client, err := zmq.NewSocket(zmq.PULL)
	if err!=nil{
		println(err)
	}
	client.Bind(remote)
	
	//creating sockets for sending
	var send_soc []*zmq.Socket
	var sendids []int
	for key:=range id_arr{
		if id_arr[key]!=id {
			sender, err := zmq.NewSocket(zmq.PUSH)
			if err!=nil{
				println(err)
			}
			sender.Connect(adds_arr[key])
			send_soc = append(send_soc,sender)
			sendids = append(sendids,id_arr[key])
		}
	}
	mynode := Node{id,id_arr,adds_arr,sendids,client,send_soc,make(chan *Envelope),make(chan *Envelope)}
	go SendMessage(mynode)
	go RecvMessage(mynode)
	return &mynode
}

func SendMessage(n Node) {
//println("hello")
	//var remote1 string	
	for {
		select {
			case tobesent, err := <- n.Outbox():
				if err == false {
					return
				} else { 
					sid:=tobesent.Pid
					tobesent.Pid=n.p_id
					//encoding the envelope
					bs :=new(bytes.Buffer) 
					enc := gob.NewEncoder(bs) 
					err := enc.Encode(tobesent)
					if err != nil {
						println("encode error:", err)
					}      
					if sid==-1{		//if -1 then broadcast
						for key:= range n.send_soc{
								//send the data
								_,err1 :=n.send_soc[key].SendBytes(bs.Bytes(), 0)
								if err1!=nil{
									println(err1)
								}													
							}
					} else {		
						//finding address of my sid and send the msg
						for key:= range n.sendids{
							if n.sendids[key] == sid{
								//send the data to destination server
								_,err1 :=n.send_soc[key].SendBytes(bs.Bytes(), 0)
								if err1!=nil{
									println(err1)
								}
								break
							}
						}			
					}		
				}
		}
	}
	return
}
func RecvMessage(n Node) {
	for {		
			msg, _ :=n.rec_soc.RecvBytes(0)
			//decoding the envelope
			b:=bytes.NewBuffer(msg)
			dec := gob.NewDecoder(b) 
			var dat Envelope
			err := dec.Decode(&dat)
			if err != nil {
				println("decode error:", err)
			}
			//json.Unmarshal(msg,&dat) //converting into envelope
			n.Inbox()<-&dat		//put in inbox
		}
	return
}
	
