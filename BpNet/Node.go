package BpNet

import (
	"errors"
	"fmt"
)

const initWeight float32 = 0.5

type nodeInfo struct {
	Weight  float32
	Channel chan float32
}
type DoFuncByInput func(float32) float32
type ModifyValue func(float32, interface{}) float32
type ModifyConnWeight func(float32, string, interface{}) float32

type Node struct {
	Name    string
	Value   float32
	Channel chan float32

	DoModifyValue      ModifyValue
	DoModifyConnWeight ModifyConnWeight
	DoFunc DoFuncByInput

	//private
	nodeMap map[string]nodeInfo
}

func (node *Node) Init(chanSize int) {
	newChannel := make(chan float32, chanSize)
	node.Channel = newChannel
	node.nodeMap = make(map[string]nodeInfo)
}

func (node *Node) Connect(name string, nodeChan chan float32) error {

	if node.nodeMap == nil {
		return errors.New("node haven't init")
	}

	var newNodeInfo nodeInfo
	newNodeInfo.Channel = nodeChan
	newNodeInfo.Weight = initWeight

	node.nodeMap[name] = newNodeInfo
	fmt.Println("connect:", name)
	return nil
}

func (node *Node) GetWeigh(name string) float32 {
	var weight float32 = 0
	tmp, ok := node.nodeMap[name]
	if ok {
		weight = tmp.Weight
	}
	return weight * node.Value
}

func (node *Node) SetWeigh(name string, weightDelta float32) {
	tmp, ok := node.nodeMap[name]
	if ok {
		tmp.Weight = weightDelta
		node.nodeMap[name] = tmp
	}
}

func (node *Node) Update() {
	fmt.Println(node.Name,"start work")
	for {
		input := <- node.Channel
		fmt.Println(node.Name,"receive data:", input)
		node.Value = node.DoFunc(input)
	
		for _, nodeInfo := range node.nodeMap{
			channel := nodeInfo.Channel
			channel <- nodeInfo.Weight * node.Value
		}
	}
	
}
