package retoenfa

import "testing"

func TestBasicRegex(t *testing.T) {
	trans := NewReToeNFA("1.0.1")
	trans.StartParse()
	enfa := trans.GetEpsNFA()
	enfa.GenerateFormattedTransitionTable()
}

func TestConcatenationRegex(t *testing.T) {
	trans := NewReToeNFA("0+1.0.1")
	trans.StartParse()
	enfa := trans.GetEpsNFA()
	enfa.GenerateFormattedTransitionTable()
}

func TestComplexRegex(t *testing.T) {
	trans := NewReToeNFA("(0+1.0)*.(e+1)")
	trans.StartParse()
	enfa := trans.GetEpsNFA()
	enfa.GenerateFormattedTransitionTable()
}
