FROM golang:alpine as builder
WORKDIR /go/src/github.com/spaceshuttl/terraform-provider-discord
RUN apk add --no-cache make git
ENV GO111MODULE=on
COPY . .
RUN make bin

FROM hashicorp/terraform:light
RUN mkdir -p /root/.terraform.d/plugins
COPY --from=builder /go/src/github.com/spaceshuttl/terraform-provider-discord/terraform-provider-discord /root/.terraform.d/plugins/
WORKDIR /opt/workspace