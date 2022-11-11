FROM golang:1.19 AS build 

WORKDIR /srv/build 

COPY go.sum go.mod Makefile ./

RUN make deps 

COPY . . 

RUN make build 

FROM ubuntu:20.04 AS final 

COPY --from=build /srv/build/linkshortener /bin/
COPY --from=build /srv/build/.cobra.yaml /

ENTRYPOINT ["/bin/linkshortener", "serve"]
