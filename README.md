# 🌀 Magalu Cloud - Pulse Ingestor

Este projeto implementa o componente **Ingestor** do desafio técnico do Magalu Cloud. Ele é responsável por consumir mensagens de uso de produtos (pulses) via RabbitMQ, agregá-las temporariamente no Redis e posteriormente persistir os dados no PostgreSQL por meio de um cron worker.

---

## 📚 Sumário

- [🧠 Visão Geral](#-visão-geral)
- [⚙️ Tecnologias Utilizadas](#-tecnologias-utilizadas)
- [📁 Estrutura do Projeto](#-estrutura-do-projeto)
- [🚀 Executando o Projeto](#-executando-o-projeto)
- [🧪 Executando os Testes](#-executando-os-testes)
- [📦 Simulador de Eventos (Publisher)](#-simulador-de-eventos-publisher)
- [🔁 Fluxo do Sistema](#-fluxo-do-sistema)
- [🧱 Arquitetura](#-arquitetura)
- [📌 Considerações Finais](#-considerações-finais)

---

## 🧠 Visão Geral

Este serviço é composto por dois principais workers:

- `PulseTask` (Ingestor): Consome mensagens do RabbitMQ, agrega os dados por chave composta (tenant, SKU, unidade) e armazena temporariamente no Redis.
- `SavePulseTask`: A cada 1 hora, busca os dados agregados no Redis e persiste no PostgreSQL, limpando os dados do Redis após a persistência com sucesso.

---

## ⚙️ Tecnologias Utilizadas

- **Go 1.23+**
- **RabbitMQ**
- **Redis**
- **PostgreSQL**
- **Docker & Docker Compose**
- **Logrus** (logger estruturado)
- **Gomock** (testes unitários)

---

## 📁 Estrutura do Projeto

```bash
.
├── cmd/                   # Ponto de entrada da aplicação
│   ├── worker/            # Configurações dos workers (worker e cron worker)
│   ├── simulator/         # Configurações do simulator (Popula a fila com dados fakes e simula alto volume de dados)
│   ├── api/               # Configurações da api
├── application/           # Contém a lógica de aplicação
│   ├── usecase/           # Casos de uso da aplicação, focados em operações específicas do domínio.
│   ├── service/           # Serviços que encapsulam integrações externas.
├── domain/                # Contém os contratos e estruturas centrais da aplicação
│   ├── dto/               # Define os Data Transfer Objects, usados para transportar dados entre camadas da aplicação.
│   ├── entity/            # Contém as entidades de domínio, que representam os objetos persistidos no banco de dados.   
│   ├── mapper/            # Realiza a conversão entre DTOs e Entities, garantindo o isolamento entre camadas.       
├── infrastructure/        # Implementações técnicas
│   ├── database/          # Configurações e conexões com o bancos de dados.
│   ├── repository/        # Repositórios com lógica de persistência no banco de dados.
├── internal/              # Código interno da aplicação que não deve ser exposto para uso externo.
│   ├── configs/           # Leitura e organização das variáveis de ambiente. 
│   ├── constants/         # Constantes reutilizadas no projeto.
│   ├── dependency/        # Injeção de dependências e construção de componentes da aplicação.
├── utils/                 # Funções utilitárias
├── vendor/                # Pacotes externos baixados pelo Go Modules (go mod vendor), caso use vendoring.
├── docker-compose.yml     # Subida dos serviços
├── go.mod
└── README.md
```

---

## 🚀 Executando o Projeto

1. Clone o repositório

```
    git clone https://github.com/seuusuario/pulse-ingestor.git
    cd pulse-ingestor
```

2. Suba os serviços com Docker Compose

```
    docker-compose up --build
```

Esse comando irá iniciar:

- Redis
- RabbitMQ
- PostgreSQL
- Aplicação Go (cmd/main.go)



---
