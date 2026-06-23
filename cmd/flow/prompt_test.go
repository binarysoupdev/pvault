package flow_test

import (
	"testing"

	"pvault/cmd/flow"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
)

func TestPromptReturnsInput(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	PROMPT := rand.ASCII(10)
	INPUT := rand.ASCII(15)

	io := pipe.OpenStdio(1, 1, true)
	defer io.Close()

	//-- act
	io.Queue(PROMPT, INPUT)
	io.EndQueue()

	res := flow.Prompt(PROMPT)

	//-- assert
	assert.Equal(t, INPUT, res)
	assert.Contains(t, io.ReadLine(), PROMPT+INPUT)
}
