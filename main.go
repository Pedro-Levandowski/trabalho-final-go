package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"api-gin/models"
)

const versaoAPI = "1.2.0"

func configurarRotas() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	salaRepository := novoSalaRepository()
	alunoRepository := novoAlunoRepository()

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
	}

	return r
}

func main() {
	r := configurarRotas()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
