package BpNet

import (
	"fmt"
)
type GetRealIn func(float64) float64

type NodeMgr struct {
	nodeList           []*Node
	DoModifyValue      ModifyValue
	DoModifyConnWeight ModifyConnWeight
	DoFunc             DoFuncByInput
	mgrName string
	StepLen float64
	DoGetRealIn GetRealIn
}

func (mgr *NodeMgr) initNodes( size int) {
	mgr.StepLen = 0.3

	for i:=0; i < size; i++ {
		node := new(Node)
		node.Init()
		node.Name = fmt.Sprintf("%s_%d", mgr.mgrName, i)
		node.DoModifyValue = mgr.DoModifyValue
		node.DoModifyConnWeight = mgr.DoModifyConnWeight
		node.DoFunc = mgr.DoFunc
		node.Weight = float64(i)*0.1
		mgr.nodeList = append(mgr.nodeList, node)
	}
}

//
func (mgr *NodeMgr) AddConnect(connNode *Node) {
	for index, node := range mgr.nodeList {
		err := node.Connect(index, connNode.RecvData)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}


func (mgr *NodeMgr) DoModify(connNode Node, nodeName string, valueParams interface{}, weightParams interface{}) {
	for index, node := range mgr.nodeList {
		node.Value = node.DoModifyValue(node.Value, valueParams)
		weight := node.DoModifyConnWeight(node.Value, weightParams)
		node.SetWeigh(index, weight)
	}
}

func (mgr *NodeMgr) DoSend() {
	for _, node := range mgr.nodeList {
		node.SendData()
	}
}
