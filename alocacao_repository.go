package main

import (
	"errors"
	"sync"
	"time"

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

func novoAlocacaoRepository() *AlocacaoRepository {
	return &AlocacaoRepository{alocacoes: make([]models.Alocacao, 0)}
}

func (r *AlocacaoRepository) criar(novaAlocacao models.Alocacao) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, alocacao := range r.alocacoes {
		if alocacao.TurmaID == novaAlocacao.TurmaID {
			return ErrTurmaJaAlocada
		}

		mesmaSalaEDia := alocacao.SalaID == novaAlocacao.SalaID && alocacao.DiaSemana == novaAlocacao.DiaSemana
		horariosSobrepostos := novaAlocacao.InicioMinutos < alocacao.FimMinutos && novaAlocacao.FimMinutos > alocacao.InicioMinutos
		if mesmaSalaEDia && horariosSobrepostos {
			return ErrConflitoHorarioSala
		}
	}

	r.alocacoes = append(r.alocacoes, novaAlocacao)
	return nil
}

func horarioEmMinutos(horario string) (int, error) {
	valor, err := time.Parse("15:04", horario)
	if err != nil {
		return 0, err
	}

	return valor.Hour()*60 + valor.Minute(), nil
}
