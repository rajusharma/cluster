cluster
=======

COMMANDS TO RUN:







IMPLEMENTATION:

/*struct which contains process id and message*/
type Envelope struct {
}

/*interface for the server*/
type Server interface {
}

/*struct of a process which sends and recieves messages*/
type Node struct{
   p_id int
   peers []int //array of process id's
   adds []string  //array of process address's
   chout chan *Envelope //channesl for output
   chin chan *Envelope	//channel for input
   
}

/*New() function initializes a new process Node.*/
/*There are 2 Go routines one is for listening and one for sending messages*/
/*First read the json file which contains all the peer process id's and address's then store it to an array and pass these arrays to Node() which creates and initializes process*/
func New(id int,infile string) Server{
	
   /*code for reading the json file and storing it in array*/
   
	mynode := Node{id,id_arr,adds_arr,make(chan *Envelope),make(chan *Envelope)}
	var remote string
	/*finding address of current process which is my process by looping in array*/ 
	for key:= range id_arr{
		if id_arr[key] == id{
			remote=adds_arr[key]
			break
		}
	}

	/*Listener accept the connection and convert the byte msg to envelope and stores it in inbox() channel*/
	go func(){
	
	
		/*CODE*/
		
		
	}()
	/* Sender converts the envelope to byte using marshal command and then changes the pid of envelope and sends to the destination address, if the address is -1 then it sends to all process's by looping*/
	
	go func(){
	
	/*CODE*/
	
	
	}
	return mynode
}

