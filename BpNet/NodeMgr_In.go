package BpNet

import (
	"fmt"
)

type InNodeMgr struct {
	NodeMgr
}

type Y_Param struct {
	Real    float32
	Get     float32
}

type InWeightParams struct {
	StepLen float32
	B       float32
	Y_param []Y_Param
}

func InModifyWeight(x_value float32, inpurParam interface{}) float32 {
	params := inpurParam.(InWeightParams)
	var Y_deltaSum float32 = 0
	for _, value := range params.Y_param {
		g := value.Get * (1 - value.Get) * (value.Real - value.Get)
		Y_deltaSum += g * params.B
	}
	E := params.B * (1 - params.B) * Y_deltaSum
	retn := params.StepLen * E * x_value
	return retn
}

func InModifyXValue(input float32, inpurParam interface{}) float32 {
	return input
}

func InDoInputFunc(input, inValue float32) float32 {
	return input
}

func (mgr *InNodeMgr) InitInNodeMgr(dataStruct []string) {
	mgr.DoModifyConnWeight = InModifyWeight
	mgr.DoModifyValue = InModifyXValue
	mgr.DoFunc = InDoInputFunc
	mgr.initNodes(dataStruct)
}

func (mgr *InNodeMgr) DoInput(dataMap map[string]float32) {
	for name, value := range dataMap {
		node, ok := mgr.nodeMap[name]
		if !ok {
			fmt.Println(name, "not in map")
			return
		}
		node.RecvData(value)
	}
}
