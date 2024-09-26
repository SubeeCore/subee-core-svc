FROM golang:1 AS build

ARG GO_BUILDER_GITHUB_TOKEN=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
RUN git config --global "url.https://$GO_BUILDER_GITHUB_TOKEN@github.com/".insteadOf https://github.com/
ENV GOPRIVATE=github.com/subeecore

WORKDIR /app

COPY . .

RUN go install github.com/cespare/reflex@latest
RUN go mod download

RUN go build -o main ./cmd

EXPOSE 3000

CMD ["./main"]