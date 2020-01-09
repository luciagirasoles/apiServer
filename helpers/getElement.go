package helpers

import (
	models "apiServer/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//ElementInfo return title, icon url and server status
func ElementInfo(url string) (string, string, bool, error) {
	resp, err := http.Get(url)
	var title, iconURL string
	var zError error
	statusURL := false

	if err != nil {
		// Catch(err)
		return "", "", false, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if answer, ok := getHTMLTitle(resp.Body); ok {
		title = answer
	} else {
		println("Fail to get HTML title")
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if answer, ok := obtainIcon(resp.Body); ok {
		iconURL = answer
	} else {
		println("Fail to get HTML icon")
	}

	return title, iconURL, statusURL, zError
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"

}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func getHTMLTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func obtainIcon(r io.Reader) (string, bool) {

	var imgSrcURL, imgDataOriginal string

	if r != nil {
		log.Println("Page response is NOT nil")
		// --------------

		data, _ := ioutil.ReadAll(r)

		hdata := strings.Replace(string(data), "<noscript>", "", -1)
		hdata = strings.Replace(hdata, "</noscript>", "", -1)
		// --------------

		if document, err := html.Parse(strings.NewReader(hdata)); err == nil {
			var parser func(*html.Node)
			parser = func(n *html.Node) {

				if n.Type == html.ElementNode && n.Data == "link" {

					for _, element := range n.Attr {
						if element.Key == "href" {
							imgSrcURL = element.Val
						}
						if element.Key == "type" && element.Val == "image/x-icon" {
							imgDataOriginal = element.Val

						}
						if imgDataOriginal != "" && imgSrcURL != "" {
							break
						}

					}

				}

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parser(c)
					if imgDataOriginal != "" && imgSrcURL != "" {
						break
					}
				}

			}
			parser(document)
		} else {
			fmt.Println("Parse html error", err)
		}

	} else {
		fmt.Println("Page response IS nil")
	}
	return imgSrcURL, true
}

func getFromWHOIS(domain string) (string, string) {
	var country, owner string

	whoisRaw := runWhoisCommand(domain)
	reader := bytes.NewReader(whoisRaw.Bytes())
	scanner := bufio.NewScanner(reader)
	// Scan lines
	scanner.Split(bufio.ScanLines)

	// Scan through lines and find refer server
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "No match for domain ") {
			fmt.Println("No match domain")
		}
		if strings.Contains(line, "Country:") {
			// // Trim the refer: on left
			country = strings.TrimPrefix(line, "Country:")
			// // Trim whitespace
			country = strings.TrimSpace(country)

		}
		if strings.Contains(line, "OrgName: ") {
			owner = strings.TrimPrefix(line, "OrgName:")
			owner = strings.TrimSpace(owner)

		}

	}

	return country, owner
}
func defineSslGrade(listSsl []string) string {
	grade := ""
	sslRange := map[string]int{
		"A+": 8,
		"A":  7,
		"A-": 6,
		"B":  5,
		"C":  4,
		"D":  3,
		"E":  2,
		"F":  1,
	}
	min := 9
	if len(listSsl) != 0 {

		for _, v := range listSsl {
			if sslRange[v] < min {
				min = sslRange[v]
				grade = v
			}
		}
	}
	return grade
}

//SslInfo return list of srvers, grade of SSL, previos SL grade and if server is down
func SslInfo(urlAnalyze string) ([]models.Server, string, string, bool, error) {
	var result, resultPast1hour map[string]interface{}
	var sslChange bool
	mservers := []models.Server{}
	// mservers_past_1hour := []Server{}
	var grade, previousGrade string
	var err error

	urlSsl := fmt.Sprintf("https://api.dev.ssllabs.com/api/v3/analyze?host=%s", urlAnalyze)
	urlSslPast1hour := "https://api.dev.ssllabs.com/api/v3/analyze?host=" + urlAnalyze + "&maxAge=1"

	body, readErr := requestInfo(urlSsl)
	if readErr != nil {
		return mservers, grade, previousGrade, sslChange, readErr
	}
	json.Unmarshal(body, &result)

	endpoints, ok := result["endpoints"].([]interface{})
	if ok {
		mservers, grade = subServers(endpoints)
	}

	bodyPast1hour, readErr := requestInfo(urlSslPast1hour)
	if readErr != nil {
		return mservers, grade, previousGrade, sslChange, readErr
	}

	json.Unmarshal(bodyPast1hour, &resultPast1hour)
	endpoints, ok = result["endpoints"].([]interface{})

	if ok {
		_, previousGrade = subServers(endpoints)

	}
	c1, _ := Unmarshal(body)
	c2, _ := Unmarshal(bodyPast1hour)
	sslChange = !equal(c1, c2)

	return mservers, grade, previousGrade, sslChange, err

}
func requestInfo(url string) ([]byte, error) {
	var body []byte
	var err error
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return body, err
	}

	res, err := spaceClient.Do(req)
	if err != nil {
		return body, err
	}

	body, err = ioutil.ReadAll(res.Body)

	return body, err

}

func runWhoisCommand(args ...string) bytes.Buffer {
	// Store output on buffer
	var out bytes.Buffer

	// Execute command
	cmd := exec.Command("whois", args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()

	return out
}

func subServers(serversList []interface{}) ([]models.Server, string) {
	mservers := []models.Server{}
	var sslServers []string
	var grade string

	for _, value := range serversList {
		server := models.Server{}
		dataEndpoint := value.(map[string]interface{})

		serverAddress, ok := dataEndpoint["serverName"]
		server.Address = (map[bool]string{true: fmt.Sprintf("%v", serverAddress), false: ""})[ok]

		sslGrade, ok := dataEndpoint["grade"]
		server.SSLGrade = (map[bool]string{true: fmt.Sprintf("%v", sslGrade), false: ""})[ok]

		ipaddress, ok := dataEndpoint["ipAddress"]
		url := (map[bool]string{true: fmt.Sprintf("%v", ipaddress), false: ""})[ok]

		server.Country, server.Owner = getFromWHOIS(url)
		sslServers = append(sslServers, server.SSLGrade)
		grade = defineSslGrade(sslServers)
		mservers = append(mservers, server)

	}

	return mservers, grade
}

func equal(vx, vy interface{}) bool {
	if reflect.TypeOf(vx) != reflect.TypeOf(vy) {
		return false
	}

	switch x := vx.(type) {
	case map[string]interface{}:
		y := vy.(map[string]interface{})

		if len(x) != len(y) {
			return false
		}

		for k, v := range x {
			val2 := y[k]

			if (v == nil) != (val2 == nil) {
				return false
			}

			if !equal(v, val2) {
				return false
			}
		}

		return true
	case []interface{}:
		y := vy.([]interface{})

		if len(x) != len(y) {
			return false
		}

		var matches int
		flagged := make([]bool, len(y))
		for _, v := range x {
			for i, v2 := range y {
				if equal(v, v2) && !flagged[i] {
					matches++
					flagged[i] = true
					break
				}
			}
		}
		return matches == len(x)
	default:
		return vx == vy
	}
}

// Unmarshal parses the Body-encoded data into an interface{}
func Unmarshal(b []byte) (interface{}, error) {
	var j interface{}

	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return j, nil
}
