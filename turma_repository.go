package main

import (
	"errors"
	"sync"

	"api-gin/models"
)

var (
	ErrTurmaNaoEncontrada = errors.New("turma não encontrada")
	ErrAlunoJaMatriculado = errors.New("aluno já matriculado na turma")
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

	novaTurma.AlunosIDs = make([]string, 0)
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

func (r *TurmaRepository) buscarPorID(id string) (models.Turma, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, turma := range r.turmas {
		if turma.ID == id {
			return turma, true
		}
	}

	return models.Turma{}, false
}

func (r *TurmaRepository) marcarComoAlocada(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for indice := range r.turmas {
		if r.turmas[indice].ID == id {
			r.turmas[indice].Alocada = true
			return nil
		}
	}

	return ErrTurmaNaoEncontrada
}

func (r *TurmaRepository) matricularAluno(turmaID, alunoID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for indice := range r.turmas {
		if r.turmas[indice].ID != turmaID {
			continue
		}

		for _, matriculadoID := range r.turmas[indice].AlunosIDs {
			if matriculadoID == alunoID {
				return ErrAlunoJaMatriculado
			}
		}

		r.turmas[indice].AlunosIDs = append(r.turmas[indice].AlunosIDs, alunoID)
		r.turmas[indice].QuantidadeAlunos = len(r.turmas[indice].AlunosIDs)
		return nil
	}

	return ErrTurmaNaoEncontrada
}

func (r *TurmaRepository) listarAlunosIDs(turmaID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, turma := range r.turmas {
		if turma.ID == turmaID {
			alunosIDs := make([]string, len(turma.AlunosIDs))
			copy(alunosIDs, turma.AlunosIDs)
			return alunosIDs, nil
		}
	}

	return nil, ErrTurmaNaoEncontrada
}
