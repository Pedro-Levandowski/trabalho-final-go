package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"api-gin/models"
)

const versaoAPI = "1.8.0"

type criarTurmaRequest struct {
	ID         string `json:"id" binding:"required"`
	Nome       string `json:"nome" binding:"required"`
	Disciplina string `json:"disciplina" binding:"required"`
	Professor  string `json:"professor" binding:"required"`
}

type matricularAlunoRequest struct {
	AlunoID string `json:"aluno_id" binding:"required"`
}

type alocarSalaRequest struct {
	SalaID        string           `json:"sala_id" binding:"required"`
	DiaSemana     models.DiaSemana `json:"dia_semana" binding:"required,oneof=segunda terca quarta quinta sexta sabado domingo"`
	HorarioInicio string           `json:"horario_inicio" binding:"required"`
	HorarioFim    string           `json:"horario_fim" binding:"required"`
}

func configurarRotas() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	salaRepository := novoSalaRepository()
	alunoRepository := novoAlunoRepository()
	turmaRepository := novoTurmaRepository()
	alocacaoRepository := novoAlocacaoRepository()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   versaoAPI,
			})
		})

		v1.POST("/salas", func(c *gin.Context) {
			var novaSala models.Sala
			if err := c.ShouldBindJSON(&novaSala); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "dados inválidos: informe id, nome, capacidade maior que zero e recursos permitidos",
				})
				return
			}

			novaSala.ID = strings.TrimSpace(novaSala.ID)
			novaSala.Nome = strings.TrimSpace(novaSala.Nome)
			if novaSala.ID == "" || novaSala.Nome == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "id e nome não podem conter apenas espaços",
				})
				return
			}
			novaSala.Ativa = true

			if cadastrada := salaRepository.criar(novaSala); !cadastrada {
				c.JSON(http.StatusConflict, gin.H{
					"erro": "já existe uma sala com o identificador informado",
				})
				return
			}

			if novaSala.Recursos == nil {
				novaSala.Recursos = []models.RecursoSala{}
			}

			c.JSON(http.StatusCreated, novaSala)
		})

		v1.GET("/salas", func(c *gin.Context) {
			c.JSON(http.StatusOK, salaRepository.listar())
		})

		v1.GET("/salas/:id/alocacoes", func(c *gin.Context) {
			sala, encontrada := salaRepository.buscarPorID(c.Param("id"))
			if !encontrada {
				c.JSON(http.StatusNotFound, gin.H{"erro": "sala não encontrada"})
				return
			}

			diaValor, temDia := c.GetQuery("dia_semana")
			inicioValor, temInicio := c.GetQuery("horario_inicio")
			fimValor, temFim := c.GetQuery("horario_fim")

			if !temDia && !temInicio && !temFim {
				c.JSON(http.StatusOK, gin.H{
					"sala_id":   sala.ID,
					"alocacoes": alocacaoRepository.listarPorSala(sala.ID),
				})
				return
			}

			if !temDia || !temInicio || !temFim {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "para consultar a disponibilidade, informe dia_semana, horario_inicio e horario_fim",
				})
				return
			}

			diaSemana := models.DiaSemana(strings.TrimSpace(diaValor))
			inicioValor = strings.TrimSpace(inicioValor)
			fimValor = strings.TrimSpace(fimValor)
			inicioMinutos, errInicio := horarioEmMinutos(inicioValor)
			fimMinutos, errFim := horarioEmMinutos(fimValor)
			if !diaSemana.Valido() || errInicio != nil || errFim != nil || fimMinutos <= inicioMinutos {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "dia da semana ou horários inválidos",
				})
				return
			}

			conflitos := alocacaoRepository.buscarConflitos(sala.ID, diaSemana, inicioMinutos, fimMinutos)
			c.JSON(http.StatusOK, gin.H{
				"sala_id":        sala.ID,
				"dia_semana":     diaSemana,
				"horario_inicio": inicioValor,
				"horario_fim":    fimValor,
				"disponivel":     len(conflitos) == 0,
				"conflitos":      conflitos,
			})
		})

		v1.POST("/alunos", func(c *gin.Context) {
			var novoAluno models.Aluno
			if err := c.ShouldBindJSON(&novoAluno); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "dados inválidos: informe id, nome completo e e-mail válido",
				})
				return
			}

			novoAluno.ID = strings.TrimSpace(novoAluno.ID)
			novoAluno.Nome = strings.TrimSpace(novoAluno.Nome)
			novoAluno.Email = strings.ToLower(strings.TrimSpace(novoAluno.Email))
			if novoAluno.ID == "" || novoAluno.Nome == "" || novoAluno.Email == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "id, nome completo e e-mail não podem conter apenas espaços",
				})
				return
			}

			if cadastrado := alunoRepository.criar(novoAluno); !cadastrado {
				c.JSON(http.StatusConflict, gin.H{
					"erro": "já existe um aluno com o identificador informado",
				})
				return
			}

			c.JSON(http.StatusCreated, novoAluno)
		})

		v1.GET("/alunos", func(c *gin.Context) {
			c.JSON(http.StatusOK, alunoRepository.listar())
		})

		v1.GET("/alunos/:id", func(c *gin.Context) {
			aluno, encontrado := alunoRepository.buscarPorID(c.Param("id"))
			if !encontrado {
				c.JSON(http.StatusNotFound, gin.H{
					"erro": "aluno não encontrado",
				})
				return
			}

			c.JSON(http.StatusOK, aluno)
		})

		v1.POST("/turmas", func(c *gin.Context) {
			var request criarTurmaRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "dados inválidos: informe id, nome, disciplina e professor",
				})
				return
			}

			request.ID = strings.TrimSpace(request.ID)
			request.Nome = strings.TrimSpace(request.Nome)
			request.Disciplina = strings.TrimSpace(request.Disciplina)
			request.Professor = strings.TrimSpace(request.Professor)
			if request.ID == "" || request.Nome == "" || request.Disciplina == "" || request.Professor == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "id, nome, disciplina e professor não podem conter apenas espaços",
				})
				return
			}

			novaTurma := models.Turma{
				ID:         request.ID,
				Nome:       request.Nome,
				Disciplina: request.Disciplina,
				Professor:  request.Professor,
				Ativa:      true,
			}

			if cadastrada := turmaRepository.criar(novaTurma); !cadastrada {
				c.JSON(http.StatusConflict, gin.H{
					"erro": "já existe uma turma com o identificador informado",
				})
				return
			}

			c.JSON(http.StatusCreated, novaTurma)
		})

		v1.GET("/turmas", func(c *gin.Context) {
			c.JSON(http.StatusOK, turmaRepository.listar())
		})

		v1.POST("/turmas/:id/alocar", func(c *gin.Context) {
			var request alocarSalaRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "informe sala_id, dia_semana e horários no formato HH:MM",
				})
				return
			}

			request.SalaID = strings.TrimSpace(request.SalaID)
			request.HorarioInicio = strings.TrimSpace(request.HorarioInicio)
			request.HorarioFim = strings.TrimSpace(request.HorarioFim)

			inicioMinutos, errInicio := horarioEmMinutos(request.HorarioInicio)
			fimMinutos, errFim := horarioEmMinutos(request.HorarioFim)
			if request.SalaID == "" || errInicio != nil || errFim != nil || fimMinutos <= inicioMinutos {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "sala_id deve ser informado e o horário final deve ser posterior ao inicial no formato HH:MM",
				})
				return
			}

			turma, encontrada := turmaRepository.buscarPorID(c.Param("id"))
			if !encontrada {
				c.JSON(http.StatusNotFound, gin.H{"erro": "turma não encontrada"})
				return
			}

			sala, encontrada := salaRepository.buscarPorID(request.SalaID)
			if !encontrada {
				c.JSON(http.StatusNotFound, gin.H{"erro": "sala não encontrada"})
				return
			}

			if !turma.Ativa || !sala.Ativa {
				c.JSON(http.StatusConflict, gin.H{"erro": "a turma e a sala devem estar ativas"})
				return
			}

			if turma.QuantidadeAlunos > sala.Capacidade {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"erro": "a capacidade da sala é menor que a quantidade de alunos da turma",
				})
				return
			}

			alocacao := models.Alocacao{
				TurmaID:       turma.ID,
				SalaID:        sala.ID,
				DiaSemana:     request.DiaSemana,
				HorarioInicio: request.HorarioInicio,
				HorarioFim:    request.HorarioFim,
				InicioMinutos: inicioMinutos,
				FimMinutos:    fimMinutos,
			}

			for _, alunoID := range turma.AlunosIDs {
				if alunoPossuiConflitoAgenda(alunoID, turma.ID, alocacao, turmaRepository, alocacaoRepository) {
					c.JSON(http.StatusConflict, gin.H{"erro": ErrConflitoAgendaAluno.Error()})
					return
				}
			}

			if err := alocacaoRepository.criar(alocacao); err != nil {
				switch {
				case errors.Is(err, ErrTurmaJaAlocada), errors.Is(err, ErrConflitoHorarioSala):
					c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível alocar a sala"})
				}
				return
			}

			if err := turmaRepository.marcarComoAlocada(turma.ID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível atualizar a turma"})
				return
			}

			c.JSON(http.StatusCreated, alocacao)
		})

		v1.POST("/turmas/:id/alunos", func(c *gin.Context) {
			var request matricularAlunoRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "informe o identificador do aluno",
				})
				return
			}

			request.AlunoID = strings.TrimSpace(request.AlunoID)
			if request.AlunoID == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "o identificador do aluno não pode conter apenas espaços",
				})
				return
			}

			aluno, encontrado := alunoRepository.buscarPorID(request.AlunoID)
			if !encontrado {
				c.JSON(http.StatusNotFound, gin.H{
					"erro": "aluno não encontrado",
				})
				return
			}

			turma, encontrada := turmaRepository.buscarPorID(c.Param("id"))
			if !encontrada {
				c.JSON(http.StatusNotFound, gin.H{"erro": "turma não encontrada"})
				return
			}

			for _, alunoID := range turma.AlunosIDs {
				if alunoID == request.AlunoID {
					c.JSON(http.StatusConflict, gin.H{"erro": ErrAlunoJaMatriculado.Error()})
					return
				}
			}

			if alocacao, alocada := alocacaoRepository.buscarPorTurma(turma.ID); alocada {
				sala, encontrada := salaRepository.buscarPorID(alocacao.SalaID)
				if !encontrada {
					c.JSON(http.StatusInternalServerError, gin.H{"erro": "sala da alocação não encontrada"})
					return
				}

				if turma.QuantidadeAlunos >= sala.Capacidade {
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"erro": "a matrícula excederia a capacidade da sala alocada",
					})
					return
				}

				if alunoPossuiConflitoAgenda(request.AlunoID, turma.ID, alocacao, turmaRepository, alocacaoRepository) {
					c.JSON(http.StatusConflict, gin.H{"erro": ErrConflitoAgendaAluno.Error()})
					return
				}
			}

			if err := turmaRepository.matricularAluno(c.Param("id"), request.AlunoID); err != nil {
				switch {
				case errors.Is(err, ErrTurmaNaoEncontrada):
					c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
				case errors.Is(err, ErrAlunoJaMatriculado):
					c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível matricular o aluno"})
				}
				return
			}

			c.JSON(http.StatusCreated, aluno)
		})

		v1.GET("/turmas/:id/alunos", func(c *gin.Context) {
			alunosIDs, err := turmaRepository.listarAlunosIDs(c.Param("id"))
			if errors.Is(err, ErrTurmaNaoEncontrada) {
				c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
				return
			}

			alunos := make([]models.Aluno, 0, len(alunosIDs))
			for _, alunoID := range alunosIDs {
				if aluno, encontrado := alunoRepository.buscarPorID(alunoID); encontrado {
					alunos = append(alunos, aluno)
				}
			}

			c.JSON(http.StatusOK, alunos)
		})
	}

	return r
}

func main() {
	r := configurarRotas()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
