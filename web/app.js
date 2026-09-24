const painelResultado = document.querySelector("#resultado");
const titulo = document.querySelector("#resultado-titulo");
const status = document.querySelector("#resultado-status");
const conteudo = document.querySelector("#resultado-conteudo");

async function chamarAPI(descricao, caminho, opcoes = {}) {
    painelResultado.hidden = false;
    titulo.textContent = descricao;
    status.textContent = "Aguardando...";
    status.className = "";
    conteudo.textContent = `${opcoes.method || "GET"} ${caminho}`;

    try {
        const resposta = await fetch(caminho, {
            ...opcoes,
            headers: opcoes.body ? { "Content-Type": "application/json" } : undefined,
        });
        const texto = await resposta.text();
        let corpo = texto;

        try {
            corpo = JSON.stringify(JSON.parse(texto), null, 2);
        } catch {
            corpo = texto || "Resposta sem conteúdo.";
        }

        painelResultado.hidden = false;
        status.textContent = `${resposta.status} ${resposta.statusText}`;
        status.className = resposta.ok ? "sucesso" : "erro";
        conteudo.textContent = corpo;
    } catch (erro) {
        painelResultado.hidden = false;
        status.textContent = "Falha de conexão";
        status.className = "erro";
        conteudo.textContent = erro.message;
    }
}

document.querySelector("#fechar-resultado").addEventListener("click", () => {
    painelResultado.hidden = true;
});

function dados(formulario) {
    return new FormData(formulario);
}

function enviarJSON(descricao, caminho, corpo) {
    return chamarAPI(descricao, caminho, {
        method: "POST",
        body: JSON.stringify(corpo),
    });
}

document.querySelector("#verificar-api").addEventListener("click", () =>
    chamarAPI("Monitoramento da API", "/api/v1/health")
);

document.querySelector("#form-sala").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const form = evento.currentTarget;
    const valores = dados(form);
    enviarJSON("Cadastro de sala", "/api/v1/salas", {
        id: valores.get("id"),
        nome: valores.get("nome"),
        capacidade: Number(valores.get("capacidade")),
        recursos: valores.getAll("recursos"),
    });
});

document.querySelector("#listar-salas").addEventListener("click", () =>
    chamarAPI("Salas cadastradas", "/api/v1/salas")
);

document.querySelector("#form-agenda-sala").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const valores = dados(evento.currentTarget);
    const salaID = encodeURIComponent(valores.get("sala_id"));
    const consulta = new URLSearchParams({
        dia_semana: valores.get("dia_semana"),
        horario_inicio: valores.get("horario_inicio"),
        horario_fim: valores.get("horario_fim"),
    });
    chamarAPI("Disponibilidade da sala", `/api/v1/salas/${salaID}/alocacoes?${consulta}`);
});

document.querySelector("#listar-agenda").addEventListener("click", () => {
    const form = document.querySelector("#form-agenda-sala");
    const campoSala = form.querySelector('[name="sala_id"]');
    if (!campoSala.reportValidity()) return;
    const salaID = encodeURIComponent(campoSala.value);
    chamarAPI("Alocações da sala", `/api/v1/salas/${salaID}/alocacoes`);
});

document.querySelector("#form-aluno").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const valores = dados(evento.currentTarget);
    enviarJSON("Cadastro de aluno", "/api/v1/alunos", {
        id: valores.get("id"),
        nome: valores.get("nome"),
        email: valores.get("email"),
    });
});

document.querySelector("#listar-alunos").addEventListener("click", () =>
    chamarAPI("Alunos cadastrados", "/api/v1/alunos")
);

document.querySelector("#form-buscar-aluno").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const alunoID = encodeURIComponent(dados(evento.currentTarget).get("id"));
    chamarAPI("Consulta de aluno", `/api/v1/alunos/${alunoID}`);
});

document.querySelector("#form-turma").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const valores = dados(evento.currentTarget);
    enviarJSON("Cadastro de turma", "/api/v1/turmas", {
        id: valores.get("id"),
        nome: valores.get("nome"),
        disciplina: valores.get("disciplina"),
        professor: valores.get("professor"),
    });
});

document.querySelector("#listar-turmas").addEventListener("click", () =>
    chamarAPI("Turmas cadastradas", "/api/v1/turmas")
);

document.querySelector("#form-matricula").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const valores = dados(evento.currentTarget);
    const turmaID = encodeURIComponent(valores.get("turma_id"));
    enviarJSON("Matrícula de aluno", `/api/v1/turmas/${turmaID}/alunos`, {
        aluno_id: valores.get("aluno_id"),
    });
});

document.querySelector("#listar-matriculados").addEventListener("click", () => {
    const form = document.querySelector("#form-matricula");
    const campoTurma = form.querySelector('[name="turma_id"]');
    if (!campoTurma.reportValidity()) return;
    const turmaID = encodeURIComponent(campoTurma.value);
    chamarAPI("Alunos matriculados", `/api/v1/turmas/${turmaID}/alunos`);
});

document.querySelector("#form-alocacao").addEventListener("submit", (evento) => {
    evento.preventDefault();
    const valores = dados(evento.currentTarget);
    const turmaID = encodeURIComponent(valores.get("turma_id"));
    enviarJSON("Alocação de turma", `/api/v1/turmas/${turmaID}/alocar`, {
        sala_id: valores.get("sala_id"),
        dia_semana: valores.get("dia_semana"),
        horario_inicio: valores.get("horario_inicio"),
        horario_fim: valores.get("horario_fim"),
    });
});
