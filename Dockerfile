# ---------- STAGE 1: build ----------
    FROM golang:1.22 AS build

    # Set working directory inside the container for building the app
    WORKDIR /src
    
    # Copy go.mod and go.sum first to leverage Docker cache for dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the entire source code into the container
    COPY . .
    
    # Build a static Go binary for Linux
    # - CGO_ENABLED=0: disables CGO for a fully static binary
    # - GOOS=linux, GOARCH=amd64: target Linux x86_64
    # - ldflags "-s -w": strip debug info → smaller binary size
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -ldflags="-s -w" -o /bin/ipbot ./cmd/bot
    
    # ---------- STAGE 2: final ----------
    # Use a minimal and secure runtime image without shell or root access
    FROM gcr.io/distroless/static:nonroot
    
    # Set the working directory for the app
    WORKDIR /app
    
    # Copy the compiled Go binary from the build stage
    COPY --from=build /bin/ipbot /app/ipbot
    
    # Copy database migrations into the image (if you use them)
    COPY ./migrations /app/migrations
    
    # Run the app as a non-root user for better security
    USER nonroot:nonroot
    
    # Start the application
    ENTRYPOINT ["/app/ipbot"]
    