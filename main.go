package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RecursoSala string

const (
	RecursoProjetor       RecursoSala = "projetor"
	RecursoComputadores   RecursoSala = "computadores"
	RecursoSistemaAudio   RecursoSala = "sistema_audio"
	RecursoArCondicionado RecursoSala = "ar_condicionado"
)

type Sala struct {
	ID         string        `json:"id" binding:"required"`
	Nome       string        `json:"nome" binding:"required"`
	Capacidade int           `json:"capacidade" binding:"required,gt=0"`
	Recursos   []RecursoSala `json:"recursos" binding:"omitempty,dive,oneof=projetor computadores sistema_audio ar_condicionado"`
}

type SalaRepository struct {
	mu    sync.RWMutex
	salas []Sala
}

func novoSalaRepository() *SalaRepository {
	return &SalaRepository{salas: make([]Sala, 0)}
}

func (r *SalaRepository) criar(novaSala Sala) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, sala := range r.salas {
		if sala.ID == novaSala.ID {
			return false
		}
	}

	if novaSala.Recursos == nil {
		novaSala.Recursos = []RecursoSala{}
	}

	r.salas = append(r.salas, novaSala)
	return true
}

func (r *SalaRepository) listar() []Sala {
	r.mu.RLock()
	defer r.mu.RUnlock()

	salas := make([]Sala, len(r.salas))
	copy(salas, r.salas)
	return salas
}

func configurarRotas() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	repository := novoSalaRepository()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		v1.POST("/salas", func(c *gin.Context) {
			var novaSala Sala
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

			if cadastrada := repository.criar(novaSala); !cadastrada {
				c.JSON(http.StatusConflict, gin.H{
					"erro": "já existe uma sala com o identificador informado",
				})
				return
			}

			if novaSala.Recursos == nil {
				novaSala.Recursos = []RecursoSala{}
			}

			c.JSON(http.StatusCreated, novaSala)
		})

		v1.GET("/salas", func(c *gin.Context) {
			c.JSON(http.StatusOK, repository.listar())
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
