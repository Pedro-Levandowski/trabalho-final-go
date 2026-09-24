package main

import (
	"errors"

	"api-gin/models"
)

var ErrConflitoAgendaAluno = errors.New("aluno possui conflito de agenda com outra turma")

func alunoPossuiConflitoAgenda(
	alunoID string,
	turmaDestinoID string,
	alocacaoDestino models.Alocacao,
	turmaRepository *TurmaRepository,
	alocacaoRepository *AlocacaoRepository,
) bool {
	for _, turma := range turmaRepository.listarPorAluno(alunoID) {
		if turma.ID == turmaDestinoID {
			continue
		}

		alocacaoExistente, alocada := alocacaoRepository.buscarPorTurma(turma.ID)
		if !alocada || alocacaoExistente.DiaSemana != alocacaoDestino.DiaSemana {
			continue
		}

		if existeSobreposicao(
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
