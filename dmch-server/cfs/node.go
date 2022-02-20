package cfs

type INode struct {
	Name string

	Type    NodeType
	Subtype NodeSubtype

	Children []INode
}
