package commands

import (
	"encoding/json"
)

type UserCommand struct {
	Command string
	Uuid    string
}

func New(str string) *UserCommand {
	var u = &UserCommand{}
	err := json.Unmarshal([]byte(str), u)
	if err != nil {
		// ничего не делаем
	}
	return u
}

func (c *UserCommand) ToJson() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (c *UserCommand) IsNotEmpty() bool {
	return c.Command != ""
}
