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

func ListInstances(project string) (models.GCloudSQLInstanceList, error) {
	args := []string {
		"sql",
		"instances",
		"list",
		fmt.Sprintf("--project=%s", project),
		"--format=json",
	}

	output, err := runCommand(exec.Command("gcloud", args...))

	if err != nil {
		return nil, err
	}

	instances := make(models.GCloudSQLInstanceList, 0)
	err = json.Unmarshal(output, &instances)

	if err != nil {
		panic(err)
	}

	return instances, err
}

func ActivateServiceAccount(serviceAccount string) error {
	saPath := WriteServiceAccountToFile(serviceAccount)

	args := []string {
		"auth",
		"activate-service-account",
		fmt.Sprintf("--key-file=%s", saPath),
		"--quiet",
	}

	_, err := runCommand(exec.Command("gcloud", args...))

	return err
}

func GetInstanceInfo(name string, project string) (models.GCloudSQLInstance, error) {
	args := []string {
		"sql",
		"instances",
		"describe",
		name,
		fmt.Sprintf("--project=%s", project),
		"--format=json",
	}

	output, err := runCommand(exec.Command("gcloud", args...))

	if err != nil {
		panic(err)
	}

	var instance models.GCloudSQLInstance
	err = json.Unmarshal(output, &instance)

	if err != nil {
		panic(err)
	}

	return instance, nil
}

func WriteServiceAccountToFile(serviceAccount string) string {
	saUuid, err := uuid.NewRandom()
	Check(err)

	saPath := "/tmp/%s" + saUuid.String()

	err = ioutil.WriteFile(saPath, []byte(serviceAccount), 0644)
	Check(err)

	return saPath
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
