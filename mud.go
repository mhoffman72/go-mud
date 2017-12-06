package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

//!+gameloop
type client chan<- string // an outgoing message channel

type Direction int

const (
	North Direction = iota
	South
	East
	West
	Up
	Down
)

type Character struct {
	Name       string
	Room       int
	RemoteAddr string
}
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
	Players     map[string]*Character
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func gameloop() {
	clients := make(map[client]bool) // all connected clients
	//	players := make(map[client]*Character) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-gameloop

//!+handleConn
func handleConn(conn net.Conn, world map[int]*Room) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "\nWelcome to goMud!\n\nLogin: "
	//	messages <- "New connection from " + who
	entering <- ch

	input := bufio.NewScanner(conn)

	input.Scan()
	username := input.Text()

	ch <- "Password: "
	input.Scan()
	password := input.Text()

	var player Character
	player.Name = username
	player.Room = 3001
	player.RemoteAddr = who
	//world[player.Room].Players[player.Name] = &player

	ch <- "Welcome back " + player.Name + " (password: " + password + ")\n\n"
	//	entering <- &player

	//messages <- player.Name + " has entered the game.\n"

	msg := ""
	done := false
	cmd := ""

	ch <- world[player.Room].Name + "\n\n" + world[player.Room].Description + "\n> "

	//	input := bufio.NewScanner(conn)
	for input.Scan() {

		//fmt.Println("Input = " + input.Text())
		if strings.Index(input.Text(), " ") != -1 {
			parts := strings.Split(input.Text(), " ")
			cmd = parts[0]
		} else {
			cmd = input.Text()
		}

		var buf string

		//fmt.Println("Command = " + cmd)
		switch cmd {
		case "bye":
			msg = "Goodbye, " + player.Name + "!\n"
			done = true
		case "l":
			fallthrough
		case "look":
			ch <- world[player.Room].Name + "\n\n" + world[player.Room].Description + "\n"
		case "exits":
			for k, v := range world[player.Room].Exits {
				switch k {
				case 0:
					_, ok := world[v.Destination]
					if ok {
						buf += "North: " + world[v.Destination].Name + "\n"
					} else {
						buf += "North: " + strconv.Itoa(v.Destination) + "\n"
					}
				case 1:
					_, ok := world[v.Destination]
					if ok {
						buf += "East: " + world[v.Destination].Name + "\n"
					} else {
						buf += "East: " + strconv.Itoa(v.Destination) + "\n"
					}
				case 2:
					_, ok := world[v.Destination]
					if ok {
						buf += "South: " + world[v.Destination].Name + "\n"
					} else {
						buf += "South: " + strconv.Itoa(v.Destination) + "\n"
					}
				case 3:
					_, ok := world[v.Destination]
					if ok {
						buf += "West: " + world[v.Destination].Name + "\n"
					} else {
						buf += "West: " + strconv.Itoa(v.Destination) + "\n"
					}
				case 4:
					_, ok := world[v.Destination]
					if ok {
						buf += "Up: " + world[v.Destination].Name + "\n"
					} else {
						buf += "Up: " + strconv.Itoa(v.Destination) + "\n"
					}
				case 5:
					_, ok := world[v.Destination]
					if ok {
						buf += "Down: " + world[v.Destination].Name + "\n"
					} else {
						buf += "Down: " + strconv.Itoa(v.Destination) + "\n"
					}
				}
			}
			ch <- buf
		case "n":
			_, ok := world[player.Room].Exits[0]
			if ok {
				_, ok = world[world[player.Room].Exits[0].Destination]
				if ok {
					player.Room = world[player.Room].Exits[0].Destination
				} else {
					ch <- "\n***Room not defined in world***\n\n"
				}
				ch <- world[player.Room].Name + "\n\n" + world[player.Room].Description + "\n"
			} else {
				ch <- "You can't go that way\n"
			}
		case "e":
			_, ok := world[player.Room].Exits[1]
			if ok {
				_, ok = world[world[player.Room].Exits[1].Destination]
				if ok {
					//					delete(world[player.Room].Players, player.Name)
					player.Room = world[player.Room].Exits[1].Destination
					//					world[player.Room].Players[player.Name] = &player
				} else {
					ch <- "\n***Room not defined in world***\n\n"
				}
				ch <- world[player.Room].Name + "\n\n" + world[player.Room].Description + "\n"
			} else {
				ch <- "You can't go that way\n"
			}
		case "s":
			_, ok := world[player.Room].Exits[2]
			if ok {
				_, ok = world[world[player.Room].Exits[2].Destination]
				if ok {
					player.Room = world[player.Room].Exits[2].Destination
				} else {
					ch <- "\n***Room not defined in world***\n\n"
				}
				ch <- world[player.Room].Name + "\n\n" + world[player.Room].Description + "\n"
			} else {
				ch <- "You can't go that way\n"
			}
		case "w":
			_, ok := world[player.Room].Exits[3]
			if ok {
				_, ok = world[world[player.Room].Exits[3].Destination]
				if ok {
					//					delete(world[player.Room].Players, player.Name)
					player.Room = world[player.Room].Exits[3].Destination
					//					world[player.Room].Players[player.Name] = &player
				} else {
					ch <- "\n***Room not defined in world***\n\n"
				}
				ch <- world[player.Room].Name + "\n\n" + world[player.Room].Description + "\n"
			} else {
				ch <- "You can't go that way\n"
			}
		case "":
			msg = "> "
		default:
			msg = "I don't understant that command\n"
		}

		//		if strings.EqualFold(parts[0], "bye") {
		//			break
		//		}
		//
		//		messages <- who + ": " + input.Text()
		ch <- msg

		if done {
			time.Sleep(1000 * time.Millisecond)
			break
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprint(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

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
	/*
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
				//fmt.Println(rooms[i])
				split = strings.Index(rooms[i], "\n")
				//fmt.Println(strconv.Itoa(split))
				//tmp = rooms[i]
				//fmt.Println(tmp[:split])
				roomNum, err = strconv.Atoi(rooms[i][:split])
				//fmt.Println(rooms[i][:split])
				fmt.Println("Room number: " + strconv.Itoa(roomNum))
				tmp = rooms[i][split+1:]
				world[roomNum] = new(Room)
				split = strings.Index(tmp, "~")
				world[roomNum].Name = tmp[:split]
				world[roomNum].Num = roomNum
				tmp = tmp[split+1:]
				split = strings.Index(tmp, "~")
				world[roomNum].Description = tmp[:split]
			}
		}
	*/
	world := make(map[int]*Room)

	worldFile, err := ioutil.ReadFile("./world.json")

	//n := bytes.Index(yamlFile, []byte{0})
	//log.Printf(string(yamlFile[:n]))

	if err != nil {
		log.Printf("worldFile.Get err   #%v ", err)
	}
	//	err = yaml.Unmarshal(yamlFile, world)
	err = json.Unmarshal(worldFile, &world)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return world, nil
}

//!+main
func main() {

	// Build the world
	world, err := buildWorld()

	if err != nil {
		fmt.Println("Error building world")
		return
	}

	//	world := make(map[int]*Room)
	//	world[1000] = new(Room)
	//	world[1000].Name = "The Temple of Midgaard"
	//	world[1000].Num = 1000
	//	world[1000].Exits[East] = 1001
	//	world[1001] = new(Room)
	//	world[1001].Name = "The Village Square"
	//	world[1001].Num = 1001
	//	world[1001].Exits[West] = 1000

	fmt.Println("World has been built (" + strconv.Itoa(len(world)) + " rooms)")

	fmt.Println(world)

	for k, v := range world {
		fmt.Printf("key[%s] value[%s]\n", k, v)

		for k2, v2 := range v.Exits {
			fmt.Printf("   key[%s] value[%s]\n", k2, v2)
		}
	}

	// Open the socket
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on: " + listener.Addr().String())

	// Start the Game Loop
	go gameloop()

	// Wait for players to connect
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, world)
	}
}

//!-main
