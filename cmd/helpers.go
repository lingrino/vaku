package cmd

import (
	"fmt"
)

// newVakuClient returns a vaku client configured with our preferences
// func newVakuClient(cmd *cobra.Command) (*vaku.Client, error) {
// 	var opts []vaku.Option

// 	return vaku.NewClient(
// 		vaku.WithWorkers(flagWorkers),
// 	)
// }

// func authVakuClient(c *vaku.Client) {

// }

func output(val interface{}) {
	switch t := val.(type) {
	case []string:
		outputList(t)
	default:
		fmt.Println(val)
	}
}

func outputList(l []string) {
	for _, v := range l {
		fmt.Println(v)
	}
}
