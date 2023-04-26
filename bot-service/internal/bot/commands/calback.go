package commands

import (
	"encoding/json"
	"fmt"
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

func (c *UserCommand) ToJson() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal object to JSON: %#v", err.Error())
	}
	return string(bytes), nil
}

func (c *UserCommand) IsNotEmpty() bool {
	return c.Command != ""
}
