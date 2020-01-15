package main
import(
   . "bpNet/BpNet"
   "fmt"
   "time"
) 

func main()  {
    var inNode InNodeMgr
    var bNode BNodeMgr
    var yNode YNodeMgr

	nodeNameInit := "inNode"
    var nodeNameList []string
    nodeMap := make(map[string]float32)
	for i:=0; i< 4; i++{
		nodeName := fmt.Sprintf("%s_%d", nodeNameInit, i)
        nodeNameList = append(nodeNameList, nodeName)
        nodeMap[nodeName] = 0.1 + 0.1*float32(i)      
    }
    
    inNode.InitInNodeMgr(nodeNameList)
    
    bNode.InitBNodeMgr(5)
    bNode.SetInNode(&inNode)

    yNode.InitYNodeMgr(1)
    yNode.SetBNode(&bNode)

    inNode.DoInput(nodeMap)
    inNode.DoSend()
    bNode.DoSend()
    fmt.Println(yNode.CheckResult())

    time.Sleep(1 * time.Second) 
}