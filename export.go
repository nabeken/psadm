package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nabeken/psadm/ps"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ExportCommand struct {
	KeyPrefix string `long:"key-prefix" description:"Specify a key prefix to be exported"`
}

func (cmd *ExportCommand) Execute(args []string) error {
	client := ps.NewClient(session.Must(session.NewSession()))

	params, err := client.GetAllParameters(cmd.KeyPrefix)
	if err != nil {
		return errors.Wrap(err, "failed to get parameters from Parameter Store")
	}

	out, err := yaml.Marshal(params)
	if err != nil {
		return errors.Wrap(err, "failed to marshal into YAML")
	}

	fmt.Print(string(out))

	return nil
}

func init() {
	parser.AddCommand(
		"export",
		"Export parameters",
		"The export command exports parameters from Parameter Store.",
		&ExportCommand{},
	)
}
