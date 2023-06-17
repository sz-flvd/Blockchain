package client

import (
    "fmt"
    "net/http"
    "encoding/json"
	"bufio"
	"os"
	"strings"
	"io/ioutil"
	"bytes"
	"krypto.blockchain/src/api"
)

func HandleInput() {
	reader := bufio.NewReader(os.Stdin)


	for {
		fmt.Println("Give me a command and I'll do what you ask")

		cmd, _ := reader.ReadString('\n')

		switch {
		case strings.HasPrefix(cmd, "ls b"):
			getBlock()
		case strings.HasPrefix(cmd, "ls c"):
			getCount()
		case strings.HasPrefix(cmd, "ls l"):
			getLocalChain()
		case strings.HasPrefix(cmd, "add b"):
			addBlock()
		default:
			fmt.Println("Oh, mamma-mamma mia, command is not recognized!")
		}

	}
}

func getBlock() {
    var id string

	fmt.Println("Provide block id")

	fmt.Scanln(&id)
	
	response, err := http.Get("http://localhost:8080/blockchain/" + id)
    if err != nil {
        fmt.Println(err.Error())
		return
    }

	responseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}

func getCount() {
	response, err := http.Get("http://localhost:8080/blockchain/count")
    if err != nil {
        fmt.Println(err.Error())
		return
    }

	responseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Count = " + string(responseData))
}

func getLocalChain() {
	var id string

	fmt.Println("Provide node id")

	fmt.Scanln(&id)
	response, err := http.Get("http://localhost:8080/blockchain/local/" + id)
    if err != nil {
        fmt.Println(err.Error())
		return
    }

	responseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}

func addBlock() {
	var data string

	fmt.Println("Provide data for record")

	fmt.Scanln(&data)

	// for now (and possibly for ever) we send requests to all nodes
	request := api.AddRecordRequest{Content: data, Receivers: []int{}}
	json, _ := json.Marshal(request)

	response, err := http.Post("http://localhost:8080/blockchain", "application/json", bytes.NewBuffer(json))
    if err != nil {
        fmt.Println("Ooops..." + err.Error())
		return
    }

	responseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}