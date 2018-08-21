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
		panic("parse error:" + err.Error())
	}

	err = request.Source.Validate()
	if err != nil {
		panic("invalid configuration:" + err.Error())
	}

	commands.ActivateServiceAccount(*request.Source.ServiceAccount)

	instanceInfo := commands.GetInstanceInfo(request.Version.Instance, *request.Source.Project)

	WriteInfoToOutputFile(writeDir + "/helm_release_name", instanceInfo.Name + "-gcloudsql-proxy")
	WriteInfoToOutputFile(writeDir + "/instance_name", instanceInfo.Name)
	WriteInfoToOutputFile(writeDir + "/port", strconv.Itoa(InstanceTypeToPort(instanceInfo.DatabaseVersion)))
	WriteInfoToOutputFile(writeDir + "/project", *request.Source.Project)
	WriteInfoToOutputFile(writeDir + "/region", instanceInfo.Region)

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
	commands.CheckError(err)

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

func WriteInfoToOutputFile(path string, data string) {
	err := ioutil.WriteFile(path, []byte(data), 0644)
	commands.CheckError(err)
}