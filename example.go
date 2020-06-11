package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)



type Board struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Move struct {
	X1 int `json:"x_1"`
	Y1 int `json:"y_1"`
	X2 int `json:"x_2"`
	Y2 int `json:"y_2"`
}

type Message struct {
	Board []Board `json:"board"`
	Move  Move    `json:"move"`
	GameType int  `json:"game_type"`
}

const (
	WIDTH = 12
	HEIGHT = 12
)

func validateBorder(x1 int, x2 int, y1 int, y2 int, gameType int) bool {
	fmt.Println("inValidate")
	if  gameType == 4 && (x1 < 0 || x2 < 0 || y1 < 0 || y2 < 0 || ((y2 == 0 || y2 == 1 || y2 == 10 || y2 == 11) &&
		(x2 == 0 || x2 == 1 || x2 == 10 || x2 == 11))){
			return false
	}
	if gameType == 2 && (x1<0 || x2<0 || y1<0 || y2<0 || x1>7 || x2>7 || y1>7 || y2>7){
		fmt.Println("in game 2")
		return false
	}
	return true
}

func isEmpty(x int, y int, board []Board, c chan bool, p chan Board) {
	for _, b := range board{
		if b.X == x && b.Y == y{
			fmt.Println(b)
			c <-false
			p <- b
			return
		}
	}
	c <-true
}

func isEmptyOnRoad(x int, y int, board []Board) bool{
	for _, b := range board{
		if b.X == x && b.Y == y{
			return false
		}
	}
	return true
}

func checkPawnMoves(move *Move, piece *Board, board []Board, gameType int) bool {
	fmt.Println("in pawn moves")
	var x = move.X1
	var y = move.Y1
	var c = make(chan bool)
	var p = make(chan Board)
	if gameType == 4 {
		// налево и направо
		if (x+1 == move.X2 || x-1 == move.X2) && move.Y2 == y {
			fmt.Println("vverh")
			go isEmpty(move.X2, move.Y2, board, c, p)
			if <-c {
				fmt.Println("vverh")
				return true
			}
		} else if x == move.X2 && (move.Y2 == y+1 || move.Y2 == y-1) { //вверх и вниз
			fmt.Println("right")
			go isEmpty(move.X2, move.Y2, board, c, p)
			if <-c {
				fmt.Println("right")
				return true
			}
		} else if (x+1 == move.X2 && y+1 == move.Y2) || (x+1 == move.X2 && y-1 == move.Y2) ||
			(x-1 == move.X2 && y-1 == move.Y2) || (x-1 == move.X2 && y+1 == move.Y2) {
			fmt.Println("zheu")
			go isEmpty(move.X2, move.Y2, board, c, p)
			if !<-c {
				var result = <-p
				if piece.Color != result.Color {
					return true
				}
			}
		}else if (((move.X1 == 1 || move.X1 == 10) && (math.Abs(float64(move.X2-move.X1)) == 2)) ||
			((move.Y1 == 1 || move.Y1 == 10) && math.Abs(float64(move.Y2-move.Y1))==2)) &&
			(move.X1 == move.X2 || move.Y2 == move.Y1){
			fmt.Println("in function else if")
			return true
		}
		return false
	}else if gameType == 2 {
		if ((move.Y1 == 1 || move.Y1 == 6) && (math.Abs(float64(move.Y1-move.Y2)) == 2)) &&
			(move.X1 == move.X2){
			return true
		}else if (math.Abs(float64(move.Y2-move.Y1)) == 1) && (move.X1 == move.X2) {
			return true
		}
		go isEmpty(move.X2, move.Y2, board, c, p)
		if !<-c {
			var result = <-p
			if piece.Color != result.Color {
				return true
			}
		}
		return false
	}
	return false
}

func checkKnightMoves(move *Move, piece *Board, board []Board) bool {
	//(math.Abs(float64(move.Y1-move.Y2)) != math.Abs(float64(move.X1-move.X2)))
	if (math.Abs(float64(move.X1 - move.X2)) == 2 && math.Abs(float64(move.Y1-move.Y2)) == 1) ||
			(math.Abs(float64(move.X1 - move.X2)) == 1 && math.Abs(float64(move.Y1-move.Y2)) == 2) {
		var c = make(chan bool)
		var p = make(chan Board)
		go isEmpty(move.X2, move.Y2, board, c, p)
		if <-c{
			return true
		}else {
			var result = <-p
			if piece.Color != result.Color{
				return true
			}
		}
	}
	return false
}


func checkRookMoves(move *Move, piece *Board, board []Board) bool {
	var c = make(chan bool)
	var p = make(chan Board)
	var x = move.X1
	var y = move.Y1
	go isEmpty(move.X2, move.Y2, board, c, p)
	var dif = 0
	if move.X2 == move.X1{
		dif = int(math.Abs(float64(move.Y2 - move.Y1))) // difference
	}else {
		dif = int(math.Abs(float64(move.X1 - move.X2)))
	}
	for i := 0; i < dif-1; i++ {
		if move.X2 == move.X1{
			if move.Y1 > move.Y2{
				y = y - 1
			}else {
				y = y + 1
			}
		}else {
			if move.X1 > move.X2{
				x = x - 1
			}else {
				x = x + 1
			}
		}
		if !(isEmptyOnRoad(x, y, board)) {
			return false
		}
	}
	if <-c{
		return true
	}else {
		var result = <-p
		fmt.Println(piece.Color)
		if piece.Color != result.Color{
			return true
		}
	}
	return false
}

func checkQueenMoves(move *Move, piece *Board, board []Board) bool {
	var c = make(chan bool)
	var p = make(chan Board)
	go isEmpty(move.X2, move.Y2, board, c, p)
	var x = move.X1
	var y = move.Y1
	var dif = 0
	if move.X2 == move.X1 || move.Y1 == move.Y2{
		if move.X2 == move.X1{
			dif = int(math.Abs(float64(move.Y2 - move.Y1))) // difference
		}else {
			dif = int(math.Abs(float64(move.X1 - move.X2)))
		}
		for i := 0; i < dif-1; i++ {
			if move.X2 == move.X1{
				if move.Y1 > move.Y2{
					y = y - 1
				}else {
					y = y + 1
				}
			}else {
				if move.X1 > move.X2{
					x = x - 1
				}else {
					x = x + 1
				}
			}
			if !(isEmptyOnRoad(x, y, board)) {
				return false
			}
		}
	}else {
		dif = int(math.Abs(float64(move.X1 - move.X2)))
		for i := 0; i < dif-1; i++{
			if move.X1 < move.X2{
				x = x + 1
			}else {
				x = x - 1
			}
			if move.Y2 > move.Y1{
				y = y + 1
			}else {
				y = y - 1
			}
			if !(isEmptyOnRoad(x, y, board)){
				return false
			}
		}
	}
	if <-c{
		return true
	}else {
		var result = <-p
		if piece.Color != result.Color{
			return true
		}
	}
	return false
}

func checkBishopMoves(move *Move, piece *Board, board []Board) bool  {
	var c = make(chan bool)
	var p = make(chan Board)
	go isEmpty(move.X2, move.Y2, board, c, p)
	var dif = int(math.Abs(float64(move.X1 - move.X2)))
	var x = move.X1
	var y = move.Y1
	for i := 0; i < dif-1; i++{
		if move.X1 < move.X2{
			x = x + 1
		}else {
			x = x - 1
		}
		if move.Y2 > move.Y1{
			y = y + 1
		}else {
			y = y - 1
		}
		if !(isEmptyOnRoad(x, y, board)){
			return false
		}
	}
	if <-c{
		return true
	}else {
		var result = <-p
		if piece.Color != result.Color{
			return true
		}
	}
	return false
}

func checkKingMoves(move *Move, piece *Board, board []Board) bool {
	var x = move.X1
	var y = move.Y1
	if ((x + 1 == move.X2 || x - 1 == move.X2) && move.Y2 == y) ||(x == move.X2 && (move.Y2 == y+1 || move.Y2 == y-1)) ||
		((x+1 == move.X2 && y+1 == move.Y2) || (x+1 == move.X2 && y-1 == move.Y2) ||
			(x-1 == move.X2 && y-1 == move.Y2) || (x-1 == move.X2 && y+1 == move.Y2)){
		var c = make(chan bool)
		var p = make(chan Board)
		go isEmpty(move.X2, move.Y2, board, c, p)
		if <-c{
			return true
		}else {
			var result = <-p
			if piece.Color != result.Color{
				return true
			}
		}
	}
	return false
}

func chessValidate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var message Message
	result := json.Unmarshal(reqBody, &message)
	var board = message.Board
	var move = message.Move
	var gameType = message.GameType
	var piece Board
	for _, b := range board{
		if b.X == move.X1 && b.Y == move.Y1{
			piece = b
		}
	}
	if result == nil && validateBorder(move.X1, move.X2, move.Y1, move.Y2, gameType){
		switch piece.Name {
		// Пешка
		case "pawn":{
			fmt.Println("pawn")
			fmt.Println(checkPawnMoves(&move, &piece, board, gameType))
			if checkPawnMoves(&move, &piece, board, gameType){
				_ = json.NewEncoder(w).Encode(true)
				return
			}
			_ = json.NewEncoder(w).Encode(false)

		}
		// Конь
		case "knight":{
			fmt.Println("knight")
			if checkKnightMoves(&move, &piece, board){
				json.NewEncoder(w).Encode(true)
				return
			}
			json.NewEncoder(w).Encode(false)
			return
		}
		// Слон
		case "bishop":{
			if (math.Abs(float64(move.X2-move.X1)) == math.Abs(float64(move.Y2-move.Y1))) &&
				checkBishopMoves(&move, &piece, board){
					json.NewEncoder(w).Encode(true)
					return
			}
			json.NewEncoder(w).Encode(false)
		}
		// Ладья
		case "rook":{
			if (move.X2 == move.X1 || move.Y2 == move.Y1) && (checkRookMoves(&move, &piece, board)){
				json.NewEncoder(w).Encode(true)
				return
			}
			json.NewEncoder(w).Encode(false)
		}
		// Ферзя
		case "queen":{
			if ((move.X2 == move.X1 || move.Y2 == move.Y1) || (math.Abs(float64(move.X2-move.X1)) ==
				math.Abs(float64(move.Y2-move.Y1)))) && checkQueenMoves(&move, &piece, board){
				json.NewEncoder(w).Encode(true)
				return
			}
			json.NewEncoder(w).Encode(false)
		}
		// Король
		case "king":{
			if checkKingMoves(&move, &piece, board){
				json.NewEncoder(w).Encode(true)
				return
			}
			json.NewEncoder(w).Encode(false)
		}
		default:
			json.NewEncoder(w).Encode(false)
			return
		}
	}else {
		json.NewEncoder(w).Encode(false)
		return
	}
	return
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/validate", chessValidate).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}