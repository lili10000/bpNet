package BpNet

import(
	"fmt"
)

type NodeMgr struct {
	nodeMap map[string] Node
	DoModifyValue ModifyValue
	DoModifyConnWeight ModifyConnWeight
	DoFunc DoFuncByInput
}


func (mgr *NodeMgr) initNodes(dataStruct []string) {
	mgr.nodeMap = make(map[string] Node)
	
	for _,value := range dataStruct{
		var node Node
		node.Init(len(dataStruct))
		node.Name = value
		node.DoModifyValue = mgr.DoModifyValue
		node.DoModifyConnWeight = mgr.DoModifyConnWeight
		node.DoFunc = mgr.DoFunc
		mgr.nodeMap[value] = node
		go node.Update()
		fmt.Println("in", node.Channel)
	}

	for _, node := range mgr.nodeMap{
		fmt.Println("out", node.Channel)
	}

}

//
func (mgr *NodeMgr) AddConnect(connNode *Node) {
	for _,node := range mgr.nodeMap{
		err := node.Connect(connNode.Name, connNode.Channel)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func (mgr *NodeMgr) DoModify(connNode Node, nodeName string,valueParams interface{}, weightParams interface{}){
	for _,node := range mgr.nodeMap{
		node.Value = node.DoModifyValue(node.Value, valueParams)
		weight := node.DoModifyConnWeight(node.Value, nodeName, weightParams)
		node.SetWeigh(nodeName, weight)
	}
}

