package BpNet

import (
	// "fmt"
	"math"
)

type BNodeMgr struct {
	NodeMgr
	NodeMgr_In *InNodeMgr
}

type BWeightParams struct {
	StepLen float64
	B_weight []float64
	YParam  Y_Param
	YParamList []Y_Param
}

func BModifyWeight(b_get float64, inpurParam interface{}) float64 {
	params := inpurParam.(BWeightParams)
	value := params.YParam
	g := value.Get * (1 - value.Get) * (value.Real - value.Get)
	
	retn := params.StepLen * g * b_get
	return retn
}

func BModifyXValue(b_value float64, inpurParam interface{}) float64 {
	params := inpurParam.(BWeightParams)
	var Y_deltaSum float64 = 0
	for index, value := range params.YParamList {
		g := value.Get * (1 - value.Get) * (value.Real - value.Get)
		Y_deltaSum += g * params.B_weight[index]
	}
	E := b_value*(1-b_value)* Y_deltaSum

	return -1 * params.StepLen * E
}

func BDoInputFunc(input float64) float64 {
	// input -= delta
	// fmt.Println("BDoInputFunc in:", input)
	tmp := 1/(math.Exp(-1*input) + 1)
	// fmt.Println("BDoInputFunc out:", tmp)
	// return tmp*(1-tmp)
	return tmp
}

func (mgr *BNodeMgr) InitBNodeMgr(size int) {
	mgr.DoModifyConnWeight = BModifyWeight
	mgr.DoModifyValue = BModifyXValue
	mgr.DoFunc = BDoInputFunc
	mgr.mgrName = "Bnode"
	mgr.initNodes(size)
}

func (mgr *BNodeMgr) SetInNode(inMgr *InNodeMgr) {
	mgr.NodeMgr_In = inMgr
	for _, node := range mgr.nodeList {
		mgr.NodeMgr_In.AddConnect(node)
	}
}

func (mgr *BNodeMgr) DoDelta_B(Y_real []float64, Y_get []float64){
	for _, bNode := range mgr.nodeList {
		b_get := BDoInputFunc(bNode.DataRecv)
		var param BWeightParams
		param.StepLen = mgr.StepLen

		for index, yNode := range bNode.NodeList {
			var tmp Y_Param 
			tmp.Get = Y_get[index]
			tmp.Real = Y_real[index]
			param.YParam = tmp
			param.YParamList = append(param.YParamList, tmp)
			param.B_weight = append(param.B_weight, yNode.Weight)
			bNode.NodeList[index].Weight += BModifyWeight(b_get, param)
		}
		bNode.Value += BModifyXValue(b_get, param)
		// fmt.Println("node name:",bNode.Name, "node value:", bNode.Value)
	}
}
func (mgr *BNodeMgr) GetResult_B() []float64 {
	var retnList []float64
	for _, node := range mgr.nodeList {
		tmp := BDoInputFunc(node.DataRecv)
		retnList = append(retnList, tmp*(1-tmp))
	}
	return retnList
}

func (mgr *BNodeMgr) GetWeight_B() [][]float64 {
	var retnList [][]float64
	for _, node := range mgr.nodeList {
		
		var tmpList  []float64
		for _, YNode := range node.NodeList{
			tmpList = append(tmpList, YNode.Weight)
		}
		retnList = append(retnList, tmpList)
	}
	return retnList
}

func (mgr *BNodeMgr) Clear_B()  {
	for _, node := range mgr.nodeList {
		node.DataRecv = 0
	}
}