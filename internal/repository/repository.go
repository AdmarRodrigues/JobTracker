package repository

import (
	"JobTracker/internal/models"
	"sync"
)

type JobTrack struct {
	mu     sync.Mutex
	jb     map[int]models.Jobs
	NextId int
}

func NewJobStore() *JobTrack {
	return &JobTrack{jb: make(map[int]models.Jobs), NextId: 1}
}

func (j *JobTrack) Get(id int) (models.Jobs, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	job, exists := j.jb[id]
	return job, exists
}

func (j *JobTrack) Add(title string) models.Jobs {
	j.mu.Lock()
	defer j.mu.Unlock()
	currentId := j.NextId

	newJobTest := models.Jobs{Cargo: title}

	j.jb[currentId] = newJobTest
	j.NextId++
	return j.jb[currentId]

}

func (j *JobTrack) ListAll() []models.Jobs {
	j.mu.Lock()
	defer j.mu.Unlock()

	all := make([]models.Jobs, len(j.jb))
	for _, job := range j.jb {
		all = append(all, job)
	}
	return all

}
