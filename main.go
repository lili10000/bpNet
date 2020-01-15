package main
import(
   . "bpNet/BpNet"
   "fmt"
   "time"
) 

func main()  {
    var node InNodeMgr

	nodeNameInit := "inNode"
    var nodeNameList []string
    nodeMap := make(map[string]float32)
	for i:=0; i< 4; i++{
		nodeName := fmt.Sprintf("%s_%d", nodeNameInit, i)
        nodeNameList = append(nodeNameList, nodeName)
        nodeMap[nodeName] = 0.1 + 0.1*float32(i)      
	}
    node.Init(nodeNameList)
    node.DoInput(nodeMap)
    time.Sleep(5 * time.Second) 
}