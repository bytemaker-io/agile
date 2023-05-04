package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"sniffmac/statusinit"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://root:123456@192.168.1.196:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"macaddress",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	fmt.Println("Message sent to RabbitMQ!")

	statusinit.InitRouter()
	cmd := exec.Command("airodump-ng", "wlan1", "--essid-regex", "'  '")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	macRegexp := regexp.MustCompile(`(?:[0-9A-Fa-f]{2}[:.-]){5}[0-9A-Fa-f]{2}`)

	macAddresses := make(map[string]bool)
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		// 搜索 MAC 地址
		mac := macRegexp.FindString(line)
		if mac != "" && !macAddresses[mac] {
			macAddresses[mac] = true

			body := mac
			fmt.Println("Discovery a new MACaddress" + mac)
			err = ch.Publish(
				"",     //
				q.Name, //
				false,  //
				false,  //
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				},
			)
		}

	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}

}
