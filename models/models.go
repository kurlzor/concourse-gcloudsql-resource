package models

import (
	"errors"
	"time"
)

type Source struct {
	Project *string  		`json:"project"`
	ServiceAccount *string 	`json:"service_account"`
}

type Version struct {
	Instance string `json:"instance"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
}

func (source Source) Validate() error {
	if source.Project == nil {
		return errors.New("must set 'project'")
	}

	if source.ServiceAccount == nil {
		return errors.New("must set 'service_account'")
	}

	return nil
}

type GCloudSQLInstance struct {
	Name string `json:"instance"`
	instanceType string
	region string
	creationTime time.Time `json:"serverCaCert:createTime"`
}

type GCloudSQLInstanceList []GCloudSQLInstance

func (list GCloudSQLInstanceList) Len() int {
	return len(list)
}

func (list GCloudSQLInstanceList) Less(i, j int) bool {
	return list[i].creationTime.Before(list[j].creationTime)
}

func (list GCloudSQLInstanceList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}