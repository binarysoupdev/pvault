package flow_test

import (
	"testing"

	"pvault/cmd/flow"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
)

func TestPromptReturnsInput(t *testing.T) {
	//-- arrange
	const PROMPT = "prompt: "
	const INPUT = "input"

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

func TestPromptPasswordReturnsPassword(t *testing.T) {
	//-- arrange
	const PROMPT = "prompt: "
	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 1, true)
	defer io.Close()

	//-- act
	io.Queue(PROMPT, PASSWORD)
	io.EndQueue()

	res := flow.PromptPassword(PROMPT)

	//-- assert
	assert.Equal(t, PASSWORD, res)
	assert.Contains(t, io.ReadLine(), PROMPT+PASSWORD)
}
