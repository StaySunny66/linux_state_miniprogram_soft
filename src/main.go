package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
    "net/http"
	"time"
	"encoding/json"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/skip2/go-qrcode"
    "os"
    "encoding/base64"
    "io/ioutil"
	"github.com/labstack/gommon/log"



)


type Status struct {
  T   time.Time `json:"-"`
  Cpu float64   `json:"cpu"`
  Mem struct {
    Current uint64 `json:"current"`
    Total   uint64 `json:"total"`
  } `json:"mem"`
  Swap struct {
    Current uint64 `json:"current"`
    Total   uint64 `json:"total"`
  } `json:"swap"`
  Disk struct {
    Current uint64 `json:"current"`
    Total   uint64 `json:"total"`
  } `json:"disk"`
  Uptime   uint64    `json:"uptime"`
  Loads    []float64 `json:"loads"`
  TcpCount int       `json:"tcpCount"`
  UdpCount int       `json:"udpCount"`
  NetIO    struct {
    Up   uint64 `json:"up"`
    Down uint64 `json:"down"`
  } `json:"netIO"`
  NetTraffic struct {
    Sent uint64 `json:"sent"`
    Recv uint64 `json:"recv"`
	SentM uint64 `json:"sentM"`
    RecvM uint64 `json:"recvM"`
  } `json:"netTraffic"`
}



func show_qrcode(){

	fmt.Println("欢迎使用服务器监控脚本")
	fmt.Println("------------------稍等片刻-----------------")
	fmt.Println("获取公网IP")


	resp, err := http.Get("https://api4.ipify.org/?format=text")
	if err != nil {
		fmt.Println("无法获取IP地址：", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("无法读取响应：", err)
		return
	}

	local_ip := string(body)
	

	fmt.Println("你的外部IP地址是：", local_ip)
	fmt.Println("服务器端口号", 1323)

	url:= "http://"+local_ip+":1323/"

	data := []byte(url)
	// 进行Base64编码
	encoded := base64.StdEncoding.EncodeToString(data)
	fmt.Println("绑定ID", encoded)



	// 生成二维码
	qr, err := qrcode.New(encoded, qrcode.Medium)
	if err != nil {
		fmt.Println("生成二维码时出错：", err)
		os.Exit(1)
	}


	fmt.Println("\n")
	
	fmt.Println("请使用微信小程序 矢光小屋 扫描下方二维码 添加服务器")
	fmt.Println("\n")

	// 在终端中显示二维码
	fmt.Println(qr.ToSmallString(false))
}

func customLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取Echo的默认日志记录器
			logger := c.Echo().Logger
			// 设置日志级别为致命错误
			logger.SetLevel(log.ERROR)
			// 执行下一个中间件或处理程序
			return next(c)
		}
	}
}



func GetStatus(lastStatus *Status) *Status {
	now := time.Now()
	status := &Status{
		T: now,
	}

	percents , err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("get cpu percent failed:", err)
	} else {
		status.Cpu = percents[0]
	}

	upTime , err := host.Uptime()
	if err != nil {
		fmt.Println("get uptime failed:", err)
	} else {
		status.Uptime = upTime
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("get virtual memory failed:",err)
	} else {
		status.Mem.Current = memInfo.Used
		status.Mem.Total = memInfo.Total
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		fmt.Println("get swap memory failed:", err)
	} else {
		status.Swap.Current = swapInfo.Used
		status.Swap.Total = swapInfo.Total
	}

	distInfo,err := disk.Usage("/")
	if err != nil {
		fmt.Println("get dist usage failed:", err)
	} else {
		status.Disk.Current = distInfo.Used/1024/1024
		status.Disk.Total = distInfo.Total/1024/1024
	}

	avgState, err := load.Avg()
	if err != nil {
		fmt.Println("get load avg failed:",err)
	} else {
		status.Loads = []float64{avgState.Load1, avgState.Load5, avgState.Load15}
	}

	ioStats, err := net.IOCounters(false)
	if err != nil {
		fmt.Println("get io counters failed:", err)
	} else if len(ioStats) > 0 {
		ioStat := ioStats[0]
		status.NetTraffic.Sent = ioStat.BytesSent
		status.NetTraffic.Recv = ioStat.BytesRecv
		status.NetTraffic.SentM = ioStat.BytesSent/1024/1024
		status.NetTraffic.RecvM = ioStat.BytesRecv/1024/1024

		if lastStatus != nil {
			duration := now.Sub(lastStatus.T)
			seconds := float64(duration) / float64(time.Second)
			up := uint64(float64(status.NetTraffic.Sent-lastStatus.NetTraffic.Sent) / seconds)
			down := uint64(float64(status.NetTraffic.Recv-lastStatus.NetTraffic.Recv) / seconds)
			status.NetIO.Up = up
			status.NetIO.Down = down
		}
	} else {
		fmt.Println("can not find io counters")
	}

	status.TcpCount, err = 0,nil // GetTCPCount()
	if err != nil {
		fmt.Println("get tcp connections failed:", err)
	}

	status.UdpCount, err = 0,nil //  GetUDPCount()
	if err != nil {
		fmt.Println("get udp connections failed:", err)
	}

	return status
}


var S *Status



func main() {


  // Echo instance
  e := echo.New()

  // Middleware
  //e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(customLogger())

  // Routes
  e.GET("/", hello)
  e.GET("/", state)

  show_qrcode()
  
  // Start server
  e.Logger.Fatal(e.Start(":1323"))

  e.Logger.SetOutput(ioutil.Discard)

}

// Handler
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}

func state(c echo.Context) error {

  
  S = GetStatus(S);

  bytes, err := json.Marshal(S)
  if err != nil {
    panic(err)
  }
  return c.String(http.StatusOK, string(bytes))

}