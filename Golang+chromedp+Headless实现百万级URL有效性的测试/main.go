package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
	"strings"
	"io/ioutil"
	"os"
	// "sync"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Score struct {
	Num int
	Url string
}

var  wg sync.WaitGroup

var statusMap map[int]string = map[int]string{
	301: "301",
	302: "302",
	304: "304",
	400: "400",
	401: "401",
	403: "403",
	404: "404",
	405: "405",
	408: "408",
	500: "500",
	502: "502",
	503: "503",
	504: "504",
}

var flag bool

var DispatchNumControl = make(chan bool, 10000)

func (s *Score) Do() {
	fmt.Printf("num:%d  url:%s\n", s.Num, s.Url)
	// fmt.Println("num:", s.Url)
	// opts := chromedp.ExecPath("C:\\Users\\10049\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe")

	// allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts)
	// defer cancel()

	// taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	// defer cancel()

	var buf []byte
	// ch := make(chan []byte, 1024)
	// create chrome instance
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
		// fmt.Println(err)
	}

	// listen network event
	listenForNetworkEvent(taskCtx) //,ch

	chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(s.Url),
		// chromedp.WaitVisible(`#doc-info`), //, chromedp.BySearch
		chromedp.CaptureScreenshot(&buf),
	)

	log.Println("hello")
	if flag == true{
		pngname := strings.Replace(s.Url,".","1",-1)
		pngname = strings.Replace(pngname,"http://","1",-1)
		log.Println(buf)
		if err1 := ioutil.WriteFile(pngname+".png", buf, 0644); err1 != nil {
			log.Fatal(err1)
		}
	}

	time.Sleep(30 * time.Second)
}

func listenForNetworkEvent(ctx context.Context) { //, ch chan []byte
	count := 0
	chromedp.ListenTarget(ctx, func(ev interface{}) {

		switch ev := ev.(type) {

		case *network.EventResponseReceived:
			count++
			// log.Println("count=", count)
			resp := ev.Response
			if count < 2 {
				if resp.Status == 200 {
					log.Printf("Success, url:%s; status:%d; text:%s", resp.URL, resp.Status, resp.StatusText)
					writefile("./success.txt", "Url:"+resp.URL+" Status:200"+" Text:"+resp.StatusText+"\r\n")
					flag = false
				} else {
					log.Printf("Failed, url:%s; status:%d; text:%s", resp.URL, resp.Status, resp.StatusText)
					writefile("./failed.txt", "Url:"+resp.URL+" Status:"+statusMap[int(resp.Status)]+" Text:"+resp.StatusText+"\r\n")
					flag = true
				}
			}
		}
		// other needed network Event
	})
}

func writefile(filename, msg string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644) //表示最佳的方式打开文件，如果不存在就创建，打开的模式是可读可写，权限是644
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(msg)
	f.Close()
}

func main() {
	num := 500
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	// wg.Add(1)
	p := NewWorkerPool(num)
	p.Run()

	// 读取文件，以切片的形式返回 URLS
	txtpath := "F:\\GO_code\\test_xlsx\\top-1m.xlsx"
	urls := loadfile(txtpath)
	log.Printf("urls[0]:%s  type:%T   urls'length:%d", urls[0], urls[0], len(urls)) // print the content as a 'string'
	datanum := len(urls)
	go func() {
		for i := 0; i < datanum; i++ {
			// wg.Add(1)
			sc := &Score{Num: i, Url: urls[i]}
			p.JobQueue <- sc
			// fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
			// wg.Wait()
			time.Sleep(100 * time.Millisecond) // 700 * time.Millisecond
		}
	}()
	// wg.Wait()
	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}

}
