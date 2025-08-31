FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build the Go application
# CGO_ENABLED=0: disables CGO, making the binary statically linked. This is CRUCIAL for distroless images.
# GOOS=linux: ensures the binary is compiled for Linux.
# -a -installsuffix nocgo: ensures all packages are rebuilt and avoids issues with CGO-linked libraries.
# -o argocd-repositorty-generator: specifies the output file name as 'main'.
# -ldflags "-s -w": reduces the binary size by omitting debugging information (symbol table and DWARF sections).
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o argocd-repositorty-generator .

FROM gcr.io/distroless/static
LABEL org.opencontainers.image.source = "https://github.com/thorbenbelow/argocd-repository-generator"
WORKDIR /
COPY --from=builder /app/argocd-repositorty-generator .

EXPOSE 8080

ENTRYPOINT ["/argocd-repositorty-generator"]
