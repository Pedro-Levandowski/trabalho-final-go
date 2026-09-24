package repositories

import (
	"errors"
	"sync"

	"api-gin/models"
)

var (
	ErrTurmaJaAlocada      = errors.New("turma já possui uma alocação")
	ErrConflitoHorarioSala = errors.New("a sala já está ocupada nesse dia e horário")
)

type AlocacaoRepository struct {
	mu        sync.RWMutex
	alocacoes []models.Alocacao
}

func NovoAlocacaoRepository() *AlocacaoRepository {
	return &AlocacaoRepository{alocacoes: make([]models.Alocacao, 0)}
}

func (r *AlocacaoRepository) Criar(novaAlocacao models.Alocacao) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, alocacao := range r.alocacoes {
		if alocacao.TurmaID == novaAlocacao.TurmaID {
			return ErrTurmaJaAlocada
		}

		mesmaSalaEDia := alocacao.SalaID == novaAlocacao.SalaID && alocacao.DiaSemana == novaAlocacao.DiaSemana
		horariosSobrepostos := models.HorariosSobrepostos(
			novaAlocacao.InicioMinutos,
			novaAlocacao.FimMinutos,
			alocacao.InicioMinutos,
			alocacao.FimMinutos,
		)
		if mesmaSalaEDia && horariosSobrepostos {
			return ErrConflitoHorarioSala
		}
	}

	r.alocacoes = append(r.alocacoes, novaAlocacao)
	return nil
}

func (r *AlocacaoRepository) ListarPorSala(salaID string) []models.Alocacao {
	r.mu.RLock()
	defer r.mu.RUnlock()

	alocacoes := make([]models.Alocacao, 0)
	for _, alocacao := range r.alocacoes {
		if alocacao.SalaID == salaID {
			alocacoes = append(alocacoes, alocacao)
		}
	}

	return alocacoes
}

func (r *AlocacaoRepository) BuscarPorTurma(turmaID string) (models.Alocacao, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, alocacao := range r.alocacoes {
		if alocacao.TurmaID == turmaID {
			return alocacao, true
		}
	}

	return models.Alocacao{}, false
}

func (r *AlocacaoRepository) BuscarConflitos(salaID string, diaSemana models.DiaSemana, inicioMinutos, fimMinutos int) []models.Alocacao {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conflitos := make([]models.Alocacao, 0)
	for _, alocacao := range r.alocacoes {
		mesmaSalaEDia := alocacao.SalaID == salaID && alocacao.DiaSemana == diaSemana
		horariosSobrepostos := models.HorariosSobrepostos(inicioMinutos, fimMinutos, alocacao.InicioMinutos, alocacao.FimMinutos)
		if mesmaSalaEDia && horariosSobrepostos {
			conflitos = append(conflitos, alocacao)
		}
	}

	return conflitos
}
