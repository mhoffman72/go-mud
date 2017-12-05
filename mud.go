// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
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
	NorthEast
	NorthWest
	SouthEast
	SouthWest
	Up
	Down
)

type Character struct {
	Name		string
	Room		int
	RemoteAddr	string
}

type Room struct {
	Name		string
	Num		int
	Description	string
	Exits		[10]int // n, s, e, w, ne, nw, se, sw, u, d
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
	ch <- "Welcome to goMud!\n\nLogin: " 
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

	ch <- "Welcome back " + player.Name + " (password: " + password + ")\n" 
//	entering <- &player

	messages <- player.Name + " has entered the game.\n"

	msg := ""
	done := false
	cmd := ""

	ch <- world[player.Room].Name + "\n" + world[player.Room].Description + "\n" + "> "

//	input := bufio.NewScanner(conn)
	for input.Scan() {

		//fmt.Println("Input = " + input.Text())
		if strings.Index(input.Text(), " ") != -1 {
			parts := strings.Split(input.Text(), " ")
			cmd = parts[0]
		} else {
			cmd = input.Text()
		}
		//fmt.Println("Command = " + cmd)
		switch cmd {
		case "bye":
			msg = "Goodbye, " + player.Name + "!\n"
			done = true
		case "l":
		case "look":
			ch <- world[player.Room].Name + "\n" + world[player.Room].Description + "\n" + "> "
		case "e":
			if world[player.Room].Exits[East] != 0 {
				player.Room = world[player.Room].Exits[East]
				ch <- world[player.Room].Name + "\n"
			} else {
				ch <- "You can't go that way\n> "
			}
//			player.Room++
//			msg = "You are in room " + strconv.Itoa(player.Room) + "\n"
		case "w":
			if world[player.Room].Exits[West] != 0 {
				player.Room = world[player.Room].Exits[West]
				ch <- world[player.Room].Name + "\n"
			} else {
				ch <- "You can't go that way\n> "
			}
//			player.Room--
//			msg = "You are in room " + strconv.Itoa(player.Room) + "\n"
		case "":
			msg = "> "
		default:
			msg = "I don't understant that command\n> "
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

	// Open the socket
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

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
