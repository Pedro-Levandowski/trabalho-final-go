package repositories

import (
	"sync"

	"api-gin/models"
)

type SalaRepository struct {
	mu    sync.RWMutex
	salas []models.Sala
}

func NovoSalaRepository() *SalaRepository {
	return &SalaRepository{salas: make([]models.Sala, 0)}
}

func (r *SalaRepository) Criar(novaSala models.Sala) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, sala := range r.salas {
		if sala.ID == novaSala.ID {
			return false
		}
	}

	if novaSala.Recursos == nil {
		novaSala.Recursos = []models.RecursoSala{}
	}

	r.salas = append(r.salas, novaSala)
	return true
}

func (r *SalaRepository) Listar() []models.Sala {
	r.mu.RLock()
	defer r.mu.RUnlock()

	salas := make([]models.Sala, len(r.salas))
	copy(salas, r.salas)
	return salas
}

func (r *SalaRepository) BuscarPorID(id string) (models.Sala, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, sala := range r.salas {
		if sala.ID == id {
			return sala, true
		}
	}

	return models.Sala{}, false
}
