package BpNet

import (
	"errors"
	// "fmt"
	"math/rand"
)

const initWeight float32 = 0.5

type DataSync func(float32)

type DoFuncByInput func(float32, float32) float32
type ModifyValue func(float32, interface{}) float32
type ModifyConnWeight func(float32, interface{}) float32

type nodeInfo struct {
	Weight  float32
	Channel DataSync
}

type Node struct {
	Name  string
	Value float32
	DataRecv float32

	DoModifyValue      ModifyValue
	DoModifyConnWeight ModifyConnWeight
	DoFunc             DoFuncByInput

	//private
	nodeMap  map[string]nodeInfo
	
}

func (node *Node) Init() {
	node.nodeMap = make(map[string]nodeInfo)
}

func (node *Node) Connect(name string, channel DataSync) error {

	if node.nodeMap == nil {
		return errors.New("node haven't init")
	}

	var newNodeInfo nodeInfo
	newNodeInfo.Channel = channel
	newNodeInfo.Weight = rand.Float32()

	node.nodeMap[name] = newNodeInfo
	return nil
}

func (node *Node) Clear(name string) {
	_, ok := node.nodeMap[name]
	if ok {
		delete(node.nodeMap, name)
	}
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

func (node *Node) RecvData(input float32) {
	// fmt.Println("out name:", node.Name, "RecvData:", input)
	node.DataRecv += input
}

func (node *Node) SendData() {
	for _, connNode := range node.nodeMap {
		connNode.Channel(node.DataRecv * connNode.Weight)
	}
	node.DataRecv = 0
}
