package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"os"
	"crypto/tls"
	"strings"
	"runtime"
	"strconv"
	utilnet "k8s.io/apimachinery/pkg/util/net"
)


func do_request(ip string, host string, proto string, pattern string) {
	var url string 
	url = proto + "://" + ip + ":443/"
	resp, err := httpGetNoConnectionPoolTimeout(url, 5*time.Second)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyString := string(body)
		if strings.Contains(bodyString, pattern) {
			fmt.Printf("Found at %s\n", ip)
		} else {
			//fmt.Printf("Nothing at %s\n", ip)
		}

	}


}



func httpGetNoConnectionPoolTimeout(url string, timeout time.Duration) (*http.Response, error) {
	tr := utilnet.SetTransportDefaults(&http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse},
	}
	return client.Get(url)
}


func get_request(ip *string, host *string, proto *string, pattern *string, tr http.RoundTripper) {
	var url string

	if *proto == "http" {
		url = *proto + "://" + *ip + ":80/"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		req.Host = *host
		req.Header.Set("Connection", "Close")
		client := &http.Client{Timeout: time.Second * 4, CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse}, }
		res, err := client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		bodyString := string(body)
		if strings.Contains(bodyString, *pattern) {
			fmt.Printf("Found at %s\n", *ip)
		} else {
			//fmt.Printf("Nothing at %s\n", *ip)
		}

	} else {

		
		

		client := &http.Client{Timeout: time.Second * 8, Transport: tr, CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse}, }

		url = *proto + "://" + *ip + ":443/"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		req.Host = *host
		req.Header.Set("Connection", "close")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")



		res, err := client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		bodyString := string(body)
		//fmt.Printf("content at %s:\n", *ip)
		//fmt.Println(bodyString)
		if strings.Contains(bodyString, *pattern) {
			fmt.Printf("Found at %s\n", *ip)
		} else {
			//fmt.Printf("Nothing at %s\n", *ip)
		}

		//*client.CloseIdleConnections()


	}

	//fmt.Printf("sent request to %s\n", *ip)

}



func main() {

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s IPLIST HOST PROTO PATTERN START END\n\n", os.Args[0])
		return
	}
	host := os.Args[2]
	proto := os.Args[3]
	pattern := os.Args[4]
	starter, _ := strconv.Atoi(os.Args[5])
	ender, _ := strconv.Atoi(os.Args[6])
	runtime.GOMAXPROCS(15000)

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
    panic(err.Error())
	}
	lines := strings.Split(string(content), "\n")
	counter := 0

	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true,
    	}
	


  	i := 0
	for i = starter; i < ender; i++ {
		if len(lines[i]) > 0 {

			ip := lines[i]
			if counter > 2000 {
				time.Sleep(8 * time.Second)
				counter = 0
				fmt.Println("Waiting 8 secs...")
				}
			go get_request(&ip, &host, &proto, &pattern, tr)
			//go do_request(ip, host, proto, pattern)
			//fmt.Println(counter)
			counter++

		}
	}

	time.Sleep(8 * time.Second)

}
