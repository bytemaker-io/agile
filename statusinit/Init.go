package statusinit

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

/*
*
Author:kalean
Time:2023/05/03
*
*/

const InterfaceName string = "wlan1"

/**
 * This function will check the monitor mode is on or not
**/
func InitRouter() bool {
	modetype := getInterfacemode()
	if modetype == "monitor" {
		log.Println("the interface is in monitor mode")
		status := getInterfaceStatus()
		if status == "UP" {
			log.Println("the interface is up")
			log.Println("init finished")
			return true
		}
		log.Println("the interface is not up")
		//start the interface
		log.Println("start the interface")
		result := openInterface()
		if result == true {
			log.Println("the interface is up")
			log.Println("init finished")
			return true
		} else {
			log.Println("the interface is not up")
			log.Println("init failed")
			return false
		}
	}

	log.Println("the interface is not in monitor mode")
	//close the interface
	log.Println("close the interface")
	result := closeInterface()
	if result == true {
		log.Println("the interface is closed")
		//set the interface to monitor mode
		log.Println("set the interface to monitor mode")
		result = setInterfacemode()
		if result == true {
			log.Println("the interface is in monitor mode")
			//start the interface
			log.Println("start the interface")
			result = openInterface()
			if result == true {
				log.Println("the interface is up")
				log.Println("init finished")
				return true
			} else {
				log.Println("the interface is not up")
				log.Println("init failed")
				return false
			}
		} else {
			log.Println("init failed")
			return false
		}

	}
	log.Println("the interface is not in monitor mode")
	log.Println("init failed")
	return false
}

/**
check the interface is up or down
**/

func setInterfacemode() bool {
	//set the interface to monitor mode
	cmd := exec.Command("iw", "dev", InterfaceName, "set", "type", "monitor")
	_, err := cmd.Output()
	if err != nil {
		log.Panic("command execution failed:", err)
		return false
	}
	result := getInterfacemode()

	if strings.Contains(result, "type monitor") {
		log.Println("the interface is in monitor mode")
		return true
	} else {
		log.Println("the interface is not in monitor mode")
		return false
	}
}

func getInterfacemode() string {
	log.Println("check the interface mode")
	cmd := exec.Command("iw", "dev", InterfaceName, "info")
	output, err := cmd.Output()

	// 处理错误
	if err != nil {
		log.Panic("command execution failed:", err)
	}

	// 将输出转换为字符串
	outputStr := string(output)

	fields := strings.Fields(outputStr)
	typeIdx := indexOf(fields, "type")
	if typeIdx == -1 || typeIdx+1 >= len(fields) {
		log.Panic("Failed to get type")

	}
	typeStr := fields[typeIdx+1]
	if match, _ := regexp.MatchString("^[a-zA-Z]+$", typeStr); !match {
		log.Panic("Failed to get type")
	}
	return typeStr
}

func getInterfaceStatus() string {
	//check the interface is up or down
	cmd := exec.Command("ip", "a", "show", "dev", InterfaceName)
	output, err := cmd.Output()
	if err != nil {
		log.Panic("command execution failed:", err)
	}
	outputStr := string(output)
	fields := strings.Fields(outputStr)  // 按空格分割字符串
	stateIdx := indexOf(fields, "state") // 获取 "state" 的位置
	if stateIdx == -1 || stateIdx+1 >= len(fields) {
		log.Panic("Failed to get state")
	}
	state := fields[stateIdx+1]
	if match, _ := regexp.MatchString("^[a-zA-Z]+$", state); !match {
		log.Panic("Failed to get state")
	}
	return state
}

func indexOf(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

/**
*	This function will close the interface
**/
func closeInterface() bool {
	result := getInterfaceStatus()
	if result != "DOWN" {
		//close the interface
		cmd := exec.Command("ip", "link", "set", InterfaceName, "down")
		_, err := cmd.Output()
		if err != nil {
			log.Panic("command execution failed:", err)
			return false
		}
		result := getInterfaceStatus()
		log.Println("the interface status is " + result)
		if strings.Contains(result, "DOWN") {
			log.Println("the interface is down")
			return true
		}
		log.Panic("the interface is still up,plesae check the interface")
	}
	return true
}

/*
*
/This function will open the interface
*
*/
func openInterface() bool {
	cmd := exec.Command("ip", "a", "show", "dev", InterfaceName)
	_, err := cmd.Output()
	if err != nil {
		log.Panic("command execution failed:", err)
		return false
	}

	result := getInterfaceStatus()

	if strings.Contains(result, "DOWN") {

		cmd = exec.Command("ip", "link", "set", InterfaceName, "up")
		_, err = cmd.Output()
		if err != nil {
			log.Panic("command execution failed:", err)
			return false
		}
		result := getInterfaceStatus()
		if strings.Contains(result, "UP") {
			log.Println("the interface is Up")
			return true
		}
		log.Panic("the interface is still down,plesae check the interface")
	}
	log.Println("unknown reason can't start the interface")
	return false
}
