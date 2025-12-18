FROM golang:1.25.5-alpine3.23 AS build

ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -ldflags "\
    -X 'main.Version=${VERSION}' \
    -X 'main.GitCommit=${GIT_COMMIT}' \
    -X 'main.BuildDate=${BUILD_DATE}'" \
  -o /build/axosec-vault /build/cmd/api

FROM alpine:3.23 AS runner
WORKDIR /
COPY --from=build /build/axosec-vault /axosec-vault
CMD ["/axosec-vault"]
