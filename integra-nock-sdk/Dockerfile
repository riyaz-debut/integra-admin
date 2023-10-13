FROM golang:1.18.2-alpine

RUN mkdir -p /integra-admin-sdk

WORKDIR /integra-admin-sdk

RUN apk update && \
    apk add curl git pkgconfig && \
    apk add openssl-dev && \
    apk add musl-dev && \
    apk add gcc 
    # CGO_ENABLED=0 && \
    
# Cache and install dependencies integra-nock-sdk/go.mod
COPY . .
# COPY integra-nock-sdk/go.mod integra-nock-sdk/go.sum ./
RUN go mod download




################     
# FROM node:16.16.0

# RUN mkdir -p /hyperledger-sdk

# WORKDIR /hyperldger-sdk

# COPY package*.json ./

# RUN npm install
# COPY . .
# CMD [ "npm", "start" ]

