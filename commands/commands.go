package commands

import (
	"os/exec"
	"fmt"
	"encoding/json"
	"github.com/Evaneos/concourse-gcloudsql-resource/models"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
)

func runCommand(command *exec.Cmd) ([]byte, error) {
	output, err := command.CombinedOutput()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", output, err.Error()))
	}

	return output, err
}

func ListInstances(project string) (models.GCloudSQLInstanceList) {
	args := []string {
		"sql",
		"instances",
		"list",
		fmt.Sprintf("--project=%s", project),
		"--format=json",
	}

	output, err := runCommand(exec.Command("gcloud", args...))
	CheckError(err)

	instances := make(models.GCloudSQLInstanceList, 0)
	err = json.Unmarshal(output, &instances)
	CheckError(err)

	return instances
}

func ActivateServiceAccount(serviceAccount string) {
	saPath := WriteServiceAccountToFile(serviceAccount)

	args := []string {
		"auth",
		"activate-service-account",
		fmt.Sprintf("--key-file=%s", saPath),
		"--quiet",
	}

	_, err := runCommand(exec.Command("gcloud", args...))
	CheckError(err)
}

func GetInstanceInfo(name string, project string) (models.GCloudSQLInstance) {
	args := []string {
		"sql",
		"instances",
		"describe",
		name,
		fmt.Sprintf("--project=%s", project),
		"--format=json",
	}

	output, err := runCommand(exec.Command("gcloud", args...))
	CheckError(err)

	var instance models.GCloudSQLInstance
	err = json.Unmarshal(output, &instance)
	CheckError(err)

	return instance
}

func WriteServiceAccountToFile(serviceAccount string) string {
	saUuid, err := uuid.NewRandom()
	CheckError(err)

	saPath := "/tmp/%s" + saUuid.String()

	err = ioutil.WriteFile(saPath, []byte(serviceAccount), 0644)
	CheckError(err)

	return saPath
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
