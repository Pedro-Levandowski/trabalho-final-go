package repositories

import (
	"sync"

	"api-gin/models"
)

type AlunoRepository struct {
	mu     sync.RWMutex
	alunos []models.Aluno
}

func NovoAlunoRepository() *AlunoRepository {
	return &AlunoRepository{alunos: make([]models.Aluno, 0)}
}

func (r *AlunoRepository) Criar(novoAluno models.Aluno) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, aluno := range r.alunos {
		if aluno.ID == novoAluno.ID {
			return false
		}
	}

	r.alunos = append(r.alunos, novoAluno)
	return true
}

func (r *AlunoRepository) Listar() []models.Aluno {
	r.mu.RLock()
	defer r.mu.RUnlock()

	alunos := make([]models.Aluno, len(r.alunos))
	copy(alunos, r.alunos)
	return alunos
}

func (r *AlunoRepository) BuscarPorID(id string) (models.Aluno, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, aluno := range r.alunos {
		if aluno.ID == id {
			return aluno, true
		}
	}

	return models.Aluno{}, false
}
