package models

type Turma struct {
	ID               string   `json:"id"`
	Nome             string   `json:"nome"`
	Disciplina       string   `json:"disciplina"`
	Professor        string   `json:"professor"`
	QuantidadeAlunos int      `json:"quantidade_alunos"`
	Alocada          bool     `json:"alocada"`
	AlunosIDs        []string `json:"-"`
}
