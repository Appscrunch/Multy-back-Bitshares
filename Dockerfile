FROM golang:1.9.4

ENV REPO "github.com/Appscrunch/Multy-Back-Bitshares"

COPY ./ "$GOPATH/src/$REPO"

RUN cd $GOPATH/src/$REPO && \
    make all-with-deps

RUN ls $GOPATH/src/

WORKDIR /go/src/github.com/Appscrunch/Multy-Back-Bitshares/cmd

EXPOSE 8080

ENTRYPOINT $GOPATH/src/$REPO/multy-bitshares
