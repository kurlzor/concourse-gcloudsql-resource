package commands

import (
	"os/exec"
	"fmt"
	"encoding/json"
	"github.com/Evaneos/concourse-gcloudsql-resource/models"
	"errors"
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
	json.Unmarshal(output, instances)

	return instances, err
}

func ActivateServiceAccount(serviceAccount string) error {
	args := []string {
		"auth",
		"activate-service-account",
		fmt.Sprintf("--key-file=%s", serviceAccount),
		"--quiet",
	}

	_, err := runCommand(exec.Command("gcloud", args...))

	return err
}

func GetInstanceInfo(name string, project string) (*models.GCloudSQLInstance, error) {
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
		return nil, err
	}

	var instance *models.GCloudSQLInstance
	json.Unmarshal(output, instance)

	return instance, nil
}
