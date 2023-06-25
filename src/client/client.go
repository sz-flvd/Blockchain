package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	responseNode, err := http.Get("http://localhost:8080/blockchain/public/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	responseData, _ := ioutil.ReadAll(response.Body)
	blockchainStr := string(responseData)

	responseKeyData, _ := ioutil.ReadAll(responseNode.Body)
	publicKeyStr := string(responseKeyData)
	// fmt.Println(string(responseKeyData) + "XXDDD")

	lines := strings.Split(blockchainStr, "\n")

	walletValue := 0
	for lineIdx, line := range lines {
		if strings.Contains(line, "\"Records\": [") {
			myIdx := lineIdx
			for {
				if strings.Contains(lines[myIdx], "]") {
					break
				}
				recordLine := strings.Split(lines[myIdx], "$")
				if len(recordLine) == 5 && strings.Contains(recordLine[0], "Content") {
					// If he's a receiver
					if strings.Compare(recordLine[2], publicKeyStr) == 0 {
						value, err := strconv.Atoi(recordLine[3])
						if err == nil {
							// fmt.Println("!!!1")
							walletValue += value
						}
					}
					// If he's a sender
					if strings.Compare(recordLine[1], publicKeyStr) == 0 {
						value, err := strconv.Atoi(recordLine[3])
						if err == nil {
							// fmt.Println("!!!2")
							walletValue -= value
						}
					}
				}
				myIdx++
			}
		}
	}
	// fmt.Println("!!!3")
	fmt.Print("Balance: ")
	fmt.Println(walletValue)
}

func addTransfer() {
	var senderId string
	var receiverId string
	var amount string
	var data []string

	fmt.Println("Provide sender id")

	fmt.Scanln(&senderId)

	fmt.Println("Provide receiver id")

	fmt.Scanln(&receiverId)

	fmt.Println("Provide amount of money to transfer")

	fmt.Scanln(&amount)

	senderPublicKeyResponse, err := http.Get("http://localhost:8080/blockchain/public/" + senderId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	receiverPublicKeyResponse, err := http.Get("http://localhost:8080/blockchain/public/" + receiverId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	senderPublicKey, _ := ioutil.ReadAll(senderPublicKeyResponse.Body)
	receiverPublicKey, _ := ioutil.ReadAll(receiverPublicKeyResponse.Body)
	recordInfo := "$" + string(senderPublicKey) + "$" + string(receiverPublicKey) + "$" + string(amount) + "$"
	data = append(data, recordInfo)

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
