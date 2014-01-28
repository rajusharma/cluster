package main
import (
    "cluster"
    "time"
    "fmt"
)

func main() {
 
   // parse argument flags and get this server's id into myid
   var ser_arr []cluster.Server
   p:=10
   for i:=1;i<6;i++{
   		ser_arr=append(ser_arr,cluster.New(p,"./test.json"))
   		p++
   }
 
   // the returned server object obeys the Server interface above.
  
   //wait for keystroke to start.
   // Let each server broadcast a message
   ser_arr[0].Outbox() <- &cluster.Envelope{Pid:11, Msg: "helloaaaa there"}
  
   select {
       case envelope := <- ser_arr[1].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(10 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   select {
       case envelope := <- ser_arr[2].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(10 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   select {
       case envelope := <- ser_arr[3].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(10 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   select {
       case envelope := <- ser_arr[4].Inbox(): 
           fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
  
       case <- time.After(10 * time.Second): 
           println("Waited and waited. Ab thak gaya\n")
   }
   
}
