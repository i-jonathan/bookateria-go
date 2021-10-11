#FROM golang:1.16.8-alpine AS builder
#
#WORKDIR /bookateria
#
#COPY . .
#
#RUN apk --no-cache -U add libc-dev build-base
#RUN go mod download && go mod tidy

FROM bookateria-base AS builder

WORKDIR /bookateria
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -o main .

FROM scratch
COPY --from=builder /bookateria/main ./main
COPY .env .
COPY docs docs
CMD [ "./main" ]