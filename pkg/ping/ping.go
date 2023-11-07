package ping

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ewilan-riviere/spotlight/pkg/terminal"
)

type WebsiteHealth struct {
	Domain       string
	Http         string
	Https        string
	UseHttp2     bool
	UseHttps     bool
	HttpRedirect bool
	HttpCode     int
	PingSuccess  bool
	CurlSuccess  bool
	Server       string
	IpAdress     string
	Time         string
	Ok           bool
}

type CurlOutput struct {
	HttpCode     int
	HttpRedirect bool
	UseHttp2     bool
	Server       string
}

func Make(domains []string) []WebsiteHealth {
	var websites []WebsiteHealth

	for _, domain := range domains {
		websites = append(websites, WebsiteHealth{
			Domain: domain,
			Http:   "http://" + domain,
			Https:  "https://" + domain,
		})
	}

	websites = ping(websites)
	websites = curl(websites)

	return websites

}

func ping(websites []WebsiteHealth) []WebsiteHealth {
	items := []WebsiteHealth{}
	for _, website := range websites {
		output := terminal.ExecCommand("ping -c 5 " + website.Domain)

		statistics := strings.Split(output, "\n")
		statisticsLine := statistics[len(statistics)-3] // 5 packets transmitted, 5 packets received, 0.0% packet loss

		statisticsLine = strings.Replace(statisticsLine, " packets transmitted", "", -1)
		statisticsLine = strings.Replace(statisticsLine, " packets received", "", -1)
		statisticsLine = strings.Replace(statisticsLine, "% packet loss", "", -1)
		statList := strings.Split(statisticsLine, ", ")

		statMap := map[string]string{
			"transmitted": statList[0],
			"received":    statList[1],
			"loss":        statList[2],
		}
		loss := statMap["loss"]
		lossFloat, err := strconv.ParseFloat(loss, 64)
		if err != nil {
			fmt.Println(err)
		}

		firstPingLine := statistics[1] // 64 bytes from IP_ADDRESS: icmp_seq=0 ttl=54 time=32.918 ms
		ipAddressSplit := strings.Split(firstPingLine, "bytes from")
		ipAddressSplitLastPart := ipAddressSplit[1]
		ipAddressSplit2 := strings.Split(ipAddressSplitLastPart, ": icmp_seq=")
		ipAddressFirstPart := ipAddressSplit2[0]
		ipAddress := strings.TrimSpace(ipAddressFirstPart)

		timeSplit := strings.Split(firstPingLine, "time=") // 64 bytes from IP_ADDRESS: icmp_seq=0 ttl=54 time=32.918 ms
		time := timeSplit[1]

		website.PingSuccess = lossFloat < 50
		website.IpAdress = ipAddress
		website.Time = time

		items = append(items, website)
	}

	return items
}

func curl(websites []WebsiteHealth) []WebsiteHealth {
	items := []WebsiteHealth{}
	for _, website := range websites {
		outputHttp := terminal.ExecCommand("curl -I " + website.Http)
		curlHttp := parseCurl(outputHttp)
		outputHttps := terminal.ExecCommand("curl -I " + website.Https)
		curlHttps := parseCurl(outputHttps)

		if curlHttps.HttpCode == 200 {
			website.CurlSuccess = curlHttps.HttpCode == 200
			website.HttpCode = curlHttps.HttpCode
			website.UseHttp2 = curlHttps.UseHttp2
			website.UseHttps = true
			website.HttpRedirect = curlHttp.HttpRedirect
			website.Server = curlHttps.Server
			website.Ok = true
		} else {
			website.CurlSuccess = curlHttp.HttpCode == 200
			website.HttpCode = curlHttp.HttpCode
			website.UseHttp2 = curlHttp.UseHttp2
			website.UseHttps = false
			website.HttpRedirect = curlHttp.HttpRedirect
			website.Server = curlHttp.Server
			website.Ok = false
		}

		items = append(items, website)
	}

	return items
}

func parseCurl(output string) CurlOutput {
	split := strings.Split(output, "\n")
	httpCode := strings.Split(split[0], " ")[1]

	httpCodeI, err := strconv.Atoi(httpCode)
	if err != nil {
		fmt.Println(err)
	}

	server := strings.Split(split[1], " ")[1]

	return CurlOutput{
		HttpCode:     httpCodeI,
		HttpRedirect: httpCodeI == 301 || httpCodeI == 302,
		UseHttp2:     strings.Contains(output, "HTTP/2"),
		Server:       server,
	}
}

func (w WebsiteHealth) ToString() string {
	return fmt.Sprintf("%+v", w)
}
