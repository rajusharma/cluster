package main
import (
    "fmt"
	zmq "github.com/pebbe/zmq4"
    "os"
	"time"
    "encoding/json"
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
   peers []int
   adds []string
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
    for key :=range jsontype.Peers{
    	id_arr=append(id_arr,jsontype.Peers[key].P_id)
    	adds_arr=append(adds_arr,jsontype.Peers[key].Host)
    }
    
    
	mynode := Node{id,id_arr,adds_arr,make(chan *Envelope),make(chan *Envelope)}
	go SendMessage(mynode)
	go RecvMessage(mynode)
	return &mynode
}

func SendMessage(n Node) {
	var remote1 string	
	for {
		select {
			case tobesent, err := <- n.Outbox():
				if err == false {
					return
				} else { 
					sid:=tobesent.Pid
					tobesent.Pid=n.p_id
					b, _ := json.Marshal(tobesent)

					if sid==-1{		//if -1 then broadcast
						for key:= range n.peers{
							if n.peers[key]!=n.p_id{
								remote1=n.adds[key]
								
								client, _ := zmq.NewSocket(zmq.PUSH)
								defer client.Close()
								client.Connect(remote1)
								client.SendBytes(b,0)
								//client.Close()
								//println("hello")	
								//client.RecvBytes(0)								
							}
						}
					}else{		
						//finding address of my sid
						for key:= range n.peers{
							if n.peers[key] == n.p_id{
								remote1=n.adds[key]
								break
							}
						}
						
						client, _ := zmq.NewSocket(zmq.PUSH)
						defer client.Close()
						client.Connect(remote1)
						client.SendBytes(b,0)
						//client.RecvBytes(0)			
						
					}		
				}
		}
	}
	return
}
func RecvMessage(n Node) {
	var remote string
	//finding address of n's process id
	for key:= range n.peers{
		if n.peers[key] == n.p_id{
			remote=n.adds[key]
			break
		}
	}
	ser, _ := zmq.NewSocket(zmq.PULL)
	ser.Bind(remote)

	for {		
			msg, _ :=ser.RecvBytes(0)
			/*the code is not running after this recvbytes*/
			println("hello")
			var dat Envelope
			json.Unmarshal(msg,&dat) //converting into envelope
			n.Inbox()<-&dat		//put in inbox
			ser.SendBytes(msg,0)
		}
	return
}
	
func main() {
 
  //making servers and putting in array
   var ser_arr []Server
   p:=10
   for i:=1;i<6;i++{
   		ser_arr=append(ser_arr,New(p,"./peers.json"))
   		p++
   }
 
  //putting envelope in outbox 
   ser_arr[0].Outbox() <- &Envelope{Pid:-1, Msg: "helloaaaa there"}
  
   select {
       case envelope := <- ser_arr[1].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(5 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   select {
       case envelope := <- ser_arr[2].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(5 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   select {
       case envelope := <- ser_arr[3].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(5 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   select {
       case envelope := <- ser_arr[4].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(5 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   
}
