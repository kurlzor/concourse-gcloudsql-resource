package main

import (
	"github.com/Evaneos/concourse-gcloudsql-resource/models"
	"encoding/json"
	"os"
	"fmt"
	"github.com/Evaneos/concourse-gcloudsql-resource/commands"
	"io/ioutil"
	"strings"
	"strconv"
)

func main() {
	var writeDir = os.Args[1]
	var request models.InRequest

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

	instanceInfo, err := commands.GetInstanceInfo(request.Version.Instance, *request.Source.Project)

	err = ioutil.WriteFile(writeDir + "/helm_release_name", []byte(instanceInfo.Name + "-gcloudsql-proxy"), 0644)
	commands.Check(err)

	err = ioutil.WriteFile(writeDir + "/instance_name", []byte(instanceInfo.Name), 0644)
	commands.Check(err)

	err = ioutil.WriteFile(writeDir + "/port", []byte(strconv.Itoa(InstanceTypeToPort(instanceInfo.DatabaseVersion))), 0644)
	commands.Check(err)

	err = ioutil.WriteFile(writeDir + "/project", []byte(*request.Source.Project), 0644)
	commands.Check(err)

	err = ioutil.WriteFile(writeDir + "/region", []byte(instanceInfo.Region), 0644)
	commands.Check(err)

	var version = models.ConcourseInOutput{
		Version: models.Version{
			Instance: instanceInfo.Name,
		},
		Metadata: []models.Metadata{
			{
				Name:  "region",
				Value: instanceInfo.Region,
			},
			{
				Name: "instanceType",
				Value: strings.ToLower(instanceInfo.DatabaseVersion),
			},
		},
	}

	jsonOutput, err := json.Marshal(version)

	fmt.Fprintf(os.Stdout, "%s\n", jsonOutput)
}

func InstanceTypeToPort(databaseVersion string) int {
	if strings.HasPrefix(databaseVersion, "POSTGRES") {
		return 5432
	} else if strings.HasPrefix(databaseVersion, "MYSQL") {
		return 3306
	} else {
		panic("unknown database type")
	}
}