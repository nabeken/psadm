package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nabeken/psadm"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type GetCommand struct {
	At         string `long:"at" description:"Specify a time of a snapshot of value. Default is now."`
	PrintValue bool   `long:"print-value" description:"Print a value instead of YAML"`
}

func (cmd *GetCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("You must specify a KEY to get.")
	}

	client := psadm.NewClient(session.Must(session.NewSession()))

	var param *psadm.Parameter
	var err error
	if cmd.At == "" {
		param, err = client.GetParameterWithDescription(args[0])
	} else {
		at, err := time.Parse(time.RFC3339, cmd.At)
		if err != nil {
			return errors.Wrap(err, "failed to parse `at'.")
		}
		param, err = client.GetParameterByTime(args[0], at)
	}
	if err != nil {
		return err
	}

	if cmd.PrintValue {
		fmt.Println(param.Value)
	} else {
		out, err := yaml.Marshal([]*psadm.Parameter{param})
		if err != nil {
			return errors.Wrap(err, "failed to marshal into YAML")
		}

		fmt.Print(string(out))
	}

	return nil
}

func init() {
	parser.AddCommand(
		"get",
		"Get a parameter at given time",
		"The get command gets a parameter in YAML at given time from PS.",
		&GetCommand{},
	)
}
