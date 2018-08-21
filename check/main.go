package main

import (
	"os"
	"fmt"
	"encoding/json"
	"github.com/Evaneos/concourse-gcloudsql-resource/models"
	"github.com/Evaneos/concourse-gcloudsql-resource/commands"
	"sort"
)

func main() {
	var request models.CheckRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)

	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}

	err = request.Source.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid configuration:", err)
		os.Exit(1)
	}

	commands.ActivateServiceAccount(*request.Source.ServiceAccount)

	instances, err := commands.ListInstances(*request.Source.Project)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error while listing instances:", err.Error())
		os.Exit(1)
	}

	sort.Sort(instances)
	output, err := json.Marshal(instances)

	fmt.Fprintf(os.Stdout, "%s\n", output)
}
