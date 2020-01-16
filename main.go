package main
import(
   . "bpNet/BpNet"
   "bufio"
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
    // "math"
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
    // fmt.Println(dataList)

    var inNode InNodeMgr
    var bNode BNodeMgr
    var yNode YNodeMgr
    inNode.InitInNodeMgr(4)

    bNode.InitBNodeMgr(6)
    bNode.SetInNode(&inNode)

    yNode.InitYNodeMgr(3)
    yNode.SetBNode(&bNode)

    for i:= 0; i < 100000; i++{
        for _, data := range dataList{
            if len(data) != 5{
                break
            }

            InData := data[:len(data)-1]
            resultTmp := data[len(data)-1]

            // result := []float64{resultTmp}
            result := []float64{0,0,0}
            result[int(resultTmp)-1] = 1
    
            inNode.DoInput(InData)
            inNode.DoSend()
            bNode.DoSend()
            Y_get := yNode.GetResult_Y()
            // fmt.Println("GET ", Y_get)
            B_get := bNode.GetResult_B()
            B_weight := bNode.GetWeight_B()  
            yNode.DoDelta_Y(result)
            bNode.DoDelta_B(result, Y_get)
            inNode.DoDelta_In(result, Y_get, B_get, B_weight)

            yNode.Clear_Y()
            bNode.Clear_B()
            inNode.Clear_In()
            // fmt.Println("REAL", result)
            // break
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

        result := []float64{0,0,0}

        if int(resultTmp) > 3 {
            fmt.Println(data)
        }
        result[int(resultTmp)-1] = 1

        inNode.DoInput(InData)
        inNode.DoSend()
        bNode.DoSend()
        Y_get := yNode.GetResult_Y()

        getIndex := checkBig(Y_get)
        realIndex := checkBig(result)
        
        if getIndex != realIndex {
            errSum += 1
        }else{
            okSum +=1
        }
        
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