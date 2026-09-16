package main

import (
	"sync"

	"api-gin/models"
)

type TurmaRepository struct {
	mu     sync.RWMutex
	turmas []models.Turma
}

func novoTurmaRepository() *TurmaRepository {
	return &TurmaRepository{turmas: make([]models.Turma, 0)}
}

func (r *TurmaRepository) criar(novaTurma models.Turma) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, turma := range r.turmas {
		if turma.ID == novaTurma.ID {
			return false
		}
	}

	r.turmas = append(r.turmas, novaTurma)
	return true
}

func (r *TurmaRepository) listar() []models.Turma {
	r.mu.RLock()
	defer r.mu.RUnlock()

	turmas := make([]models.Turma, len(r.turmas))
	copy(turmas, r.turmas)
	return turmas
}
