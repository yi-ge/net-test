package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/xuri/excelize/v2"
)

const (
	pingCount = 60
)

func main() {
	fmt.Println("程序正在运行，请打开当前目录下的'网络检测记录.xlsx'查看运行结果")

	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "测速开始时间")
	xlsx.SetCellValue("Sheet1", "B1", "测速结束时间")
	xlsx.SetCellValue("Sheet1", "C1", "主路由每分钟平均延迟")
	xlsx.SetCellValue("Sheet1", "D1", "主路由每分钟丢包率")
	xlsx.SetCellValue("Sheet1", "E1", "网关每分钟平均延迟")
	xlsx.SetCellValue("Sheet1", "F1", "网关每分钟丢包率")
	xlsx.SetCellValue("Sheet1", "G1", "百度每分钟平均延迟")
	xlsx.SetCellValue("Sheet1", "H1", "百度每分钟丢包率")
	xlsx.SetCellValue("Sheet1", "I1", "服务器每分钟平均延迟")
	xlsx.SetCellValue("Sheet1", "J1", "服务器每分钟丢包率")
	xlsx.SetCellValue("Sheet1", "K1", "电信机房网速下行")
	xlsx.SetCellValue("Sheet1", "L1", "全球网测网速下行")
	// 调整列宽
	const cellWidth = 2 // 每个字的宽度

	for i := 'A'; i <= 'L'; i++ { // 遍历列
		cellWidthValue := cellWidth
		if i < 'C' {
			cellWidthValue = cellWidth + 1
		} else {
			cellWidthValue = cellWidth
		}
		colName := string(i) + "1"                              // 获取单元格名称
		cellText, _ := xlsx.GetCellValue("Sheet1", colName)     // 获取单元格的文本值
		textLength := len([]rune(cellText))                     // 获取文本的长度
		width := float64(textLength) * float64(cellWidthValue)  // 计算单元格宽度
		xlsx.SetColWidth("Sheet1", string(i), string(i), width) // 设置单元格宽度
	}
	xlsx.SaveAs("网络检测记录.xlsx")

	for i := 2; ; i++ {
		startTime := time.Now()
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), startTime.Format("2006-01-02 15:04:05"))

		// Ping
		routerAvgRtt, routerPktLoss := pingTarget("192.168.1.1")
		gatewayAvgRtt, gatewayPktLoss := pingTarget("192.168.1.254")
		baiduAvgRtt, baiduPktLoss := pingTarget("www.baidu.com")
		serverAvgRtt, serverPktLoss := pingTarget("1.2.4.8")

		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i), routerAvgRtt)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i), routerPktLoss)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i), gatewayAvgRtt)
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i), gatewayPktLoss)
		xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(i), baiduAvgRtt)
		xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(i), baiduPktLoss)
		xlsx.SetCellValue("Sheet1", "I"+strconv.Itoa(i), serverAvgRtt)
		xlsx.SetCellValue("Sheet1", "J"+strconv.Itoa(i), serverPktLoss)

		xlsx.Save()

		// 下载测试
		downloadSpeed1 := downloadFile("https://example.com/file.zip?timestamp=" + strconv.FormatFloat(float64(time.Now().UnixNano()/1e9), 'f', -1, 64))
		downloadSpeed2 := downloadFile("https://example.com/file.zip")

		endTime := time.Now()
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), endTime.Format("2006-01-02 15:04:05"))

		// 将下载速度结果添加到xlsx中
		xlsx.SetCellValue("Sheet1", "K"+strconv.Itoa(i), downloadSpeed1)
		xlsx.SetCellValue("Sheet1", "L"+strconv.Itoa(i), downloadSpeed2)

		xlsx.Save()

		time.Sleep(time.Minute * 5)
	}
}

func pingTarget(target string) (avgRtt float64, pktLoss float64) {
	pinger, err := probing.NewPinger(target)
	if err != nil {
		fmt.Printf("Error creating pinger for %s: %v\n", target, err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = pingCount
	pinger.Timeout = time.Second
	pinger.Run()

	stats := pinger.Statistics()
	avgRtt = float64(stats.AvgRtt) / float64(time.Millisecond)
	pktLoss = stats.PacketLoss
	return
}

func downloadFile(url string) float64 {
	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return 0
	}
	defer resp.Body.Close()

	out, err := os.Create("temp_download")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return 0
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error writing to temp file:", err)
		return 0
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	fileInfo, err := out.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return 0
	}

	fileSize := float64(fileInfo.Size())
	downloadSpeed := fileSize / duration / 1024 / 1024 * 8 // Mbps

	err = os.Remove("temp_download")
	if err != nil {
		fmt.Println("Error deleting temp file:", err)
	}

	return downloadSpeed
}
