package main

import (
	"errors"

	"api-gin/models"
	"api-gin/repositories"
)

var ErrConflitoAgendaAluno = errors.New("aluno possui conflito de agenda com outra turma")

func alunoPossuiConflitoAgenda(
	alunoID string,
	turmaDestinoID string,
	alocacaoDestino models.Alocacao,
	turmaRepository *repositories.TurmaRepository,
	alocacaoRepository *repositories.AlocacaoRepository,
) bool {
	for _, turma := range turmaRepository.ListarPorAluno(alunoID) {
		if turma.ID == turmaDestinoID {
			continue
		}

		alocacaoExistente, alocada := alocacaoRepository.BuscarPorTurma(turma.ID)
		if !alocada || alocacaoExistente.DiaSemana != alocacaoDestino.DiaSemana {
			continue
		}

		if models.HorariosSobrepostos(
			alocacaoDestino.InicioMinutos,
			alocacaoDestino.FimMinutos,
			alocacaoExistente.InicioMinutos,
			alocacaoExistente.FimMinutos,
		) {
			return true
		}
	}

	return false
}
