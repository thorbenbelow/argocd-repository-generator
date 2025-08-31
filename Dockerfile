FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o argocd-repository-generator .

FROM gcr.io/distroless/static

WORKDIR /
COPY --from=builder /app/argocd-repository-generator .

EXPOSE 8080

ENTRYPOINT ["/argocd-repository-generator"]
