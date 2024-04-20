package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/goccy/go-yaml"
	"github.com/nabeken/psadm"
	"github.com/pkg/errors"
)

type ExportCommand struct {
	KeyPrefix string `long:"key-prefix" description:"Specify a key prefix to be exported"`
}

func (cmd *ExportCommand) Execute(args []string) error {
	client := psadm.NewClient(session.Must(session.NewSession()))

	params, err := client.GetParametersByPath(cmd.KeyPrefix)
	if err != nil {
		return err
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
