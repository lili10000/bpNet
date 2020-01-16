package BpNet

import (
	// "fmt"
	"math"
)

type YWeightParams struct {
	StepLen float64
	Y  Y_Param
}

func YModifyWeight(b_value float64, inpurParam interface{}) float64 {
	return b_value
}

func YModifyXValue(input float64, inpurParam interface{}) float64 {
	params := inpurParam.(YWeightParams)
	value := params.Y
	g := value.Get * (1 - value.Get) * (value.Real - value.Get)
	
	retn := -1 * params.StepLen * g
	return retn
}

func YDoInputFunc(input float64) float64 {
	// input -= delta
	retn := 1/(math.Exp(-1*input) + 1)
	return retn
}

type YNodeMgr struct {
	NodeMgr
	NodeRecvData map[string]float64
	NodeMgr_B    *BNodeMgr
}

func (mgr *YNodeMgr) InitYNodeMgr(size int) {
	mgr.DoModifyConnWeight = YModifyWeight
	mgr.DoModifyValue = YModifyXValue
	mgr.DoFunc = YDoInputFunc
	mgr.mgrName = "Ynode"
	mgr.initNodes(size)
}

func (mgr *YNodeMgr) SetBNode(inMgr *BNodeMgr) {
	mgr.NodeMgr_B = inMgr
	for _, node := range mgr.nodeList {
		mgr.NodeMgr_B.AddConnect(node)
	}
}

func (mgr *YNodeMgr) GetResult_Y() []float64 {
	var retnList []float64
	for _, node := range mgr.nodeList {
		retnList = append(retnList, YDoInputFunc(node.DataRecv))
	}
	return retnList
}

// func (mgr *YNodeMgr) GetDelta(realValue []float64) float64 {
// 	if len(realValue) != len(mgr.nodeList){
// 		fmt.Println("GetDelta: input dataList size: ",len(realValue), " != nodeList size:", len(mgr.nodeList))
// 		return 0
// 	}
// 	deltaSum := float64(0)

// 	for index, node := range mgr.nodeList {
// 		deltaSum += math.Pow((node.Value - realValue[index]), 2)
// 	}

// 	return deltaSum/2
// }
func (mgr *YNodeMgr) DoDelta_Y(Y_real []float64){

	for index, node := range mgr.nodeList {
		var param YWeightParams
		param.StepLen = mgr.StepLen
		param.Y.Get = YDoInputFunc(node.DataRecv)
		param.Y.Real = Y_real[index]
		node.Value += YModifyXValue(node.Value, param)
	}
}

func (mgr *YNodeMgr) Clear_Y()  {
	for _, node := range mgr.nodeList {
		node.DataRecv = 0
	}
}