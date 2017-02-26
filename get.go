package main

import (
	"fmt"

	yaml "gopkg.in/yaml.v1"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nabeken/psadm/ps"
	"github.com/pkg/errors"
)

type GetCommand struct {
	At int64 `long:"at" description:"Specify a time of a snapshot of value. Default is now."`
}

func (cmd *GetCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("You must specify a KEY to get.")
	}

	client := ps.NewClient(session.Must(session.NewSession()))

	param, err := client.GetParameter(args[0])
	if err != nil {
		return errors.Wrap(err, "failed to get a parameter")
	}

	out, err := yaml.Marshal([]*ps.Parameter{param})
	if err != nil {
		return errors.Wrap(err, "failed to marshal into YAML")
	}

	fmt.Print(string(out))

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
