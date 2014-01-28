package cluster
import (
    "fmt"
    "os"
    "net"
    "io"
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
func New(id int,infile string) Server{
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
	var remote string
	//finding address of my pid
	for key:= range id_arr{
		if id_arr[key] == id{
			remote=adds_arr[key]
			break
		}
	}
	
	//listener
	go func(){
		var data = make([]byte, 1024)
		lis, error := net.Listen("tcp", remote)
		defer lis.Close()
		if error != nil { 
			fmt.Printf("Error creating listener: %s\n", error )
			os.Exit(1); 
		}
		for {
			var read = true
			con, error := lis.Accept()
			if error != nil { fmt.Printf("Error: Accepting data: %s\n", error); os.Exit(2); } 
			for read {
				n, error := con.Read(data); //reading
				switch error { 
				case io.EOF:
					read = false;
				case nil:
					var dat Envelope
    				json.Unmarshal(data[0:n],&dat) //converting into envelope
					mynode.Inbox()<-&dat		//put in inbox
				default:
					fmt.Printf("Error: Reading data : %s \n", error); 
					read = false;
				}
			}
			con.Close();
		}
	}()
	//sender
	go func(){
		var remote1 string
		tobesent:=<-mynode.Outbox() 
		sid:=tobesent.Pid
		tobesent.Pid=id
		b, _ := json.Marshal(tobesent)//convert the envelope to byte
		if sid==-1{		//if -1 then broadcast
			for key:= range id_arr{
				if id_arr[key]!=id{
					remote1=adds_arr[key]
					con, error := net.Dial("tcp",remote1);
					if error != nil { fmt.Printf("Host not found: %s\n", error ); os.Exit(1); }
		
					
					in, error := con.Write(b);
					if error != nil { fmt.Printf("Error sending data: %s, in: %d\n", error, in ); os.Exit(2); }				
				
					con.Close();
					
				}
			}
		}else{
		
			//finding address of my sid
			for key:= range id_arr{
				if id_arr[key] == sid{
					remote1=adds_arr[key]
					break
				}
			}
			con, error := net.Dial("tcp",remote1);
			if error != nil { fmt.Printf("Host not found: %s\n", error ); os.Exit(1); }
		
			in, error := con.Write(b);
			if error != nil { fmt.Printf("Error sending data: %s, in: %d\n", error, in ); os.Exit(2); }					
			con.Close();
		}		
	}()
	return mynode
}
