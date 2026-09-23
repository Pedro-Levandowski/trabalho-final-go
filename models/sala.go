package models

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
	Ativa      bool          `json:"ativa"`
}
