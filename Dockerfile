FROM node AS Frontend
ADD ./frontend /go/src/github.com/zottelchin/sam/frontend
WORKDIR /go/src/github.com/zottelchin/sam/frontend
RUN npm install -g parcel-bundler
RUN npm install && \
    parcel build index.html

FROM golang AS Backend
ADD *.go /go/src/github.com/zottelchin/sam/
COPY --from=Frontend /go/src/github.com/zottelchin/sam/frontend/dist /go/src/github.com/zottelchin/sam/frontend/dist
WORKDIR /go/src/github.com/zottelchin/sam
RUN go get github.com/go-bindata/go-bindata/... && \
    /go/bin/go-bindata ./frontend/dist/... && \
    go get
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static" -s' -installsuffix cgo -o SAM -v .

FROM scratch
COPY --from=Backend /go/src/github.com/zottelchin/sam/SAM /SAM
COPY --from=Frontend /go/src/github.com/zottelchin/sam/frontend/dist /dist
ADD ./sql /sql
ENV GIN_MODE=release
WORKDIR /
EXPOSE 2222
VOLUME [ "/data" ]

ENTRYPOINT [ "/SAM" ]