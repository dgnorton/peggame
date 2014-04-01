// You know those triangle shaped puzzles with pegs in them?  This app
// populates the board with a random open slot, finds all solutions,
// and optionally prints a specified number of them.
//
//     0              0
//    1 1           1   2
//   1 1 1        3   4   5
//  1 1 1 1     6   7   8   9
// 1 1 1 1 1  10  11  12  13  14
//
// We use a single unsigned integer to represent the state of a game.
// The value of the LSB represents hole 0 and the value of bit
// (1 << 14) represents hole 14.  Simple, small memory footprint,
// and runs reasonably fast.
package main

import (
    "flag"
	"fmt"
	"math/rand"
	"time"
)

// possible move
type move struct {
    // P2 is the peg in the middle that gets jumped by either
    // P1 or P3, depending on which has a peg in it.
	P1, P2, P3 uint
}

// all possible moves
var pmoves = []move{
	move{0, 1, 3},
	move{0, 2, 5},
	move{1, 3, 6},
	move{1, 4, 8},
	move{2, 4, 7},
	move{2, 5, 9},
	move{3, 4, 5},
	move{3, 7, 12},
	move{3, 6, 10},
	move{4, 7, 11},
	move{4, 8, 13},
	move{5, 8, 12},
	move{5, 9, 14},
	move{6, 7, 8},
	move{7, 8, 9},
	move{10, 11, 12},
	move{12, 13, 14},
	move{13, 12, 11},
}

const maxUint = ^uint(0)

// game represents the peg board in its current state.  We use a single
// unsigned integer to represent the current state of the game.
type game struct {
	Board uint
}

// creates a new instance of a game ready to play
func newGame() game {
    // generate a board with pegs in ALL holes
	g := game{maxUint}
    // remove one random peg
	n := uint(rand.Intn(15))
	g = g.toggleBit(n)
	return g
}

// get the value (1 or 0) of the specified bit
func (g game) bitValue(hole uint) uint {
	if (g.Board & (1 << hole)) > 0 {
		return 1
	}
	return 0
}

// toggle the specified bit
func (g game) toggleBit(n uint) game {
	g.Board ^= (1 << n)
	return g
}

// returns true if the specified move can be played
func (g game) CanPlay(m move) bool {
	if g.bitValue(m.P2) == 0 {
		return false
	}
	return g.bitValue(m.P1) != g.bitValue(m.P3)
}

// returns the count of pegs remaining in the game
func (g game) PegCnt() int {
	cnt := 0
	for n := uint(0); n < 15; n++ {
		if g.bitValue(n) == 1 {
			cnt++
		}
	}
	return cnt
}

// prints an ASCII art representation of the game
func (g game) Print() {
	fmt.Printf("    %d\n", g.bitValue(0))
	fmt.Printf("   %d %d\n", g.bitValue(1), g.bitValue(2))
	fmt.Printf("  %d %d %d\n", g.bitValue(3), g.bitValue(4), g.bitValue(5))
	fmt.Printf(" %d %d %d %d\n", g.bitValue(6), g.bitValue(7), g.bitValue(8), g.bitValue(9))
	fmt.Printf("%d %d %d %d %d\n", g.bitValue(10), g.bitValue(11), g.bitValue(12), g.bitValue(13), g.bitValue(14))
}

// plays the specified move
func (g game) Play(m move) game {
	g = g.toggleBit(m.P1)
	g = g.toggleBit(m.P2)
	g = g.toggleBit(m.P3)
	return g
}

// recursive function that plays all possible paths of the specified game
func play(g game, moves []uint, solvedCh chan []uint) {
	moved := false
	for _, m := range pmoves {
		if g.CanPlay(m) {
			g2 := g.Play(m)
			mvs := append([]uint(nil), moves...)
			mvs = append(mvs, g2.Board)
			play(g.Play(m), mvs, solvedCh)
			moved = true
		}
	}

	if !moved && g.PegCnt() == 1 {
		solvedCh <- moves
	}

	if len(moves) == 1 {
		close(solvedCh)
	}
}

func main() {
    // parse command line
    var printCnt int
    flag.IntVar(&printCnt, "p", 1, "number of solutions to print")
    flag.Parse()

    // seed random generator
	rand.Seed(time.Now().UTC().UnixNano())

   // create a game board and start the solver in the background
	g := newGame()
	solutionsCh := make(chan []uint)
	moves := []uint{g.Board}
	go play(g, moves, solutionsCh)

    // loop to read solutions as the solver finds them
	for i := 0; ; i++ {
		solution := <-solutionsCh

		if solution == nil {
            if printCnt < 0 {
               printCnt = i
            }
			fmt.Printf("Printed %d of %d solutions.\n", printCnt, i+1)
			return
		}

        if i < printCnt || printCnt < 0{
		 for _, move := range solution {
		    game{move}.Print()
		    fmt.Println(" ")
		 }
        }
	}
}
