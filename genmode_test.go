package mode

import (
	"math/rand"
	"reflect"
	"strings"
	"testing/quick"
)

type mode string

var _ quick.Generator = mode("")

/*
From chmod(1):

     The symbolic mode is described by the following grammar:

           mode         ::= clause [, clause ...]
           clause       ::= [who ...] [action ...] action
           action       ::= op [perm ...]
           who          ::= a | u | g | o
           op           ::= + | - | =
           perm         ::= r | s | t | w | x | X | u | g | o

*/

// Generate implements quick.Generator
func (mode) Generate(rand *rand.Rand, size int) reflect.Value {
	n := int(rand.ExpFloat64()*float64(size)/5) + 1
	clauses := make([]string, n)
	for i := range clauses {
		clauses[i] = randClause(rand)
	}
	v := mode(strings.Join(clauses, ","))
	//log.Println("mode.Generate:", size, n, v)
	return reflect.ValueOf(v)
}

func randClause(rand *rand.Rand) string {
	return randWho(rand) + randAction(rand)
}

func randWho(rand *rand.Rand) string {
	return randomFromSet(rand, "augo")
}

func randAction(rand *rand.Rand) string {
	action := ""
	for n := int(rand.ExpFloat64() * 1 / 5); n >= 0; n-- {
		action += randOp(rand) + randPerm(rand)
	}
	return action
}

func randOp(rand *rand.Rand) string {
	switch rand.Intn(7) {
	case 0:
		return "="
	case 1, 2, 3:
		return "+"
	default:
		return "-"
	}
}

func randPerm(rand *rand.Rand) string {
	return randomFromSet(rand, "rstwxXugo")
}

func randomFromSet(rand *rand.Rand, set string) string {
	max := int(rand.ExpFloat64() * float64(len(set)) / 5)
	//have := make(map[uint]bool, len(set))
	have := uint(0)
	result := ""
	for n := 0; n < max; n++ {
		i := uint(rand.Intn(len(set)))
		//if !have[i] {
		if have&(1<<i) == 0 {
			result += set[i : i+1]
			//have[i] = true
			have |= 1 << i
		}
	}
	return result
}
