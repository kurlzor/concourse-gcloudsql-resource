package commands

import (
	"os/exec"
	"fmt"
	"encoding/json"
	"github.com/Evaneos/concourse-gcloudsql-resource/models"
	"errors"
)

var baseCommand = "gcloud"
var baseArgs = []string {"sql", "instances"}

func buildCommand(args ...string) *exec.Cmd {
	return exec.Command(baseCommand, append(baseArgs, args...)...)
}

func ListInstances(project string) (models.GCloudSQLInstanceList, error) {
	output, err := buildCommand("list", fmt.Sprintf("--project=%s", project), "--format=json").CombinedOutput()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", output, err.Error()))
	}

	instances := make(models.GCloudSQLInstanceList, 0)
	json.Unmarshal(output, instances)

	return instances, err
}

//func GetInstanceInfo(instance string, project string) (error) {
//	output, err := exec.Command(fmt.Sprintf("%s describe %s --project=%s --format=json", baseCommand, instance, project)).CombinedOutput()
//
//
//	return instanceInfo, err
//}
