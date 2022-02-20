package cfs

type NodeType uint16

const (
	NodeTypeFile NodeType = iota
	NodeTypeFolder
)
