FROM docker:dind

ARG ISTIO_VERSION=1.13.1

RUN apk add --update --no-cache curl coreutils bash git ca-certificates musl-dev go

RUN curl -L https://istio.io/downloadIstio | ISTIO_VERSION=$ISTIO_VERSION TARGET_ARCH=x86_64 sh - && \
    ln -s /istio-$ISTIO_VERSION/bin/istioctl /usr/local/bin/istioctl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \ 
    mv kubectl /usr/local/bin && \
    chmod +x /usr/local/bin/kubectl && \
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64 && \
    mv kind /usr/local/bin && \
    chmod +x /usr/local/bin/kind && \
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && ./get_helm.sh


# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
# RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

ENV PATH /usr/local/go/bin:$PATH

ENV GOPATH /go
ENV ISTIO_HOME /istio-$ISTIO_VERSION
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" && mkdir -p "$GOPATH/e2e"
WORKDIR $GOPATH/e2e

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY entrypoint.sh ./

RUN chmod +x entrypoint.sh && go mod tidy

ENTRYPOINT ["./entrypoint.sh"]