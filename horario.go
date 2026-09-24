package main

import "time"

func horarioEmMinutos(horario string) (int, error) {
	valor, err := time.Parse("15:04", horario)
	if err != nil {
		return 0, err
	}

	return valor.Hour()*60 + valor.Minute(), nil
}
