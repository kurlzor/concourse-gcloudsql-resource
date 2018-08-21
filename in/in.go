package main

import (
	"github.com/Evaneos/concourse-gcloudsql-resource/models"
	"encoding/json"
	"os"
	"fmt"
	"github.com/Evaneos/concourse-gcloudsql-resource/commands"
	"io/ioutil"
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
	check(err)

	err = ioutil.WriteFile(writeDir + "/instance_name", []byte(instanceInfo.Name), 0644)
	check(err)

	err = ioutil.WriteFile(writeDir + "/port", []byte(instanceInfo.Port), 0644)
	check(err)

	err = ioutil.WriteFile(writeDir + "/project", []byte(*request.Source.Project), 0644)
	check(err)

	err = ioutil.WriteFile(writeDir + "/region", []byte(instanceInfo.Region), 0644)
	check(err)

	output, err := json.Marshal(instanceInfo)

	fmt.Fprintf(os.Stdout, "%s\n", output)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
