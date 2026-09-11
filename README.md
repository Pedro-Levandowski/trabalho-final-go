# Especificação de Negócio: API de Controle de Turmas, Salas e Alunos (SGA - Sistema de Gestão de Alocação)

Este documento apresenta a especificação funcional e de negócios para o Sistema de Gestão de Alocação (SGA), incorporando o gerenciamento do ente , o processo de matrícula de alunos em turmas, a alocação física de salas e a validação contínua de capacidade e agenda.

## 1. Visão Geral do Negócio

### O Problema

Instituições de ensino enfrentam desafios operacionais recorrentes relacionados ao choque de horários na alocação física de salas, à subutilização de ambientes acadêmicos, ao controle inadequado da lotação de alunos por turma e à falta de visibilidade centralizada sobre a ocupação do campus e das matrículas. A gestão manual por planilhas isoladas gera retrabalho, inconsistências operacionais e alocações inadequadas.

### A Solução de Negócio

O  consiste em uma solução centralizada para gerenciar o inventário de espaços físicos, o cadastro de alunos e o agendamento de turmas. A solução automatiza o processo de alocação física e de gestão de matrículas através de regras de validação de capacidade máxima, duplicidade e prevenção de conflitos temporais ().

## 2. Escopo Funcional e Regras de Negócio

### 2.1. Gestão do Inventário de Salas

- O sistema deve registrar salas de aula e laboratórios com identificador único, nome/descrição, capacidade máxima de alunos e lista de recursos disponíveis (ex.: projetores, computadores, sistemas de áudio, ar-condicionado).
- As características físicas devem ser mantidas para assegurar alocações alinhadas com as exigências técnicas das disciplinas.
- O sistema deve permitir listar as salas registradas e consultar a grade de uso para verificação de disponibilidade em períodos específicos.

### 2.2. Gestão de Alunos

- O sistema deve permitir o registro de alunos com identificador único (ex.: número de matrícula), nome completo e e-mail institucional.
- O sistema deve permitir a listagem de todos os alunos cadastrados e a busca detalhada pelo identificador individual.

### 2.3. Gestão de Turmas e Matrícula de Alunos

- O sistema deve registrar turmas acadêmicas contendo identificador único, nome da turma, disciplina associada e docente responsável.
- Uma turma pode conter múltiplos alunos matriculados. O sistema deve registrar a associação entre o aluno e a turma.
-  Um aluno não pode ser adicionado mais de uma vez na mesma turma.
-  Se a turma já estiver alocada em uma sala, a inclusão de um novo aluno não pode fazer com que o total de alunos matriculados ultrapasse a capacidade máxima da sala alocada.
- O sistema deve permitir associar uma turma a um espaço físico definindo o dia da semana, o horário de início e o horário de término da aula.

### 2.4. Motor de Validação e Regras de Alocação ()

Toda tentativa de alocação de uma turma em uma sala deve ser submetida a um motor de regras de negócio. A alocação só será confirmada se atender a todas as condições abaixo:

- A sala e a turma informadas devem estar previamente cadastradas e ativas no sistema.
- A quantidade total de alunos matriculados na turma não pode exceder a capacidade máxima de assentos da sala solicitada.
- Uma sala não pode receber duas ou mais turmas no mesmo dia da semana em horários sobrepostos.
- Considera-se haver sobreposição de horários quando o horário de início de uma nova solicitação for menor que o horário de término de uma alocação existente,  o horário de término da nova solicitação for maior que o horário de início da alocação existente para a mesma sala e mesmo dia da semana.
- Um aluno não pode estar matriculado em turmas cujos horários de aula coincidam no mesmo dia da semana e período.

### 2.5. Monitoramento Operacional da Solução ()

- O sistema deve prover uma verificação contínua de disponibilidade funcional que informe o status operacional da solução, o horário atual do servidor e a versão ativa do serviço, garantindo confiabilidade no consumo da API.

## 3\. Contratos de Interface de Negócio (API REST)

As funcionalidades da solução são expostas por meio de serviços padronizados:

-  Verificar a disponibilidade funcional do serviço.
-  Retorna o status de integridade do sistema, data/hora atual e versão ativa.
-  Registrar uma nova sala no inventário físico.
-  Identificador da sala, nome, capacidade máxima (maior que zero) e lista de recursos técnicos.
-  Consultar o inventário completo de salas cadastradas.
-  Registrar um novo aluno na instituição.
-  Identificador do aluno (matrícula), nome completo e e-mail.
-  Consultar a listagem de alunos cadastrados.
-  Registrar uma nova turma no sistema.
-  Identificador da turma, nome, disciplina e professor responsável.
-  Consultar as turmas cadastradas, exibindo a quantidade atual de alunos matriculados e o status de alocação física.
-  Matricular/Adicionar um aluno em uma turma específica.
-  Identificador do aluno (aluno\_id).

  - Se a turma ou o aluno não forem encontrados: Erro de recurso não localizado (404).
  - Se o aluno já estiver matriculado na turma: Erro de duplicidade (409).
  - Se a turma já estiver alocada em uma sala e a inclusão exceder a capacidade da sala: Erro de capacidade insuficiente (422).
  - Se houver sobreposição de horário com outra turma em que o aluno esteja matriculado: Erro de conflito de agenda do aluno (409).

-  Listar todos os alunos matriculados em uma turma específica.
-  Realizar a alocação de uma sala para uma turma específica.
-  Identificador da sala, dia da semana, horário de início (HH\:MM) e horário de término (HH\:MM).

  - Se a sala não existir ou a turma não for encontrada: Erro de recurso não localizado (404).
  - Se a capacidade da sala for menor que a quantidade atual de alunos matriculados na turma: Erro de capacidade insuficiente (422).
  - Se houver sobreposição de horário no mesmo dia e sala: Erro de conflito de agenda () (409).
  - Se todas as regras forem atendidas: Sucesso na alocação da turma.
