package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"api-gin/models"
)

const versaoAPI = "1.1.0"

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

			if cadastrada := repository.criar(novaSala); !cadastrada {
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
