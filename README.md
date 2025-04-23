> 📘 Este README está em **Português-BR** pois o desafio foi entregue em português. Caso queira uma versão em inglês, me avise!

![Go](https://img.shields.io/badge/Go-1.23-blue)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)

# 🌀 Magalu Cloud - Ingestor

Este projeto implementa o componente **Ingestor** do desafio técnico do Magalu Cloud. Ele é responsável por consumir mensagens de uso de produtos (pulses) via RabbitMQ, agregá-las temporariamente no Redis e posteriormente persistir os dados no PostgreSQL por meio de um cron worker.

---

## 📚 Sumário

- [🧠 Visão Geral](#-visão-geral)
- [⚙️ Tecnologias Utilizadas](#-tecnologias-utilizadas)
- [📁 Estrutura do Projeto](#-estrutura-do-projeto)
- [🚀 Executando o Projeto](#-executando-o-projeto)
- [🧪 Executando os Testes](#-executando-os-testes)
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

> ⚠️ **Observação:** Os comandos `make` utilizados neste README são opcionais. Eles servem apenas para facilitar a execução das tarefas mais comuns.
> 
> Se preferir (ou se estiver em um sistema que não tenha `make` instalado), você pode executar os comandos diretamente via `go run`, `go test`, etc.
>
> 💡 Dica: usuários Windows podem instalar `make` através do WSL, Chocolatey (`choco install make`) ou Git Bash.
>

1. **Clone o repositório**

```bash
git clone https://github.com/seuusuario/pulse-ingestor.git
cd pulse-ingestor
```

2. **Configure as variavies de ambiente**

   - *Crie o arquivo `.env`*

    - *Pegue os exemplos do `.env.example`*
        
```bash
cp .env.example .env
```
        
   - *Preencha os campos com os valores desejados*

3. **Suba os serviços com Docker Compose**

```
    docker-compose up -d
```

Esse comando irá iniciar:

- *`Redis`*
- *`RabbitMQ`*
- *`PostgreSQL`*
- *`Aplicação Go (main.go)`*

4. **Configurar e rodar o Script para popular a fila**

- *Você pode configurar a taxa de envio, o número de workers e os dados diretamente no script para testar concorrência e resiliência*.
- Altere as seguintes variaveis de ambiente para:

```bash
    SIMULATOR_TOTAL_MESSAGES=1000     #define o número de mensagens simuladas que vamos enviar para a fila 
    SIMULATOR_WORKERS_NUMBER=10       #/define o número de workers que vão realizar o processo de envio de mensagens
    SIMULATOR_BUFFER_SIZE=100         #define o tamanho do buffer para o envio das mensagens
```

- Rode o seguinte comando para começar a popular a fila com os dados fake:

  `Make run-script`

  ou se preferir:

  `go run main.go -script`

5. **Configurar e rodar o Worker para execução do serviço**

- Após a fila estar populada, deve rodar o seguinte comando para execução do serviço:

  `Make run-worker`

  ou se preferir:

  `go run main.go -worker`

## 🧪 Executando os Testes

- Você pode rodar todos os testes com:

  `Make test`

  ou se preferir:

  `go test ./... -v`

- Para acompanhar a cobertura de testes rode o comando:

  `Make coverage`
  
- **Os testes cobrem os fluxos de sucesso e falha dos workers,usecases ,repositorys, simulações de falhas no Redis e mensagens inválidas no RabbitMQ.**

## 🔁 Fluxo do Sistema

- **Ingestão**

    - Mensagens são consumidas do RabbitMQ por workers concorrentes.

    - Cada mensagem é desserializada e agregada em Redis usando chave composta.

- **Persistência**

    - Um cron job executa periodicamente (1h) a leitura das chaves agregadas no Redis.

    - Os dados são salvos em lote no PostgreSQL.

    - Após a persistência, as chaves correspondentes são removidas do Redis.

```bash
flowchart TD
    A = [Recebe mensagem via RabbitMQ] --> B [Parse JSON e validação]
    B --> C [Agrega dados no Redis]
    C --> D [Worker salva no Redis]
    E [Cron job a cada 1h] --> F [Lê agregados do Redis]
    F --> G [Persiste no PostgreSQL]
    G --> H [Limpa Redis após sucesso]
```

## 🧱 Arquitetura

- A aplicação segue princípios da `Clean Architecture`, separando responsabilidades em camadas:

  - **Domain:** interfaces e contratos
  - **Application:** lógica de negócio pura
  - **Infrastructure:** detalhes técnicos (drivers, clients)
  - **CMD:** ponto de entrada da aplicação
  - **Internal:** implementações internas
  - **Utils:** funções auxiliares reutilizáveis

- Benefícios:

  - Fácil testabilidade
  - Alta manutenibilidade
  - Independência de frameworks e ferramentas

## 📌 Considerações Finais

- ✅ Aplicação resiliente com `retries` e `backoff` no Redis

- ✅ Logs estruturados com `logrus`

- ✅ `Testes unitários` com cobertura dos principais fluxos

- ✅ Pronta para rodar com `docker-compose up`

- Desenvolvido com 💙 para o desafio técnico do Magalu Cloud.

## 📄 Licença

- Este projeto está licenciado sob os termos da [MIT License](LICENSE).

---
