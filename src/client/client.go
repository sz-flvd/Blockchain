package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
		// Cryptocurrency
		case strings.HasPrefix(cmd, "ls w"):
			getBalance()
		case strings.HasPrefix(cmd, "add t"):
			addTransfer()
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
	var data []string

	fmt.Println("Provide data for record")

	for {
		singleTransaction := ""
		fmt.Scanln(&singleTransaction)

		if singleTransaction == "" {
			break
		}

		data = append(data, singleTransaction)

	}

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

// Cryptocurrency
func getBalance() {
	var id string

	fmt.Println("Provide node id")

	fmt.Scanln(&id)
	response, err := http.Get("http://localhost:8080/blockchain/local/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	responseNode, err := http.Get("http://localhost:8080/blockchain/nodes/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	responseData, _ := ioutil.ReadAll(response.Body)
	blockchainStr := string(responseData)

	responseNodeData, _ := ioutil.ReadAll(responseNode.Body)
	// publicKeyStr := responseNodeData.publicKey.String()
	fmt.Println(string(responseNodeData) + "XXDDD")

	lines := strings.Split(blockchainStr, "\n")

	contentBeginIdx := 0
	// contentEndIdx := 0
	for lineIdx, line := range lines {
		// modifiedLine := line + "xD"
		// fmt.Println(modifiedLine)
		if strings.Contains(line, "\"Records\": [") {
			contentBeginIdx = lineIdx
			myIdx := contentBeginIdx + 1
			for {
				if strings.Contains(lines[myIdx], "]") {
					// contentEndIdx = myIdx
					break
				}
				myIdx++
			}
		}
	}
	// fmt.Println(string(responseData))
}

func addTransfer() {
	var senderId string
	var receiverId string
	var amount int
	var data []string

	fmt.Println("Provide sender id")

	fmt.Scanln(&senderId)

	fmt.Println("Provide receiver id")

	fmt.Scanln(&receiverId)

	fmt.Println("Provide amount of money to transfer")

	fmt.Scanf("%d", amount)

	for {
		singleTransaction := ""
		fmt.Scanln(&singleTransaction)

		if singleTransaction == "" {
			break
		}

		data = append(data, singleTransaction)

	}

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
