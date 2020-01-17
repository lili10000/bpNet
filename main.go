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

func dataChange(input float64) float64 {
	tmp := 1/(math.Exp(-1*input) + 1)
	return tmp
}

func main()  {

    dataList := getTrainData("today.txt")

    dataTest := dataList[:len(dataList)/2]
    dataTrain := dataList[len(dataList)/2:]


    var inNode InNodeMgr
    var bNode BNodeMgr
    var bNode_1 BNodeMgr
    var yNode YNodeMgr
    inNode.InitInNodeMgr(5)

    bNode.InitBNodeMgr(6)
    bNode.SetInNode(&inNode)

    bNode_1.InitBNodeMgr(7)
    bNode_1.SetBNode(&bNode)

    yNode.InitYNodeMgr(1)
    yNode.SetBNode(&bNode_1)

    for i:= 0; i < 1000; i++{
        for _, data := range dataTrain{
            if len(data) != 7{
                break
            }

            InData := data[:len(data)-2]
            

            rate := data[len(data)-2]
            rate = dataChange(rate)
            score := data[len(data)-1]
            score = dataChange(score)

            result := []float64{rate}

    
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

    // dataList = getTrainData("test.txt")

    okSum := 0
    errSum := 0
    for _, data := range dataTest{
        if len(data) != 7{
            break
        }
        InData := data[:len(data)-2]

        rate := data[len(data)-2]
        rate = dataChange(rate)
        score := data[len(data)-1]
        score = dataChange(score)

        result := []float64{rate}

        // resultTmp := data[len(data)-1]

        // resultTmp = resultTmp*0.3
        // result := []float64{resultTmp}

        inNode.DoInput(InData)
        inNode.DoSend()
        bNode.DoSend()
        bNode_1.DoSend()
        Y_get := yNode.GetResult_Y()

        // getIndex := (math.Pow((result[0] - Y_get[0]), 2) + math.Pow(math.Abs(result[1] - Y_get[1]), 2)) /2
        // fmt.Println("0:", result[0], "0:", Y_get[0],)
        // fmt.Println("1:", result[1], "1:", Y_get[1],)
        getIndex := ((result[0] < 0.5) && (Y_get[0] < 0.5)) || ((result[0] > 0.5) && (Y_get[0] > 0.5))

        if getIndex {
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
}