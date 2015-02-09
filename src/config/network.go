package config

import (
	"fmt"
)


type connections []string

func (i *connections) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *connections) Set(value string) error {
	*i = append(*i, value )
	return nil
}

var Zconnections connections

