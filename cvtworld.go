package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Exit struct {
	Look        string `json:"look"`
	Flags       string `json:"flags"`
	Destination int    `json:"destination"`
}
type Room struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Flags       string        `json:"flags"`
	Exits       map[int]*Exit `json:"exits"`
}

func readFile(path string) (string, error) {
	/*
	   file, err := os.Open(path)
	   if err != nil {
	     return nil, err
	   }
	   defer file.Close()

	   var lines []string
	   scanner := bufio.NewScanner(file)
	   for scanner.Scan() {
	     lines = append(lines, scanner.Text())
	   }
	   return lines, scanner.Err()
	*/
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		return "", err
	}

	str := string(b) // convert content to a 'string'

	return str, nil
}

func buildWorld() (map[int]*Room, error) {
	strWorld, err := readFile("./world.wld")

	if err != nil {
		fmt.Println("Error loading world file: " + err.Error())
		return nil, err
	}

	//	fmt.Println(strWorld)

	world := make(map[int]*Room)

	rooms := strings.Split(strWorld, "#")
	roomNum := 0
	split := 0
	tmp := ""

	// Build room by room
	for i := range rooms {
		if i != 0 {
			r := new(Room)
			split = strings.Index(rooms[i], "\n")
			roomNum, err = strconv.Atoi(rooms[i][:split])
			//			fmt.Println("Room number: " + strconv.Itoa(roomNum))
			tmp = rooms[i][split+1:]
			split = strings.Index(tmp, "~")
			//			fmt.Println("Room name: " + tmp[:split])
			r.Name = tmp[:split]
			fmt.Println("Room name: [" + r.Name + "]")
			tmp = tmp[split+2:]
			split = strings.Index(tmp, "~")
			//			fmt.Println("Room desc: " + tmp[:split])
			r.Description = tmp[:split-1]
			fmt.Println("Room desc: [" + r.Description + "]")
			tmp = tmp[split+1:]
			// Find out what this one does
			split = strings.Index(tmp, "~")
			split = strings.Index(tmp, "~")
			//			fmt.Println("Room flags: " + tmp[:split])
			r.Flags = tmp[:split+8]
			fmt.Println("Room flags: [" + r.Flags + "]")
			split = strings.Index(tmp, "~")

			fmt.Println(strconv.Itoa(roomNum))
			fmt.Println(r)
		}
	}

	return world, nil
}

func main() {

	// Build the world
	world, err := buildWorld()
	if err != nil {
		fmt.Println("Error building world")
		return
	}

	fmt.Println("World has been built (" + strconv.Itoa(len(world)) + " rooms)")
}

//!-main
