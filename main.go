package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response_403 struct {
	IsSuccess  bool       `json:"isSuccess"`
	Message    string     `json:"message"`
	Result     Result_403 `json:"result"`
	StatusCode uint8      `json:"statusCode"`
}

type Result_403 struct {
	SanctionStatus bool           `json:"sanction_status"`
	Support        bool           `json:"support"`
	Crossings      []Crossing_403 `json:"crossings"`
}

type Crossing_403 struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// {
// 	"isSuccess": true,
// 	"message": "result success",
// 	"result": {
// 	  "support": true,
// 	  "sanction_status": true,
// 	  "crossings": [
// 		{
// 		  "name": "DNS",
// 		  "link": ""
// 		},
// 		{
// 		  "name": "سرویس ویژه",
// 		  "link": ""
// 		}
// 	  ]
// 	},
// 	"statusCode": 200
//   }

type Config struct {
	Verbose bool
}

var config Config

func Prin403tResponse(resp Response_403) int {

	if resp.IsSuccess == false {
		fmt.Println(resp.Message)
		return 1
	}

	fmt.Println("sanction status : ", resp.Result.SanctionStatus)

	if config.Verbose == false {
		return 0
	}

	fmt.Println("403 support status : ", resp.Result.Support)
	fmt.Println("supporteed methods : ")
	for _, crossing := range resp.Result.Crossings {
		if crossing.Name != "سرویس ویژه" {
			fmt.Println("\t", crossing.Name)
		} else {
			fmt.Println("\t VIP service (free)")
		}
		fmt.Println("\t", crossing.Link)
	}

	return 0
}

func HandleArgs(args []string) []string {
	var uris []string

	for _, arg := range args[1:] {
		switch arg {
		case "-v":
			config.Verbose = true
		default:
			uris = append(uris, arg)
		}

	}
	return uris
}

func GetDomainInfo(uri403 string) (Response_403, error) {
	resp, err := http.Get("https://api.anti403.ir/api/search-filter?url=" + uri403)

	if err != nil {
		return Response_403{}, err
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return Response_403{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Response_403{}, errors.New("request status code not ok")
	}

	decoded := Response_403{}
	json.Unmarshal([]byte(string(body)), &decoded)
	return decoded, nil
}

func main() {

	uris := HandleArgs(os.Args)

	var response403 Response_403
	var err error

	for _, uri := range uris {
		fmt.Println("requesting for ", uri)
		response403, err = GetDomainInfo(uri)

		if err != nil {
			fmt.Println(err)
			continue
		}

		Prin403tResponse(response403)
		if config.Verbose {
			fmt.Println("\n\n\n")
		}
	}
}
