package BpNet

import (
	"fmt"
)

func YModifyWeight(b_value float32, inpurParam interface{}) float32 {
	return b_value
}

func YModifyXValue(input float32, inpurParam interface{}) float32 {
	return input
}

func YDoInputFunc(input, value float32,) float32 {
	return input
}

type YNodeMgr struct {
	NodeMgr
	NodeRecvData map[string]float32
	NodeMgr_B *BNodeMgr
}

func (mgr *YNodeMgr) InitYNodeMgr(size int) {
	mgr.DoModifyConnWeight = YModifyWeight
	mgr.DoModifyValue = YModifyXValue
	mgr.DoFunc = YDoInputFunc

	var dataStruct []string
	for i := 0; i < size; i++ {
		name := fmt.Sprintf("Ynode_%d", i)
		dataStruct = append(dataStruct, name)
	}
	mgr.initNodes(dataStruct)
}

func (mgr *YNodeMgr) SetBNode(inMgr *BNodeMgr) {
	mgr.NodeMgr_B = inMgr
	for _, node := range mgr.nodeMap {
		mgr.NodeMgr_B.AddConnect(node)
	}
}

func (mgr *YNodeMgr) CheckResult() []float32{
	var retnList []float32
	for _, node := range mgr.nodeMap {
		retnList = append(retnList, node.DataRecv)  
	}
	return retnList
}

func (mgr *YNodeMgr) GetDelta() float32{
	var retn float32
	return retn
}
