package BpNet

import (
	"fmt"
)

type BNodeMgr struct {
	NodeMgr
	NodeMgr_In *InNodeMgr
}

type BWeightParams struct {
	StepLen float32
	YParam  Y_Param
}

func BModifyWeight(b_value float32, inpurParam interface{}) float32 {
	params := inpurParam.(BWeightParams)
	value := params.YParam
	g := value.Get * (1 - value.Get) * (value.Real - value.Get)
	retn := b_value + params.StepLen*g*b_value
	return retn
}

func BModifyXValue(input float32, inpurParam interface{}) float32 {
	return input
}

func BDoInputFunc(input , value float32) float32 {
	return input
}

func (mgr *BNodeMgr) InitBNodeMgr(size int) {
	mgr.DoModifyConnWeight = BModifyWeight
	mgr.DoModifyValue = BModifyXValue
	mgr.DoFunc = BDoInputFunc

	var dataStruct []string
	for i := 0; i < size; i++ {
		name := fmt.Sprintf("Bnode_%d", i)
		dataStruct = append(dataStruct, name)
	}
	mgr.initNodes(dataStruct)
}

func (mgr *BNodeMgr) SetInNode(inMgr *InNodeMgr) {
	mgr.NodeMgr_In = inMgr
	for _, node := range mgr.nodeMap {
		mgr.NodeMgr_In.AddConnect(node)
	}
}

