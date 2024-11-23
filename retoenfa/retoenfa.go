package retoenfa

import (
	"fmt"
	. "github.com/jatin297/retoenfa/dto"
	"github.com/jatin297/retoenfa/enfa"
	"strconv"
)

func NewReToeNFA(str string) *ReToeNFA {
	newRe2NFA := &ReToeNFA{regexString: str}
	newRe2NFA.closureMap = make(map[Closure]bool)
	return newRe2NFA
}

type ReToeNFA struct {
	regexString     string
	nextParentheses []int
	stateCount      int
	closureMap      map[Closure]bool
	enfa            *enfa.ENFA
}

func (r *ReToeNFA) parseRE(expression string, start, end int) (int, int) {

	// Base case: single character
	if start == end {
		initialState := r.incCapacity()
		finalState := r.incCapacity()

		// Handle epsilon transition or regular character
		if expression[start] == 'e' {
			r.addEdge(initialState, Epsilon, finalState)
		} else {
			symbol, _ := strconv.Atoi(string(expression[start]))
			r.addEdge(initialState, symbol, finalState)
		}
		return initialState, finalState
	}

	// Handle grouped expressions enclosed in parentheses
	if expression[start] == '(' && expression[end] == ')' {
		if r.nextParentheses[start] == end {
			// Recursively parse the content inside parentheses
			return r.parseRE(expression, start+1, end-1)
		}
	}

	// Check for union operator (+) and split into subexpressions
	index := start
	for index <= end {
		index = r.nextParentheses[index] // Skip to the corresponding position if parentheses exist

		if index <= end && expression[index] == '+' {
			leftStart, leftEnd := r.parseRE(expression, start, index-1)
			rightStart, rightEnd := r.parseRE(expression, index+1, end)
			unionStart, unionEnd := r.doUnion(leftStart, rightStart, leftEnd, rightEnd)
			return unionStart, unionEnd
		}
		index++
	}

	// Check for concatenation operator (.) and process accordingly
	index = start
	for index <= end {
		index = r.nextParentheses[index] // Skip nested parentheses

		if index <= end && expression[index] == '.' {
			leftStart, leftEnd := r.parseRE(expression, start, index-1)
			rightStart, rightEnd := r.parseRE(expression, index+1, end)
			concatStart, concatEnd := r.doConcatenation(leftStart, rightStart, leftEnd, rightEnd)
			return concatStart, concatEnd
		}
		index++
	}

	subStart, subEnd := r.parseRE(expression, start, end-1)
	closureStart, closureEnd := r.closure(subStart, subEnd)
	return closureStart, closureEnd
}

func (r *ReToeNFA) computeParenthesesMapping(expression string) {
	length := len(expression)
	for index := 0; index < length; index++ {

		// Identify the start of a parenthesis group
		if expression[index] == '(' {
			depth := 0
			current := index

			for {
				if expression[current] == '(' {
					depth++
				}

				if expression[current] == ')' {
					depth--
				}

				// Stop when matching closing parenthesis is found
				if depth == 0 {
					break
				}
				current++
			}

			// Map the index to the corresponding closing parenthesis
			r.nextParentheses = append(r.nextParentheses, current)

		} else {
			// Non-parenthesis characters map to themselves
			r.nextParentheses = append(r.nextParentheses, index)
		}
	}
}

func (r *ReToeNFA) computeStateClosure() {
	// Initialize a temporary queue to process states
	stateQueue := make([]int, 200)

	// Iterate through all states to calculate closures
	for srcState := 0; srcState <= r.stateCount; srcState++ {

		// Assume all states are reachable initially
		for targetState := 0; targetState < r.stateCount; targetState++ {
			r.closureMap[Closure{Src: srcState, Dst: targetState}] = true
		}

		// Setup for breadth-first traversal
		queueStart := -1
		queueEnd := 0
		stateQueue[0] = srcState

		// Each state includes itself in its closure
		r.closureMap[Closure{Src: srcState, Dst: srcState}] = true

		// Process the queue for reachable states
		for queueStart < queueEnd {
			queueStart++
			currentState := stateQueue[queueStart]

			// Explore all possible transitions
			for potentialState := 0; potentialState < r.stateCount; potentialState++ {
				closureAlreadyPresent := r.checkClosureExist(srcState, potentialState)
				validEpsilonPath := r.checkPathExist(currentState, Epsilon, potentialState)

				// Add the state to the closure if it's reachable and not already present
				if !closureAlreadyPresent && validEpsilonPath {
					queueEnd++
					stateQueue[queueEnd] = potentialState

					r.closureMap[Closure{Src: srcState, Dst: potentialState}] = true
				}
			}
		}
	}
}

func (r *ReToeNFA) incCapacity() int {
	if r.enfa == nil {
		r.enfa = enfa.CreateENFA(0, false)
	} else {
		r.enfa.InsertState(r.stateCount-1, false)
	}
	r.stateCount = r.stateCount + 1
	return r.stateCount - 1
}

func (r *ReToeNFA) addEdge(stateSrc int, cInput int, stateDst int) {
	var inputString string
	if cInput != 2 {
		inputString = strconv.Itoa(cInput)
	}
	r.enfa.DefineTransition(stateSrc, inputString, stateDst)
}

func (r *ReToeNFA) doUnion(s1, s2, t1, t2 int) (int, int) {
	newStartState := r.incCapacity()
	newFinalState := r.incCapacity()

	r.addEdge(newStartState, 2, s1)
	r.addEdge(newStartState, 2, s2)

	r.addEdge(t1, 2, newFinalState)
	r.addEdge(t2, 2, newFinalState)

	return newStartState, newFinalState
}

func (r *ReToeNFA) StartParse() {
	r.computeParenthesesMapping(r.regexString)
	nfaStart, nfaFinal := r.parseRE(r.regexString, 0, len(r.regexString)-1)
	fmt.Printf("NFA s=%d, f=%d\n", nfaStart, nfaFinal)
}

func (r *ReToeNFA) GetEpsNFA() *enfa.ENFA {
	return r.enfa
}

func (r *ReToeNFA) doConcatenation(s1, s2, t1, t2 int) (int, int) {
	r.addEdge(t1, 2, s2)
	return s1, t2
}

func (r *ReToeNFA) closure(s, t int) (int, int) {
	newStartState := r.incCapacity()
	newFinalState := r.incCapacity()

	r.addEdge(newStartState, 2, s)
	r.addEdge(t, 2, newFinalState)
	r.addEdge(t, 2, s)
	r.addEdge(newStartState, 2, newFinalState)
	return newStartState, newFinalState
}

func (r *ReToeNFA) checkClosureExist(src, target int) bool {
	closureExist, _ := r.closureMap[Closure{Src: src, Dst: target}]
	return closureExist
}

func (r *ReToeNFA) checkPathExist(src, input, dst int) bool {
	if r.enfa == nil {
		return false
	}

	return r.enfa.IsPathExists(src, strconv.Itoa(input), dst)
}
