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
		panic("parse error:" + err.Error())
	}

	err = request.Source.Validate()

	if err != nil {
		panic("invalid configuration:" + err.Error())
	}

	commands.ActivateServiceAccount(*request.Source.ServiceAccount)

	instances := commands.ListInstances(*request.Source.Project)

	sort.Sort(instances)

	versions := make([]models.Version, len(instances))

	for i, instance := range instances {
		versions[i].Instance = instance.Name
	}

	output, err := json.Marshal(versions)
	commands.CheckError(err)

	fmt.Fprintf(os.Stdout, "%s\n", output)
}