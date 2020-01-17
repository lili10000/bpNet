package BpNet

import (
	// "fmt"
	"math"
)

type BNodeMgr struct {
	NodeMgr
	NodeMgr_In *InNodeMgr
	NodeMgr_B *BNodeMgr
}

type BWeightParams struct {
	StepLen float64
	B_weight []float64
	YParam  Y_Param
	YParamList []Y_Param
}

type G_data struct {
	Real float64
	Get float64
}

func BModifyWeight(b_out float64, inpurParam interface{}) float64 {
	params := inpurParam.(BWeightParams)
	value := params.YParam
	g := value.Get * (1 - value.Get) * (value.Real - value.Get)
	
	retn := params.StepLen * g *b_out
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
	tmp := 1/(math.Exp(-1*input) + 1)
	return tmp
}

func BDoRealIn(realOut float64) float64{
	return -1 * math.Log(1/realOut - 1) / math.Log(math.E)
}


func (mgr *BNodeMgr) InitBNodeMgr(size int) {
	mgr.DoModifyConnWeight = BModifyWeight
	mgr.DoModifyValue = BModifyXValue
	mgr.DoFunc = BDoInputFunc
	mgr.DoGetRealIn = BDoRealIn
	mgr.mgrName = "Bnode"
	mgr.initNodes(size)
}

func (mgr *BNodeMgr) SetInNode(inMgr *InNodeMgr) {
	mgr.NodeMgr_In = inMgr
	for _, node := range mgr.nodeList {
		mgr.NodeMgr_In.AddConnect(node)
	}
}

func (mgr *BNodeMgr) SetBNode(bMgr *BNodeMgr) {
	mgr.NodeMgr_B = bMgr
	for _, node := range mgr.nodeList {
		mgr.NodeMgr_B.AddConnect(node)
	}
}

func (mgr *BNodeMgr) DoDelta_B(Y_real []float64, Y_get []float64){
	for _, bNode := range mgr.nodeList {
		b_out := BDoInputFunc(bNode.DataRecv - bNode.Value)
		var param BWeightParams
		param.StepLen = mgr.StepLen

		for index, yNode := range bNode.NodeList {
			var tmp Y_Param 
			tmp.Get = Y_get[index]
			tmp.Real = Y_real[index]
			param.YParam = tmp
			param.YParamList = append(param.YParamList, tmp)
			param.B_weight = append(param.B_weight, yNode.Weight)
			bNode.NodeList[index].Weight += BModifyWeight(b_out, param)
		}
		bNode.Value += BModifyXValue(b_out, param)
		// fmt.Println("node name:",bNode.Name, "node value:", bNode.Value)
	}
}

func (mgr *BNodeMgr) DoDelta_Weight(G []float64){
	for _, bNode := range mgr.nodeList {
		b_out := BDoInputFunc(bNode.DataRecv - bNode.Value)
		for index, _ := range bNode.NodeList {
			bNode.NodeList[index].Weight += mgr.StepLen * G[index] * b_out
		}
	}
}

func (mgr *BNodeMgr) DoDelta_value(G []float64){
	for i, bNode := range mgr.nodeList {
		b_out := BDoInputFunc(bNode.DataRecv - bNode.Value)
		sumTmp := float64(0)
		for index, node := range bNode.NodeList {
			sumTmp += node.Weight * G[index] 
		}
		mgr.nodeList[i].Value +=  -1 * b_out*(1 - b_out) * sumTmp * mgr.StepLen
	}
}




func (mgr *BNodeMgr) GetResult_B() []float64 {
	var retnList []float64
	for _, node := range mgr.nodeList {
		tmp := BDoInputFunc(node.DataRecv - node.Value)
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

func (mgr *BNodeMgr) GetResult_G_from_Y(Y_real []float64, Y_get []float64) []float64 {
	var retnList []float64
	for index, _ := range Y_real {

		Get := Y_get[index]
		Real := Y_real[index]

		G := Get*(1-Get)*(Real - Get)
		retnList = append(retnList, G)
	}

	return retnList
}


func (mgr *BNodeMgr) GetResult_E_from_Y(Y_real []float64, Y_get []float64) []float64 {
	var retnList []float64
	for _, bNode := range mgr.nodeList {
		b_out := BDoInputFunc(bNode.DataRecv - bNode.Value)

		sumTmp := float64(0)
		for index, yNode := range bNode.NodeList {

			Get := Y_get[index]
			Real := Y_real[index]

			G := Get*(1-Get)*(Real - Get)
			weight := yNode.Weight

			sumTmp += G * weight
		}
		E := b_out*(1-b_out)*sumTmp
		retnList = append(retnList, E)
	}
	return retnList
}


func (mgr *BNodeMgr) GetResult_E_from_B(E []float64) []float64 {
	var retnList []float64
	for _, bNode := range mgr.nodeList {
		b_out := BDoInputFunc(bNode.DataRecv - bNode.Value)

		sumTmp := float64(0)
		for index, yNode := range bNode.NodeList {
			G := E[index]
			weight := yNode.Weight

			sumTmp += G * weight
		}
		E := b_out*(1-b_out)*sumTmp
		retnList = append(retnList, E)
	}
	return retnList
}




func (mgr *BNodeMgr) Clear_B()  {
	for _, node := range mgr.nodeList {
		node.DataRecv = 0
	}
}