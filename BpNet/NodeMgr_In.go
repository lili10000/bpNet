package BpNet

import (
	"fmt"
)

type InNodeMgr struct {
	NodeMgr
}

type Y_Param struct {
	Real float64
	Get  float64
}

type InWeightParams struct {
	StepLen float64
	B       float64
	GetWeight_B []float64
	Y_param []Y_Param
}

func InModifyWeight(x_value float64, inpurParam interface{}) float64 {
	params := inpurParam.(InWeightParams)
	var Y_deltaSum float64 = 0
	for index, value := range params.Y_param {
		g := value.Get * (1 - value.Get) * (value.Real - value.Get)
		Y_deltaSum += g * params.GetWeight_B[index]
	}
	E := params.B * (1 - params.B) * Y_deltaSum
	retn := params.StepLen * E * x_value
	return retn
}

func InModifyXValue(input float64, inpurParam interface{}) float64 {
	return input
}

func InDoInputFunc(input float64) float64 {
	return input
}
func InDoRealIn (realout float64) float64{
	return realout
}

func (mgr *InNodeMgr) InitInNodeMgr(size int) {
	mgr.DoModifyConnWeight = InModifyWeight
	mgr.DoModifyValue = InModifyXValue
	mgr.DoFunc = InDoInputFunc
	mgr.DoGetRealIn = InDoRealIn
	mgr.mgrName = "InNode"
	mgr.initNodes(size)
}

func (mgr *InNodeMgr) DoInput(dataList []float64) {
	if len(dataList) != len(mgr.nodeList){
		fmt.Println("DoInput: input dataList size: ",len(dataList), " != nodeList size:", len(mgr.nodeList))
		return
	}

	for index, value := range dataList {
		node := mgr.nodeList[index]
		node.RecvData(value)
	}
}

func (mgr *InNodeMgr) DoDelta_In(Y_real []float64, Y_get []float64, B_Get []float64,B_weight[][]float64){
	var tmpList []Y_Param
	for index, _ := range Y_real{
		var tmp Y_Param 
		tmp.Get = Y_get[index]
		tmp.Real = Y_real[index]
		tmpList = append(tmpList, tmp)
	}
	
	
	for _, InNode := range mgr.nodeList {
		var param InWeightParams
		param.StepLen = mgr.StepLen
		for index, _ := range InNode.NodeList {
			param.Y_param = tmpList
			param.B = B_Get[index]
			param.GetWeight_B = B_weight[index]
			InNode.NodeList[index].Weight += InModifyWeight(InNode.Value, param)
		}
	}
}

func (mgr *InNodeMgr) DoDelta_Weight(G []float64){
	for _, bNode := range mgr.nodeList {
		b_out := bNode.DataRecv
		for index, _ := range bNode.NodeList {
			bNode.NodeList[index].Weight += mgr.StepLen * G[index] * b_out
		}
	}
}


func (mgr *InNodeMgr) Clear_In()  {
	for _, node := range mgr.nodeList {
		node.DataRecv = 0
	}
}