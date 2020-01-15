package BpNet

import (
	"fmt"
)

type NodeMgr struct {
	nodeMap            map[string] *Node
	DoModifyValue      ModifyValue
	DoModifyConnWeight ModifyConnWeight
	DoFunc             DoFuncByInput
}

func (mgr *NodeMgr) initNodes(dataStruct []string) {
	mgr.nodeMap = make(map[string] *Node)

	for _, value := range dataStruct {
		node := new(Node)
		node.Init()
		node.Name = value
		node.DoModifyValue = mgr.DoModifyValue
		node.DoModifyConnWeight = mgr.DoModifyConnWeight
		node.DoFunc = mgr.DoFunc
		mgr.nodeMap[value] = node
	}
}

//
func (mgr *NodeMgr) AddConnect(connNode *Node) {
	for _, node := range mgr.nodeMap {
		err := node.Connect(connNode.Name, connNode.RecvData)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func (mgr *NodeMgr) DelConnect(connNode *Node) {
	for _, node := range mgr.nodeMap {
		node.Clear(connNode.Name)
	}
}

func (mgr *NodeMgr) DoModify(connNode Node, nodeName string, valueParams interface{}, weightParams interface{}) {
	for _, node := range mgr.nodeMap {
		node.Value = node.DoModifyValue(node.Value, valueParams)
		weight := node.DoModifyConnWeight(node.Value, weightParams)
		node.SetWeigh(nodeName, weight)
	}
}

func (mgr *NodeMgr) DoSend() {
	for _, node := range mgr.nodeMap {
		node.SendData()
	}
}