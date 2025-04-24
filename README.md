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
- [🔧 Executando o Projeto](#-executando-o-projeto)
- [🧪 Executando os Testes](#-executando-os-testes)
- [🔁 Fluxo do Sistema](#-fluxo-do-sistema)
- [🧱 Arquitetura](#-arquitetura)
- [🚀 Deploy da Aplicação](#-arquitetura-deploy)
- [🧭 API](#-api)
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
│   ├── api/               # Configurações da api (Popula a fila com dados fakes e simula alto volume de dados e consulta os dados no banco)
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

4. - **Bater na rota para pegar o token da api**

- body de exemplo:

```bash
  {
   "username":"fake",   #solicitar o username
   "password":"fake"    #solicitar a senha
}
```

- Segue a rota para login:

  - Método **POST**:

```bash
    https://ingestor-magalu-production.up.railway.app/api/v1/login
```

5. **Bater na rota para popular a fila**

- *Você pode configurar a taxa de envio, o número de workers e os dados diretamente na api para testar concorrência e resiliência*.

- body de exemplo:

```bash
  {
   "total_messages":100000,
   "workers_number":10,
   "buffer_size":100
  }
```

- Rode o seguinte comando para começar a rodar a api:

  `Make run-api`

  ou se preferir:

  `go run main.go -api`

- Bater na seguinte rota:

- Método **POST**:

```bash
    https://ingestor-magalu-production.up.railway.app/api/v1/pulses/populate
```

> ⚠️  **- Mais abaixo, na seção API, explico detalhadamente os demais endpoints disponíveis, além das configurações e funcionalidades da aplicação**.

6. **Configurar para rodar o Worker  - execução do serviço**

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

## 🚀 Deploy da Aplicação

- A aplicação foi implantada em produção utilizando a plataforma `Railway`, garantindo fácil escalabilidade, monitoramento e integração com os serviços externos (`RabbitMQ`, `Redis` e `PostgreSQL`).

- 🌐 Acesso à Aplicação

- 🔐 Link privado: apenas para avaliadores. Caso precise de acesso, entre em contato.

### ⚙️ Infraestrutura no Railway

- A infraestrutura foi provisionada com os seguintes serviços gerenciados pela Railway:

  - 🐇 **RabbitMQ** - responsável pela fila de mensagens.
  - 🧠 **Redis** - armazenamento temporário de dados agregados.
  - 🐘 **PostgreSQL** - persistência final dos dados.
  - 🔧 **Container Go Worker** - aplicação principal worker (Ingestor) construída com Docker.
  - 🧭 **Container Go API** - aplicação principal api (Ingestor) construída com Docker.

### 🗺️ Arquitetura da Solução

- A imagem abaixo resume a arquitetura geral do sistema:

<p align="center"> <img src="https://i.imgur.com/eVWtOkX.png" alt="Arquitetura da Solução" width="700"/> </p>

### 💡 Considerações sobre o deploy

  - O deploy foi configurado com `CD` via `GitHub Actions` e não pelo arquivo `ci.yml`, e o `CI` com etapas de teste, build e build automatico para produção configurados no projeto dentro do arquivo `ci.yml`.

  - O  `Redis` e o  `PostgreSQL ` utilizam volumes persistentes, garantindo integridade dos dados mesmo após reinicializações.

- O `RabbitMQ` está com painel de administração habilitado:
  
  - link: https://rabbitmq-web-ui-production-d644.up.railway.app

  - Para conseguir user e senha me solicitar.

## 🧭 API

- Como bônus, foi implementada uma pequena `API REST` apenas para fins de visualização dos dados agregados do `Pulse`. Essa `API` não faz parte da proposta original do desafio, mas pode ajudar o avaliador a consultar os dados persistidos e popular a fila diretamente via Postman ou navegador.

> ⚠️ Caso deseje acessar a rota pública (produção) do serviço hospedado, entre em contato e eu envio a URL.

### ⚙️ Como rodar a API

  - Adicione os seguinte comando no terminal: 

  `Make run-api`

  ou se preferir:

  `go run main.go -api`

### 🔍 Endpoints

  - 🔐 **LOGIN**

    - Primeiro passo é realizar o login na aplicação, a mesma retorna a seguinte estrutura:

      ```bash
      "Auth": {
          "Token": "fake",
          "ExpiresAt": 00
      }
      ```

    - Utilizar o token como `Bearer Token` em `Authorization` para realizar a autenticação

  - **✅ Listar Pulses**

    - **GET** /api/v1/pulses

    - Consulta os dados agregados de pulses com paginação.

    - Para acessar a rota é necessário incluir o token JWT no header Authorization no seguinte formato:

      - Authorization: Bearer <seu_token_aqui>

    - **Query Params:**

      - `page`: número da página (padrão: 1)

      - `limit`: número de itens por página (padrão: 10)

    - **Exemplo**:

      ```bash
      GET https://ingestor-magalu-production.up.railway.app//api/v1/pulses?page=2&limit=5
      ```

- 🔎 Buscar Pulse por ID

    - **GET**:

      - Consulta os dados de um pulse específico pelo seu ID.

      - Para acessar a rota é necessário incluir o token JWT no header Authorization no seguinte formato:

        - Authorization: Bearer <seu_token_aqui>

  - **Exemplo**:

      ```bash
    GET https://ingestor-magalu-production.up.railway.app/api/v1/pulses/42
      ```

- 📬 Popular fila com pulsos fakes

  - **POST**:

```bash
    https://ingestor-magalu-production.up.railway.app//api/v1/pulses/populate 
```

  - Popula a fila pulses com dados fakes para conseguir testar a aplicação.

  - Para acessar a rota é necessário incluir o token JWT no header Authorization no seguinte formato:

    - Authorization: Bearer <seu_token_aqui>

  - **Exemplo**:  

    - **BODY**:
    ```bash
        {
          "total_messages":100000,
          "workers_number":10,
          "buffer_size":100
        }
      ```

- Esses endpoints são acessíveis apenas para visualização dos dados no banco e para popular a fila, facilitando a validação do funcionamento da aplicação.



## 📌 Considerações Finais

- ✅ Aplicação resiliente com `retries` e `backoff` no Redis

- ✅ Logs estruturados com `logrus`

- ✅ `Testes unitários` com cobertura dos principais fluxos

- ✅ Pronta para rodar com `docker-compose up`

- Desenvolvido com 💙 para o desafio técnico do Magalu Cloud.

## 📄 Licença

- Este projeto está licenciado sob os termos da [MIT License](LICENSE).

---
