package models

type DiaSemana string

const (
	Segunda DiaSemana = "segunda"
	Terca   DiaSemana = "terca"
	Quarta  DiaSemana = "quarta"
	Quinta  DiaSemana = "quinta"
	Sexta   DiaSemana = "sexta"
	Sabado  DiaSemana = "sabado"
	Domingo DiaSemana = "domingo"
)

func (d DiaSemana) Valido() bool {
	switch d {
	case Segunda, Terca, Quarta, Quinta, Sexta, Sabado, Domingo:
		return true
	default:
		return false
	}
}

func HorariosSobrepostos(inicioA, fimA, inicioB, fimB int) bool {
	return inicioA < fimB && fimA > inicioB
}

type Alocacao struct {
	TurmaID       string    `json:"turma_id"`
	SalaID        string    `json:"sala_id"`
	DiaSemana     DiaSemana `json:"dia_semana"`
	HorarioInicio string    `json:"horario_inicio"`
	HorarioFim    string    `json:"horario_fim"`
	InicioMinutos int       `json:"-"`
	FimMinutos    int       `json:"-"`
}
