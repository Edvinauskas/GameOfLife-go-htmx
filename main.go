package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const CellWidth = 30
const CellHeight = 30
const AliveCSSClass = "alive"
const DeadCSSClass = "dead"
const TableCell = `<td ws-send hx-vals='{"event":"%s", "cellId":"%s"}' hx-trigger="click" id="%s" class="cell %s"></td>`
const StartStopButton = `<div id="start-stop-button" class="button is-%s" ws-send hx-vals='{"event":"%s"}'>%s</div>`

var AliveCells [][]int

var CellsHorizontally = 0
var CellsVertically = 0

var GameRunning = false

var Upgrader = websocket.Upgrader{}

type Event int

type GameJSON struct {
	Event        string  `json:"event"`
	CellID       string  `json:"cellId"`
	ScreenWidth  float64 `json:"screenWidth"`
	ScreenHeight float64 `json:"screenHeight"`
}

const (
	GenerateGrid Event = iota
	Start
	Stop
	AddCell
	RemoveCell
)

func (e Event) String() string {
	switch e {
	case GenerateGrid:
		return "GenerateGrid"
	case Start:
		return "Start"
	case Stop:
		return "Stop"
	case AddCell:
		return "AddCell"
	case RemoveCell:
		return "RemoveCell"
	default:
		return "unknown"
	}
}

func setCellAlive(cellId string) {
	x, y := parseCellId(cellId)
	AliveCells[x][y] = 1
}

func setCellDead(cellId string) {
	x, y := parseCellId(cellId)
	AliveCells[x][y] = 0
}

func parseCellId(cellId string) (int, int) {
	cellLocation := strings.Split(cellId, "_")
	x, _ := strconv.Atoi(cellLocation[0])
	y, _ := strconv.Atoi(cellLocation[1])
	return x, y
}

func generateCells(width float64, height float64) string {
	CellsHorizontally = int(math.Ceil(width / CellWidth))
	CellsVertically = int(math.Ceil(height / CellHeight))

	fmt.Printf("Cell v: %v - Cell h: %v", CellsVertically, CellsHorizontally)

	AliveCells = make([][]int, CellsVertically)
	for i := 0; i < CellsVertically; i++ {
		AliveCells[i] = make([]int, CellsHorizontally)
	}

	var tableHTML = fmt.Sprintf(`<table id="game-grid" style="height:%vpx;width:%vpx;">`, height, width)
	for v := 0; v < int(CellsVertically); v++ {
		tableHTML += "<tr>"
		for h := 0; h < int(CellsHorizontally); h++ {
			var cellId = strconv.Itoa(v) + "_" + strconv.Itoa(h)
			tableHTML += fmt.Sprintf(TableCell, AddCell.String(), cellId, cellId, DeadCSSClass)
		}
		tableHTML += "</tr>"
	}
	tableHTML += "</table>"
	return tableHTML
}

func life() string {
	var returnCells = ""
	cellsAliveCopy := make([][]int, len(AliveCells))
	for i := range AliveCells {
		cellsAliveCopy[i] = make([]int, len(AliveCells[i]))
		copy(cellsAliveCopy[i], AliveCells[i])
	}
	for rowIndex, row := range AliveCells {
		for colIndex, value := range row {
			neighborsAlive := 0
			neighborsDead := 0
			//North
			if rowIndex-1 >= 0 {
				if AliveCells[rowIndex-1][colIndex] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//North East
			if rowIndex-1 >= 0 && colIndex+1 < CellsHorizontally {
				if AliveCells[rowIndex-1][colIndex+1] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//East
			if colIndex+1 < CellsHorizontally {
				if AliveCells[rowIndex][colIndex+1] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//South East
			if colIndex+1 < CellsHorizontally && rowIndex+1 < CellsVertically {
				if AliveCells[rowIndex+1][colIndex+1] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//South
			if rowIndex+1 < CellsVertically {
				if AliveCells[rowIndex+1][colIndex] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//South West
			if rowIndex+1 < CellsVertically && colIndex-1 >= 0 {
				if AliveCells[rowIndex+1][colIndex-1] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//West
			if colIndex-1 >= 0 {
				if AliveCells[rowIndex][colIndex-1] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}
			//North West
			if rowIndex-1 >= 0 && colIndex-1 >= 0 {
				if AliveCells[rowIndex-1][colIndex-1] == 1 {
					neighborsAlive += 1
				} else {
					neighborsDead += 1
				}
			}

			//any live cell with less than 2 live neighbors dies
			//any live cell with 2 or 3 live neighbors lives
			//any live cell with more than 3 live neighbors dies
			//Any dead cell with 3 live cells becomes alive
			var cellId = strconv.Itoa(rowIndex) + "_" + strconv.Itoa(colIndex)

			if value == 1 && neighborsAlive < 2 {
				cellsAliveCopy[rowIndex][colIndex] = 0
				returnCells += fmt.Sprintf(TableCell, AddCell.String(), cellId, cellId, DeadCSSClass)
			} else if value == 1 && neighborsAlive == 2 || neighborsAlive == 3 {
				cellsAliveCopy[rowIndex][colIndex] = 1
				returnCells += fmt.Sprintf(TableCell, RemoveCell.String(), cellId, cellId, AliveCSSClass)
			} else if value == 1 && neighborsAlive > 3 {
				cellsAliveCopy[rowIndex][colIndex] = 0
				returnCells += fmt.Sprintf(TableCell, AddCell.String(), cellId, cellId, DeadCSSClass)
			} else if value == 0 && neighborsAlive == 3 {
				cellsAliveCopy[rowIndex][colIndex] = 1
				returnCells += fmt.Sprintf(TableCell, RemoveCell.String(), cellId, cellId, AliveCSSClass)
			}
		}
	}

	for i := range AliveCells {
		copy(AliveCells[i], cellsAliveCopy[i])
	}

	return returnCells
}

func reader(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading: %v", err)
		}

		var jsonData GameJSON

		json.Unmarshal(message, &jsonData)

		var messageReturn string

		switch jsonData.Event {
		case GenerateGrid.String():
			// messageReturn = generateCells(jsonData["screenWidth"].(float64), jsonData["screenHeight"].(float64))
			messageReturn = generateCells(jsonData.ScreenWidth, jsonData.ScreenHeight)
		case Start.String():
			GameRunning = true
			messageReturn = fmt.Sprintf(StartStopButton, "error", Stop.String(), "STOP")
		case Stop.String():
			GameRunning = false
			messageReturn = fmt.Sprintf(StartStopButton, "info", Start.String(), "START")
		case AddCell.String():
			setCellAlive(jsonData.CellID)
			messageReturn = fmt.Sprintf(TableCell, RemoveCell.String(), jsonData.CellID, jsonData.CellID, AliveCSSClass)
		case RemoveCell.String():
			setCellDead(jsonData.CellID)
			messageReturn = fmt.Sprintf(TableCell, AddCell.String(), jsonData.CellID, jsonData.CellID, DeadCSSClass)
		default:
			messageReturn = fmt.Sprintf("State doesnt match: %v", jsonData)
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(messageReturn)); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func writer(conn *websocket.Conn) {
	defer conn.Close()
	for {
		if GameRunning {
			start := time.Now()
			message := life()

			time.Sleep(100 * time.Millisecond)
			duration := time.Since(start)
			fmt.Printf("Life iteration took: %s\n", duration)

			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Printf("Life could not run :( - %v", err)
				return
			}
			// time.Sleep(300 * time.Millisecond)
		}
	}
}

func gameOfLifeWS(w http.ResponseWriter, r *http.Request) {
	c, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading ws: %s", err)
		return
	}
	defer func() {
		fmt.Println("closing connection")
		c.Close()
	}()
	go writer(c)
	go reader(c)

	//Keep handler alive
	select {}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/gameoflife", gameOfLifeWS)
	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
