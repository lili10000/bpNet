package main
import(
   . "bpNet/BpNet"
   "bufio"
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
    "math"
) 

func getTrainData(filePath string) ([][]float64){
    var retn [][]float64
    fi, err := os.Open(filePath)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return nil
    }
    defer fi.Close()

    br := bufio.NewReader(fi)
    for {
        a, _, c := br.ReadLine()
        if c == io.EOF {
            break
        }
        tmpList := strings.Split(string(a), ",")
        var lineData []float64
        for _, str := range tmpList{
            data, _ := strconv.ParseFloat(str, 64)
            lineData = append(lineData, data)
        }
        retn = append(retn, lineData)
    }
    return retn
}

func checkBig(input []float64) int{
    maxValue := float64(0)
    maxIndex := 0


    for index, value := range input{
        if value > maxValue{
            maxIndex = index
            maxValue = value
        }
    }
    return maxIndex
}


func main()  {

    dataList := getTrainData("train.txt")

    var inNode InNodeMgr
    var bNode BNodeMgr
    var bNode_1 BNodeMgr
    var yNode YNodeMgr
    inNode.InitInNodeMgr(4)

    bNode.InitBNodeMgr(5)
    bNode.SetInNode(&inNode)

    bNode_1.InitBNodeMgr(6)
    bNode_1.SetBNode(&bNode)

    yNode.InitYNodeMgr(1)
    yNode.SetBNode(&bNode_1)

    for i:= 0; i < 10000; i++{
        for _, data := range dataList{
            if len(data) != 5{
                break
            }

            InData := data[:len(data)-1]
            resultTmp := data[len(data)-1]

            resultTmp = resultTmp*0.3
            result := []float64{resultTmp}

    
            inNode.DoInput(InData)
            inNode.DoSend()
            bNode.DoSend()
            bNode_1.DoSend()
            Y_get := yNode.GetResult_Y()

            yNode.DoDelta_Y(result)
            G_get := bNode_1.GetResult_G_from_Y(result, Y_get)
            bNode_1.DoDelta_Weight(G_get)
            bNode_1.DoDelta_value(G_get)

            E_get := bNode_1.GetResult_E_from_Y(result, Y_get)
            bNode.DoDelta_Weight(E_get)
            bNode.DoDelta_value(E_get)
            E_get = bNode.GetResult_E_from_B(E_get)

            inNode.DoDelta_Weight(E_get)

            yNode.Clear_Y()
            bNode.Clear_B()
            bNode_1.Clear_B()
            inNode.Clear_In()
        }
    }

    dataList = getTrainData("test.txt")

    okSum := 0
    errSum := 0
    for _, data := range dataList{
        if len(data) != 5{
            break
        }
        InData := data[:len(data)-1]
        resultTmp := data[len(data)-1]

        resultTmp = resultTmp*0.3
        result := []float64{resultTmp}

        inNode.DoInput(InData)
        inNode.DoSend()
        bNode.DoSend()
        bNode_1.DoSend()
        Y_get := yNode.GetResult_Y()

        getIndex := math.Abs(result[0] - Y_get[0]) <= 0.02
        // fmt.Println("real:", result[0], "get:", Y_get[0], "result:", math.Abs(result[0] - Y_get[0]), getIndex )

        if getIndex  {
            okSum += 1
        }else{
            errSum  +=1
        }

        yNode.Clear_Y()
        bNode.Clear_B()
        bNode_1.Clear_B()
        inNode.Clear_In()
        
    }
    fmt.Println("ok:",okSum, "err:", errSum)



    // var nodeList []float64
    // result := []float64{0,1,0} 



    // for i:= 0; i < 10; i++{
    //     inNode.DoInput(nodeList)
    //     inNode.DoSend()
    //     bNode.DoSend()
    //     Y_get := yNode.GetResult_Y()
    //     fmt.Println(Y_get,)
    //     B_get := bNode.GetResult_B()
    //     B_weight := bNode.GetWeight_B()
    //     yNode.DoDelta_Y(result)
    //     bNode.DoDelta_B(result, Y_get)
    //     inNode.DoDelta_In(result, Y_get, B_get, B_weight)
    // }
}