package cluster
import (
	"fmt"
	"time"
	"testing"
)
var totalsent int = 0
var totalreceived int = 0
var numberofSer int = 5
var corruptmsg int = 0
var failedmsg int = 0
func Test_cluster(t *testing.T) {
	//making servers and putting in array
   var ser_arr []Server
   //pid starts from 10-14
   p:=10
   for i:=1;i<6;i++{
   		ser_arr=append(ser_arr,New(p,"./peers.json"))
   		p++
   }
   for i:=0;i<5;i++{
   		for j:=0;j<5;j++{ 		
   			go sender(ser_arr[j])			
   			go receiver(ser_arr[j])
   			time.Sleep(1*time.Second)
   		}  		
   }
   println("TotalSent = ",totalsent)
   println("TotalReceived = ",totalreceived)
   println("TotalFailed = ",failedmsg)
  
   return
   
}

//this function sends messages to a servernode
func sender(s Server){
	msg1:="hello"
	//broadcast a message
	s.Outbox()<-&Envelope{Pid:-1, Msg: msg1}
	
	//add 4 to totalsent because it cant send to itself
	totalsent = totalsent+4
	//println(totalsent)
	//sending message to all peers 10-14 
	
	//my id
	sid:=s.Pid();
	for i:=10;i<15;i++{
		if i!=sid{
		s.Outbox()<-&Envelope{Pid:i, Msg: msg1}
		totalsent = totalsent+1
		}
	}
	return
}

//this function receives messages on a servernode
func receiver(s Server){
	for{
	select {
       case envelope := <- s.Inbox(): 
           //fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
           if envelope.Msg == "hello"{
           		totalreceived = totalreceived+1
           }else{
           		corruptmsg = corruptmsg+1
           }
  
       case <- time.After(5 * time.Second): 
       		failedmsg=failedmsg+1
           //println("Waited and waited. Ab thak gaya\n")
   	}
 	}
}
