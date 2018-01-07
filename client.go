package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"os"
  "encoding/json"
  "strconv"
	"time"
	"net"
)

type Connect struct {
    Id int `json:"id"`
    Description string `json:"description"`
		Value int `json:"value"`
}

func main() {
	fmt.Println("Connect")
	url := os.Args[1]
	value := -1
	timeout := 5 * time.Second

	id, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		panic(err)
	}

	var label string
	if len(os.Args) >= 3 {
		label = os.Args[3]
	} else {
		label = "Inconnu"
	}

	for {
		log.Print("Dial")
		ws, err := websocket.Dial(url, "", "http://192.168.0.6:9000")
		if err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}
		conn := Connect{int(id), label, value}
		b, _ := json.Marshal(conn)
		if _, err := ws.Write(b); err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}
		var msg = make([]byte, 512)
		var n int
	  for {
			ws.SetReadDeadline(time.Now().Add(timeout))
		  if n, err = ws.Read(msg); err != nil {
				if e,ok := err.(net.Error); ok && e.Timeout() {
					ws.Write([]byte("{}"))
					continue
				}

				log.Print(err)
				ws.Close()
				time.Sleep(2 * time.Second)
				break
		  }
		  fmt.Printf("Received: %s.\n", msg[:n])
	    var f interface{}
	    err := json.Unmarshal(msg[:n], &f)
	    if err != nil {
				log.Print(err)
				ws.Close()
				time.Sleep(2 * time.Second)
				break
	    }
	    value = int(f.(map[string]interface{})["value"].(float64))
	    switchGpio(value);
	  }
	}
}

func switchGpio(value int) {
	export()

	err := ioutil.WriteFile("/sys/class/gpio/gpio475/direction", []byte("out"), 0)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("/sys/class/gpio/gpio475/value", []byte(strconv.Itoa(value)), 0)
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
	_, err := os.Stat("/sys/class/gpio/gpio475/value")
	if err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile("/sys/class/gpio/export", []byte("475"), 0)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
