package repository

import (
	"JobTracker/internal/models"
	"fmt"
	"log"
	"sync"
	"time"
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

func (j *JobTrack) Add(cargo string, empresa string, status string) models.Jobs {
	j.mu.Lock()
	defer j.mu.Unlock()

	currentId := j.NextId
	j.jb[currentId] = models.Jobs{Id: currentId, Cargo: cargo, Empresa: empresa, Status: status, Data: time.Now().Format("2006-01-02")}
	//j.jb[currentId] = newJobTest
	log.Printf("criando map com chave %d, e id %d", currentId, currentId)

	j.NextId++
	return j.jb[currentId]

}

func (j *JobTrack) ListAll() []models.Jobs {
	j.mu.Lock()
	defer j.mu.Unlock()

	all := make([]models.Jobs, 0, len(j.jb))
	for _, job := range j.jb {
		all = append(all, job)
	}
	fmt.Println(all)
	return all

}

func (j *JobTrack) DeleteById(id int) {
	j.mu.Lock()
	defer j.mu.Unlock()

	dissapear := j.jb[id]
	delete(j.jb, dissapear.Id)

}
