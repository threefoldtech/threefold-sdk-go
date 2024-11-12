package db

type Node struct {
	NodeId uint32 `gorm:"uniqueIndex:uni_gpu_node_twin_id"`
}
