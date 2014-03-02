- Cluster is a library which provides the interfaces for creating and using the clusters which have multiple servers
- Using these server nodes we can send messages to peer server nodes
- The cluster configuration of peers have to specified in peers.json file

How to USE
	- Initialize cluster configuraton file peers.json
	- It contains details of all peer servers in the following format
		{
       			"P_id": 10,
       			"Host":"tcp://127.0.0.1:9006"
    		}

	- Ater initializing this file pass the file name and peer id to New(peer_id,filename) for creating server node.
	- There are 2 functions sender() and receiver() in cluster_test.go file which takes a server node as input and sends messages to peers
	- Other fuctions are well documented in code cluster.go file

Cluster.go
	- Go to the file cluster.go you can see how server nodes are created and messages are sent.
