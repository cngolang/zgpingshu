// zgpingshu project main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"./utils"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Audiobooks struct {
	item_id       int
	category_name string
	save_name     string
	url           string
	down_url      string
}

func doAnalysis(siteurl string) {
	//	res, err := http.Get("http://shantianfang.zgpingshu.com/575/#play")
	//	res, err := http.Get("http://shantianfang.zgpingshu.com/1040/#play")
	//	res, err := http.Get("http://yueyu.zgpingshu.com/3130/#play")
	res, err := http.Get(siteurl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	listlength := doc.Find(".down").Length()
	categoryName := doc.Find("#categoryHeader > h1").Text()
	categoryName, _ = DecodeToGBK(categoryName)
	fmt.Printf("专辑名为:%s, 共%d个音频文件.\n", categoryName, listlength)
	var audiobooks_slice []Audiobooks

	doc.Find(".player").Each(func(i int, s *goquery.Selection) {
		nameNode := s.Find("a").First()
		title := nameNode.Text()
		title, _ = DecodeToGBK(title)
		var audiobooks Audiobooks
		audiobooks.item_id = i
		audiobooks.save_name = title + ".mp3"
		audiobooks.category_name = categoryName
		audiobooks_slice = append(audiobooks_slice, audiobooks)
	})
	doc.Find(".down").Each(func(i int, s *goquery.Selection) {
		nameNode := s.Find("a").First()
		title := nameNode.Text()
		title, _ = DecodeToGBK(title)
		itemurl, _ := nameNode.Attr("href")
		itemurl, _ = DecodeToGBK(itemurl)
		audiobooks_slice[i].url = "http:" + itemurl
		//		if i < 4 {
		savePath := audiobooks_slice[i].category_name + "/" + audiobooks_slice[i].save_name
		exist, err := PathExists(savePath)
		if err != nil {
			fmt.Printf("get dir error![%v]\n", err)
			return
		}
		if !exist {
			if max_queue_downurl == -1 {
				downurl, err := getDownUrl(itemurl)
				status := "ok"
				if err != nil {
					status = "failed"
				}
				fmt.Printf("获取第%d个音频下载地址： %s \n", i+1, status)
				audiobooks_slice[i].down_url = downurl
			} else if max_queue_downurl > 0 {
				downurl, err := getDownUrl(itemurl)
				status := "ok"
				if err != nil {
					status = "failed"
				}
				fmt.Printf("获取第%d个音频下载地址： %s \n", i+1, status)
				audiobooks_slice[i].down_url = downurl
				max_queue_downurl--
			}
		}
		//		}
	})

	for _, audiobooks := range audiobooks_slice {
		if audiobooks.down_url != "" {
			if maxqueue == -1 {
				start := time.Now()
				_, err := DownAudiobookByFileload(audiobooks)
				//			fmt.Println("地址: ", audiobooks.down_url)
				status := "ok"
				if err != nil {
					status = "failed"
				}
				fmt.Printf("下载第%d个音频： %s  Total cost: %s \n", audiobooks.item_id+1, status, time.Since(start))
			} else if maxqueue > 0 {
				start := time.Now()
				_, err := DownAudiobookByFileload(audiobooks)
				//			fmt.Println("地址: ", audiobooks.down_url)
				status := "ok"
				if err != nil {
					status = "failed"
				}
				fmt.Printf("下载第%d个音频： %s  Total cost: %s \n", audiobooks.item_id+1, status, time.Since(start))
				maxqueue--
			}
		} else {
			//			fmt.Printf("%d  %v \n", key, audiobooks)
		}
	}
}
func printSlice(x []Audiobooks) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}

var maxqueue, max_queue_downurl int

func main() {
	flag.IntVar(&maxqueue, "q", -1, "maxqueue")
	flag.Parse()
	max_queue_downurl = maxqueue
	fmt.Println("请输入专辑地址", maxqueue)
	if len(os.Args) == 1 {
		fmt.Println("请输入专辑地址")
		return
	} else {
		url := os.Args[len(os.Args)-1]
		fmt.Printf("参数长度 %d 参数1： %s 参数url : %s \n", len(os.Args), os.Args[0], url)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		doAnalysis(url)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		//	var str string
		//	fmt.Scan(&str)
	}
}

// 获取单个音频下载地址
func getDownUrl(url string) (string, error) {
	res, err := http.Get("http:" + url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	nameNode := doc.Find("#down")
	itemurl, _ := nameNode.Attr("href")
	itemurl, _ = DecodeToGBK(itemurl)
	return string(itemurl), nil
}

// 实现单个文件的下载
func DownAudiobook(audiobooks Audiobooks) (string, error) {
	savePath := audiobooks.category_name + "/" + audiobooks.save_name
	exist, err := PathExists(savePath)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return string(savePath), nil
	}
	if exist {
		//		fmt.Printf("[%v] was downloaded\n", savePath)
		return string(savePath), nil
	}
	client := &http.Client{}
	url := audiobooks.down_url
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Cookie", "status_convertTime=00%3A00; statuscurrentTime=0.6635")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	reqest.Header.Add("Referer", audiobooks.url)
	if err != nil {
		panic(err)
	}
	res, _ := client.Do(reqest)
	Mkdir(audiobooks.category_name)
	//	savePath := audiobooks.save_name

	f, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	defer res.Body.Close()
	return string(savePath), nil
}

// 使用fileload实现单个文件的下载
func DownAudiobookByFileload(audiobooks Audiobooks) (string, error) {
	savePath := audiobooks.category_name + "/" + audiobooks.save_name
	exist, err := PathExists(savePath)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return string(savePath), nil
	}
	if exist {
		//		fmt.Printf("[%v] was downloaded\n", savePath)
		return string(savePath), nil
	} else {
		Mkdir(audiobooks.category_name)
		utils.Main1(5, savePath, audiobooks.down_url)
	}
	return string(savePath), nil
}

// 使用socks5代理下载
func DownAudiobookBySocks(audiobooks Audiobooks) (string, error) {
	savePath := audiobooks.category_name + "/" + audiobooks.save_name
	exist, err := PathExists(savePath)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return string(savePath), nil
	}
	if exist {
		//		fmt.Printf("[%v] was downloaded\n", savePath)
		return string(savePath), nil
	}
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:8889", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	client := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	url := audiobooks.down_url
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Cookie", "status_convertTime=00%3A00; statuscurrentTime=0.6635")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	reqest.Header.Add("Referer", audiobooks.url)
	if err != nil {
		panic(err)
	}
	response, _ := client.Do(reqest)
	//	savePath := audiobooks.category_name + "/" + audiobooks.save_name
	//	savePath := audiobooks.save_name
	f, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, response.Body)
	defer response.Body.Close()
	return string(savePath), nil
}

//使用http代理下载
func DownAudiobookByHttp(audiobooks Audiobooks) (string, error) {
	savePath := audiobooks.category_name + "/" + audiobooks.save_name
	exist, err := PathExists(savePath)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return string(savePath), nil
	}
	if exist {
		//		fmt.Printf("[%v] was downloaded\n", savePath)
		return string(savePath), nil
	}
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8888")
	//os.Setenv("HTTPS_PROXY", "https://127.0.0.1:9743")
	client := &http.Client{}
	url := audiobooks.down_url
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Cookie", "status_convertTime=00%3A00; statuscurrentTime=0.6635")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	reqest.Header.Add("Referer", audiobooks.url)
	if err != nil {
		panic(err)
	}
	res, _ := client.Do(reqest)
	//	savePath := audiobooks.category_name + "/" + audiobooks.save_name
	//	savePath := audiobooks.save_name
	f, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	defer res.Body.Close()
	return string(savePath), nil
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 建立文件夹
func Mkdir(dir string) {
	_dir := dir
	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if exist {
	} else {
		fmt.Printf("no dir![%v]\n", _dir)
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}

// 字符串转码
func DecodeToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GBK.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}
