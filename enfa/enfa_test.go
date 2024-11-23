package enfa

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ENFATestSuite struct {
	suite.Suite
	enfa *ENFA
}

func (suite *ENFATestSuite) SetupTest() {
	suite.enfa = CreateENFA(0, false)
}

func TestENFA(t *testing.T) {
	suite.Run(t, new(ENFATestSuite))
}

func (suite *ENFATestSuite) TestBasic() {
	suite.SetupTest()
	t := suite.T()

	nfa := suite.enfa
	nfa.InsertState(1, false)
	nfa.InsertState(2, true)

	nfa.DefineTransition(0, "a", 1)
	nfa.DefineTransition(1, "b", 2)

	if ret := nfa.ProcessInput("a"); ret[0] != 1 {
		t.Errorf("Expect 1, but get %d\n", ret)
	}

	if ret := nfa.ProcessInput("b"); ret[0] != 2 {
		t.Errorf("Expect 2, but get %d\n", ret)
	}

	if !nfa.CheckIfFinalState() {
		t.Errorf("Verify is failed")
	}
}

func (suite *ENFATestSuite) TestValidateInputSequence() {
	suite.SetupTest()
	t := suite.T()

	nfa := suite.enfa
	nfa.InsertState(1, false)
	nfa.InsertState(2, true)

	nfa.DefineTransition(0, "a", 1)
	nfa.DefineTransition(1, "b", 2)

	var inputs []string
	inputs = append(inputs, "a")
	inputs = append(inputs, "b")
	if !nfa.ValidateInputSequence(inputs) {
		t.Errorf("Verify Inputs is failed")
	}

	nfa.GenerateFormattedTransitionTable()
}

func (suite *ENFATestSuite) TestAdvanceENFA() {
	suite.SetupTest()
	t := suite.T()

	nfa := suite.enfa
	nfa.InsertState(1, true)
	nfa.InsertState(2, false)

	nfa.DefineTransition(0, "0", 0)
	nfa.DefineTransition(0, "1", 1)
	nfa.DefineTransition(1, "0", 0)
	nfa.DefineTransition(1, "1", 2)
	nfa.DefineTransition(2, "0", 2)
	nfa.DefineTransition(2, "1", 2)

	nfa.GenerateFormattedTransitionTable()
	inputs := []string{"0", "0", "1", "0", "1"}

	if !nfa.ValidateInputSequence(inputs) {
		t.Errorf("Verify inputs is failed")
	}

	//Reset the nfa for another verification
	nfa.ReinitializeActiveStates()

	//Test go to dead state 2

	inputs2 := []string{"1", "1", "0", "0", "0"}

	if nfa.ValidateInputSequence(inputs2) {
		t.Errorf("Verify inputs is failed")
	}
}

func (suite *ENFATestSuite) TestNFA() {
	suite.SetupTest()
	t := suite.T()

	nfa := suite.enfa
	nfa.InsertState(1, true)
	nfa.InsertState(2, false)

	nfa.DefineTransition(0, "0", 0, 1)
	nfa.DefineTransition(0, "1", 1)
	nfa.DefineTransition(1, "0", 0)
	nfa.DefineTransition(1, "1", 2)
	nfa.DefineTransition(2, "0", 2)
	nfa.DefineTransition(2, "1", 2, 0)
	nfa.GenerateFormattedTransitionTable()
	inputs := []string{"0", "0", "1", "0", "1"}

	if !nfa.ValidateInputSequence(inputs) {
		t.Errorf("Verify inputs is failed")
	}

	inputs2 := []string{"0", "0", "0", "0", "1"}

	if !nfa.ValidateInputSequence(inputs2) {
		t.Errorf("Verify inputs2 is failed")

	}

	inputs3 := []string{"0", "1", "2"}

	if nfa.ValidateInputSequence(inputs3) {
		t.Errorf("Verify inputs3 is failed")
	}
}

func (suite *ENFATestSuite) TestEpsilonNFA() {
	suite.SetupTest()
	t := suite.T()

	nfa := suite.enfa
	nfa.InsertState(1, false)
	nfa.InsertState(2, false)
	nfa.InsertState(3, true)
	nfa.InsertState(4, false)
	nfa.InsertState(5, false)

	nfa.DefineTransition(0, "1", 1)
	nfa.DefineTransition(0, "0", 4)

	nfa.DefineTransition(1, "1", 2)
	nfa.DefineTransition(1, "", 3) //epsilon
	nfa.DefineTransition(2, "1", 3)
	nfa.DefineTransition(4, "0", 5)
	nfa.DefineTransition(4, "", 1, 2) //E -> epsilon B C
	nfa.DefineTransition(5, "0", 3)

	nfa.GenerateFormattedTransitionTable()

	if !nfa.ValidateInputSequence([]string{"1"}) {
		t.Errorf("Verify inputs is failed")
	}

	nfa.ReinitializeActiveStates()

	if !nfa.ValidateInputSequence([]string{"1", "1", "1"}) {
		t.Errorf("Verify inputs is failed")
	}

	nfa.ReinitializeActiveStates()

	if !nfa.ValidateInputSequence([]string{"0", "1"}) {
		t.Errorf("Verify inputs is failed")
	}

	nfa.ReinitializeActiveStates()
	if !nfa.ValidateInputSequence([]string{"0", "0", "0"}) {
		t.Errorf("Verify inputs is failed")
	}
}
