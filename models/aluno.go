package models

type Aluno struct {
	ID    string `json:"id" binding:"required"`
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
