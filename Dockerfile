FROM cgr.dev/chainguard/go:latest AS build

ARG APP_NAME=ingestor

WORKDIR /go/src/${APP_NAME}
COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/${APP_NAME}

FROM cgr.dev/chainguard/static:latest

ARG APP_NAME=ingestor
ENV ENV_APP_NAME=$APP_NAME

COPY --from=build /go/bin/${APP_NAME} /

ENTRYPOINT ["/ingestor"]
CMD [""]
