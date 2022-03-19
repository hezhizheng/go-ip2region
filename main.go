package main

import (
	"encoding/json"
	"flag"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/thinkeridea/go-extend/exnet"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	wg    = sync.WaitGroup{}
	port  = ""
	d     = "" // 下载标识
	dbUrl = map[string]string{
		"1": "https://ghproxy.com/?q=https://github.com/lionsoul2014/ip2region/blob/master/data/ip2region.db?raw=true",
		"2": "https://ghproxy.com/?q=https://github.com/bqf9979/ip2region/blob/master/data/ip2region.db?raw=true",
	}
)

const (
	ipDbPath     = "./ip2region.db"
	defaultDbUrl = "2" // 默认下载 来自 lionsoul2014 仓库的 ip db文件
)

type JsonRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type IpInfo struct {
	Ip       string `json:"ip"`
	Country  string `json:"country"`  // 国家
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	County   string `json:"county"`   // 县、区
	Region   string `json:"region"`   // 区域位置
	ISP      string `json:"isp"`      // 互联网服务提供商
}

func init() {
	_p := flag.String("p", "9090", "本地监听的端口")
	_d := flag.String("d", "0", "仅用于下载最新的ip地址库，保存在当前目录")
	flag.Parse()

	port = *_p
	d = *_d

	if d != "0" {
		if value, ok := dbUrl[d]; ok {
			downloadIpDb(value)
		} else {
			downloadIpDb(dbUrl[defaultDbUrl])
		}
		os.Exit(1)
	}

	checkIpDbIsExist()
}

func main() {
	http.HandleFunc("/", queryIp)

	link := "http://127.0.0.1:" + port

	log.Println("监听端口", link)
	listenErr := http.ListenAndServe(":"+port, nil)
	if listenErr != nil {
		log.Fatal("ListenAndServe: ", listenErr)
	}
}

func queryIp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")

	defer func() {
		//捕获 panic
		if err := recover(); err != nil {
			log.Println("查询ip发生错误", err)
		}
	}()

	r.ParseForm() // 解析参数

	ip := r.FormValue("ip")

	if ip == "" {
		// 获取当前客户端 IP
		ip = getIp(r)
	}

	region, err := ip2region.New(ipDbPath)
	defer region.Close()

	if err != nil {
		msg, _ := json.Marshal(&JsonRes{Code: 4001, Msg: err.Error()})
		w.Write(msg)
		return
	}

	info, searchErr := region.BinarySearch(ip)

	if searchErr != nil {
		msg, _ := json.Marshal(JsonRes{Code: 4002, Msg: searchErr.Error()})
		w.Write(msg)
		return
	}

	// 赋值查询结果
	ipInfo := &IpInfo{
		Ip:       ip,
		ISP:      info.ISP,
		Country:  info.Country,
		Province: info.Province,
		City:     info.City,
		County:   "",
		Region:   info.Region,
	}

	msg, _ := json.Marshal(JsonRes{Code: 200, Data: ipInfo})
	w.Write(msg)
	return
}

func getIp(r *http.Request) string {
	ip := exnet.ClientPublicIP(r)
	if ip == "" {
		ip = exnet.ClientIP(r)
	}
	return ip
}

func checkIpDbIsExist() {
	if _, err := os.Stat(ipDbPath); os.IsNotExist(err) {
		log.Println("ip 地址库文件不存在")
		downloadIpDb(dbUrl[defaultDbUrl])
	}
}

func downloadIpDb(url string) {
	log.Println("正在下载最新的 ip 地址库...：" + url)
	wg.Add(1)
	go func() {
		downloadFile(ipDbPath, url)
		wg.Done()
	}()
	wg.Wait()
	log.Println("下载完成")
}

// @link https://studygolang.com/articles/26441
func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
