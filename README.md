# Welcome to GoSAN!

GoSAN is the discrete-event based simulator. 

# What it simulates?

 - Resource usage
	 - CPU
	 - Links
	 - Storage
		 - SSD
		 - HDD
 - Load generation
	 - Read mode
	 - Write mode

# Installation
Execute the following command: 

    go get -u gitlab.com/lambda-hse/tatlin-hse/gosan

# Example

To demonstrate how GoSAN can be used we implemented a toy version of actual storage array network. It consists of:

 - Load generator
 - 2 controllers 
 - 8 storage units:
	 - 4 SSD
	 - 4 HDD

Client sends data blocks to controllers. They, in own turn, write data to storage units in the round-robin algorithm.  

 
![enter image description here](https://sun9-64.userapi.com/c854216/v854216025/101217/U2KF-i4FyVc.jpg)


