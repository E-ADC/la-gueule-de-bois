# --- Étape de build ---
FROM golang:1.25-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

# --- Image finale minimale ---
FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=build /out/api /app/api
COPY migrations /app/migrations

# Volume attendu pour les photos uploadées (spec : écrites tel quel sur
# disque, servies ensuite par nginx en prod).
VOLUME ["/app/uploads"]

EXPOSE 8080
ENTRYPOINT ["/app/api"]
