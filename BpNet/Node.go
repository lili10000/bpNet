package BpNet

import (
	// "errors"
	// "fmt"
	"math/rand"
)

const initWeight float64 = 1

type DataSync func(float64)

type DoFuncByInput func(float64) float64
type ModifyValue func(float64, interface{}) float64
type ModifyConnWeight func(float64, interface{}) float64

type nodeInfo struct {
	Weight  float64
	Channel DataSync
}

type Node struct {
	Name  string
	Value float64
	DataRecv float64

	DoModifyValue      ModifyValue
	DoModifyConnWeight ModifyConnWeight
	DoFunc             DoFuncByInput
	//private
	NodeList  []nodeInfo

	Weight float64
	
}

func (node *Node) Init() {
	node.Value = 0
}

func (node *Node) Connect(index int, channel DataSync) error {

	var newNodeInfo nodeInfo
	newNodeInfo.Channel = channel
	newNodeInfo.Weight = node.Weight + rand.Float64()
	// newNodeInfo.Weight = 0.1

	node.NodeList = append(node.NodeList, newNodeInfo) 
	return nil
}


func (node *Node) GetWeigh(index int) float64 {
	var weight float64 = 0
	tmp := node.NodeList[index]
	weight = tmp.Weight
	return weight * node.Value
}

func (node *Node) SetWeigh(index int, weightDelta float64) {
	tmp:= node.NodeList[index]
	tmp.Weight = weightDelta
	node.NodeList[index] = tmp
}

func (node *Node) RecvData(input float64) {
	// fmt.Println("in name:", node.Name, "Data:", input)
	node.DataRecv += input

}

func (node *Node) SendData() {
	// fmt.Println("name:", node.Name, "ready Data:", node.DataRecv)
	sendData := node.DoFunc(node.DataRecv - node.Value)
	// fmt.Println("name:", node.Name, "cov Data:", sendData)

	for _, connNode := range node.NodeList {
		// fmt.Println("out name:", node.Name, "Data:", sendData*connNode.Weight, "weight:", connNode.Weight)
		connNode.Channel(sendData * connNode.Weight)
	}
	// node.DataRecv = 0
}
