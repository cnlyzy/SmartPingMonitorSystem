package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/tidwall/gjson"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"time"
)

var alertNum = 0                 //报警次数 每次不一样
var pushTime = time.Now().Unix() //上次推送时间
var conf, _ = goconfig.LoadConfigFile("conf.ini")

func main() {
	monitorRate, _ := conf.Int64("SmartPing", "monitorRate")
	ticker := time.NewTicker(time.Second * time.Duration(monitorRate)) //可调节的监控频率
	for range ticker.C {
		monitor()
	}
	//ticker.Stop()
	fmt.Println("Ticker stopped")
}

/**
 * 监控
 */
func monitor() {
	Alerts := getAlert()
	nowAlertNum := gjson.Get(string(Alerts), "#").Int()
	fmt.Printf("nowAlert:%d  alreadyAlert:%d  nowTime:%d  pushTime:%d \r", int(nowAlertNum), int(alertNum), time.Now().Unix(), pushTime)

	// 如果报警数 < alertNum 则 alertNum = 报警数 （解决隔天报警数清空问题）
	if (int(nowAlertNum) < alertNum) {
		alertNum = 0
		pushTime = time.Now().Unix()
		return
	}
	setAlertNum, _ := conf.Int64("MonitorSystem", "setAlertNum")   //设置的报警数（指定时间内报警数量达到这个值则报警）
	setAlertTime, _ := conf.Int64("MonitorSystem", "setAlertTime") //多少时间发送一次报警邮件
	// 如果 当前时间戳 - pushTime >= 设置的推送时间(多久推送一次)
	if time.Now().Unix()-pushTime < setAlertTime {
		return
	}
	// 如果 报警数 > alertNum 则 发送报警邮件 并 更新pushTime 和 alertNum
	alertTimes := int(nowAlertNum) - alertNum
	if int64(alertTimes) < setAlertNum {
		return
	}
	//执行报警动作
	sendMail(alertTimes)
	pushTime = time.Now().Unix()
	alertNum = int(nowAlertNum)
}

/**
 * 发邮件
 */
func sendMail(alertTimes int) {
	from, _ := conf.GetValue("Email", "from")
	to, _ := conf.GetValue("Email", "to")
	title, _ := conf.GetValue("Email", "title")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	timeLayout := "2006-01-02 15:04:05"
	//邮件内容
	body := fmt.Sprintf("SmartPing 服务器从 %s 到 %s 共产生<font size='3' color='red'> %d </font> 次报警，请注意检查!", time.Unix(pushTime, 0).Format(timeLayout), time.Unix(time.Now().Unix(), 0).Format(timeLayout), alertTimes)
	m.SetBody("text/html", body)

	host, _ := conf.GetValue("Email", "host")
	port, _ := conf.Int("Email", "port")
	username, _ := conf.GetValue("Email", "username")
	password, _ := conf.GetValue("Email", "password")
	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

/**
 * 从SmartPing获取当前的报警数
 */
func getAlert() (string) {
	//url := "http://192.168.2.182:8899/api/alert.json"
	uri, _ := conf.GetValue("SmartPing", "url")
	port, _ := conf.GetValue("SmartPing", "port")
	url := fmt.Sprintf("%s:%s/api/alert.json", uri, port)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	value := gjson.Get(string(body), "1")
	return value.Raw
}
