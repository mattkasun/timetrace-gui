package database

import (
	"encoding/json"
	"errors"

	"github.com/mattkasun/timetrace-gui/models"
)

func SaveProject(p *models.Project) error {
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return insert(p.Name, string(value), PROJECT_TABLE_NAME)
}

func GetProject(name string) (models.Project, error) {
	var project models.Project
	records, err := fetch(PROJECT_TABLE_NAME)
	if err != nil {
		return project, err
	}
	for key, record := range records {
		if key == name {
			if err := json.Unmarshal([]byte(record), &project); err != nil {
				return project, err
			}
			return project, nil
		}
	}
	return project, errors.New("no such project")
}

func GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	var project models.Project
	records, err := fetch(PROJECT_TABLE_NAME)
	if err != nil {
		return projects, err
	}
	for _, record := range records {
		if err := json.Unmarshal([]byte(record), &project); err != nil {
			continue
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func DeleteProject(name string) error {
	return delete(name, PROJECT_TABLE_NAME)
}
