package models

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	ID      uuid.UUID
	Project string
	User    uuid.UUID
	Start   time.Time
	End     time.Time
}

func (r *Record) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

type ReportData struct {
	Project string
	Records []Record
	Sum     time.Duration
}

func ConvertToReport(records []Record) []ReportData {
	data := make(map[string][]Record)
	reportData := []ReportData{}
	project := ReportData{}
	for _, record := range records {
		data[record.Project] = append(data[record.Project], record)
	}
	sum := time.Duration(0)
	for k, v := range data {
		project.Project = k
		for _, item := range v {
			project.Records = append(project.Records, item)
			timeSpent := item.End.Sub(item.Start)
			sum = sum + timeSpent
		}
		project.Sum = sum
		reportData = append(reportData, project)
	}
	return reportData
}
