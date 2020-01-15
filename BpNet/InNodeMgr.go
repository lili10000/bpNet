package BpNet

import(
	"fmt"
)

type InNodeMgr struct {
    NodeMgr
}

type Y_Param struct{
	BWeight float32
	Real float32
	Get float32
}

type InWeightParams struct{
	StepLen float32
	B float32
	Y_param []Y_Param
} 

func ModifyWeight(x_value float32, connNodeName string , inpurParam interface{}) float32{
	params := inpurParam.(InWeightParams)
	var Y_deltaSum float32 = 0
	for _, value := range params.Y_param{
		g := value.Get*(1-value.Get)*(value.Real - value.Get)
		Y_deltaSum += g*value.BWeight
	}
	E := params.B*(1-params.B)*Y_deltaSum
	retn := params.StepLen * E * x_value
	return retn
}

func ModifyXValue(x_value float32, inpurParam interface{}) float32{
	return x_value
}

func DoInputFunc(x_value float32)float32{
	return x_value
}

func (mgr *InNodeMgr) Init(dataStruct []string){
	mgr.DoModifyConnWeight = ModifyWeight
	mgr.DoModifyValue = ModifyXValue
	mgr.DoFunc = DoInputFunc
	mgr.initNodes(dataStruct)


}

func (mgr *InNodeMgr) DoInput(dataMap map[string]float32){
	for name,value := range dataMap{
		node,ok := mgr.nodeMap[name]
		if !ok {
			fmt.Println(name,"not in map")
			return
		}
		node.Channel <- value
		// fmt.Println("name", node.Name, "channel addr:", node.Channel)
	}
}




