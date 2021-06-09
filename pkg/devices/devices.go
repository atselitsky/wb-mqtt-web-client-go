package devices

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type RS485 struct {
	Debug bool `json:"debug"`
	Ports []struct {
		Path    string `json:"path"`
		Devices []struct {
			SlaveID    int    `json:"slave_id"`
			DeviceType string `json:"device_type"`
		} `json:"devices"`
		PortType     string      `json:"port_type,omitempty"`
		BaudRate     int         `json:"baud_rate"`
		Parity       string      `json:"parity"`
		DataBits     int         `json:"data_bits"`
		StopBits     int         `json:"stop_bits"`
		PollInterval int         `json:"poll_interval"`
		Enabled      bool        `json:"enabled"`
		Type         interface{} `json:"type,omitempty"`
	} `json:"ports"`
}

func GetRS485Config(c *gin.Context) {
	file, err := os.Open("/etc/wb-mqtt-serial.conf")
	if err != nil {
		log.Fatalf("Cannot open RS-485 config file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.String(200, "%s", scanner.Text())
		// fmt.Fprintln(w, scanner.Text())
	}

}

func PostRS485Config(c *gin.Context) {
	// err := os.Remove("users.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// jsonFile, err := os.Create("users.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer jsonFile.Close()

	// scanner := bufio.NewScanner(c.Request.Body)
	// for scanner.Scan() {
	// 	//fmt.Println(scanner.Text())
	// 	//jsonFile.WriteString(scanner.Text())
	// 	fmt.Fprintln(jsonFile, scanner.Text())
	// }
	// //json.NewDecoder(c.Request.Body).Decode(encoder)
	// c.JSON(200, gin.H{
	// 	"status": "posted",
	// })
	// "/etc/wb-mqtt-serial.conf"
	var jsonFile RS485
	decoder := json.NewDecoder(c.Request.Body)
	decoder.Decode(&jsonFile)

	file, _ := json.MarshalIndent(jsonFile, "", " ")

	_ = ioutil.WriteFile("users.json", file, 0644)

	c.JSON(200, gin.H{
		"status": "posted",
	})
}
