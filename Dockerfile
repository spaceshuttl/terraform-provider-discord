FROM golang:alpine as builder
WORKDIR /go/src/github.com/spaceshuttl/terraform-provider-discord
RUN apk add --no-cache make git
ENV GO111MODULE=on
COPY . .
RUN make bin

FROM hashicorp/terraform:light
COPY --from=builder /go/src/github.com/spaceshuttl/terraform-provider-discord/terraform-provider-discord /usr/local/terraform-plugins/
WORKDIR /opt/workspace