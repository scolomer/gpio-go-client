package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Connect struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

var gpioid string

func main() {
	log.Println("Starting")
	url := os.Args[1]
	value := -1

	id, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		panic(err)
	}

	var label = os.Args[3]
	gpioid = os.Args[4]

	for {
		log.Printf("Dialing to %v", url)
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Printf("Connected to %v", ws.RemoteAddr())

		conn := Connect{int(id), label, value}
		b, _ := json.Marshal(conn)
		if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}

		go func() {
			time.Sleep(10 * time.Second)
			ws.WriteMessage(websocket.TextMessage, []byte("{}"))
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Print(err)
				ws.Close()
				time.Sleep(2 * time.Second)
				break
			}
			log.Printf("Received: %s.\n", msg)
			var f interface{}
			err = json.Unmarshal(msg, &f)
			if err != nil {
				log.Print(err)
				ws.Close()
				time.Sleep(2 * time.Second)
				break
			}
			value = int(f.(map[string]interface{})["value"].(float64))
			switchGpio(value)
		}
	}
}

func switchGpio(value int) {
	export()

	err := ioutil.WriteFile("/sys/class/gpio/gpio"+gpioid+"/direction", []byte("out"), 0)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("/sys/class/gpio/gpio"+gpioid+"/value", []byte(strconv.Itoa(value)), 0)
	if err != nil {
		panic(err)
	}
}

/*func getGpioValue() int {
	export()

	err := ioutil.WriteFile("/sys/class/gpio/gpio17/direction", []byte("in"), 0)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile("/sys/class/gpio/gpio17/value")

	r, _ := strconv.Atoi(string(b))
	return r
}*/

/*func getGpioValue() int {
	x := []byte("2")
	z, _ :=strconv.Atoi(string(x))
	return z
}*/

func export() {
	_, err := os.Stat("/sys/class/gpio/gpio" + gpioid + "/value")
	if err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile("/sys/class/gpio/export", []byte(gpioid), 0)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
