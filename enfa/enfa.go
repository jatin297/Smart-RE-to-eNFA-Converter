package enfa

import (
	"fmt"
	. "github.com/jatin297/retoenfa/dto"
	"strings"
)

// CreateENFA initializes a new ENFA with the given initial state and final state designation.
func CreateENFA(initialState int, isFinal bool) *ENFA {
	newENFA := &ENFA{
		initialState: initialState,
		activeStates: make(StateSet),
		states:       []int{},
		transitions:  make(map[TransitionKey]StateSet),
		inputSymbols: make(map[string]bool),
	}
	newENFA.activeStates[initialState] = true
	newENFA.InsertState(initialState, isFinal)
	return newENFA
}

// InsertState adds a new state to the ENFA.
func (e *ENFA) InsertState(state int, isFinal bool) {
	if state == -1 {
		fmt.Println("State -1 is reserved for the dead state and cannot be added.")
		return
	}
	e.states = append(e.states, state)
	if isFinal {
		e.finalStates = append(e.finalStates, state)
	}
}

// DefineTransition sets up a transition between states based on an input symbol.
func (e *ENFA) DefineTransition(startState int, symbol string, endStates ...int) {
	if _, found := e.inputSymbols[symbol]; !found {
		e.inputSymbols[symbol] = true
	}

	stateExists := false
	for _, state := range e.states {
		if state == startState {
			stateExists = true
			break
		}
	}

	if !stateExists {
		fmt.Printf("State %d does not exist in the ENFA.\n", startState)
		return
	}

	destinationStates := make(StateSet)
	for _, destination := range endStates {
		destinationStates[destination] = true
	}

	e.transitions[TransitionKey{startState, symbol}] = destinationStates
}

// IsPathExists checks if a transition exists between two states for the given input symbol.
func (e *ENFA) IsPathExists(source int, input string, destination int) bool {
	if destSet, exists := e.transitions[TransitionKey{source, input}]; exists {
		_, found := destSet[destination]
		return found
	}
	return false
}

// DisplayTransitions outputs the ENFA's transition table.
func (e *ENFA) DisplayTransitions() {
	fmt.Println("===========================================")
	var symbolList []string
	for symbol := range e.inputSymbols {
		if symbol == "" {
			fmt.Printf("\tε|")
		} else {
			fmt.Printf("\t%s|", symbol)
		}
		symbolList = append(symbolList, symbol)
	}
	fmt.Println("\n-------------------------------------------")

	for _, state := range e.states {
		fmt.Printf("%d |", state)
		for _, symbol := range symbolList {
			if destSet, exists := e.transitions[TransitionKey{state, symbol}]; exists {
				fmt.Printf("\t")
				for dest := range destSet {
					fmt.Printf("%d,", dest)
				}
				fmt.Print("|")
			} else {
				fmt.Print("\tNA|")
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------")
	fmt.Println("===========================================")
}

// ProcessInput processes a single input symbol and updates the active states of the ENFA.
func (e *ENFA) ProcessInput(input string) []int {
	newActiveStates := make(StateSet)
	for active := range e.activeStates {
		if nextStates, exists := e.transitions[TransitionKey{active, input}]; exists {
			for dest := range nextStates {
				newActiveStates[dest] = true
				// Handle epsilon transitions
				if epsilonStates, hasEpsilon := e.transitions[TransitionKey{dest, ""}]; hasEpsilon {
					for epsilonDest := range epsilonStates {
						newActiveStates[epsilonDest] = true
					}
				}
			}
		}
	}
	e.activeStates = newActiveStates
	var resultStates []int
	for state := range newActiveStates {
		resultStates = append(resultStates, state)
	}
	return resultStates
}

type ENFA struct {
	initialState int
	activeStates StateSet
	states       []int
	finalStates  []int
	transitions  map[TransitionKey]StateSet
	inputSymbols map[string]bool
}

// CheckIfFinalState verifies if any of the active states is a final state.
func (e *ENFA) CheckIfFinalState() bool {
	for _, finalState := range e.finalStates {
		if e.activeStates[finalState] {
			return true
		}
	}
	return false
}

// ReinitializeActiveStates sets the active states back to the initial state.
func (e *ENFA) ReinitializeActiveStates() {
	e.activeStates = StateSet{e.initialState: true}
}

// ValidateInputSequence determines whether the ENFA accepts a given sequence of input symbols.
func (e *ENFA) ValidateInputSequence(inputs []string) bool {
	for _, inputSymbol := range inputs {
		e.ProcessInput(inputSymbol)
	}
	return e.CheckIfFinalState()
}

// GenerateFormattedTransitionTable creates a structured view of the transition table.
func (e *ENFA) GenerateFormattedTransitionTable() []map[string]string {
	var symbolList []string
	for symbol := range e.inputSymbols {
		symbolList = append(symbolList, symbol)
	}

	var table []map[string]string
	for _, state := range e.states {
		row := make(map[string]string)
		row["state"] = fmt.Sprintf("%d", state)
		for _, symbol := range symbolList {
			destSet, exists := e.transitions[TransitionKey{state, symbol}]
			if len(symbol) == 0 {
				symbol = "ε"
			}
			if exists {
				var destList []string
				for dest := range destSet {
					destList = append(destList, fmt.Sprintf("%d", dest))
				}
				row[symbol] = strings.Join(destList, ",")
			} else {
				row[symbol] = "NA"
			}
		}
		table = append(table, row)
	}
	return table
}
