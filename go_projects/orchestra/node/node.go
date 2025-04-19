package node

//info: Defines the structure of the node
type Node struct {
	Name string //role: keeps track of the name 
	Ip string //role: keeps track of the address of the workers
	Cores int
	Memory int
	MemoryAllocated int
	Disk int 
	DiskAllocated int
	Role string
	TaskCount int //role: keeps track of the assigned tasks
}