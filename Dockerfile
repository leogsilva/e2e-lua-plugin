FROM docker:dind

ARG ISTIO_VERSION

RUN apk --update add curl coreutils bash git

RUN curl -L https://istio.io/downloadIstio | ISTIO_VERSION=$ISTIO_VERSION TARGET_ARCH=x86_64 sh - && \
    ln -s istio-$ISTIO_VERSION/bin/istioctl /usr/local/bin/
    
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \ 
    mv kubectl /usr/local/bin && \
    chmod +x /usr/local/bin/kubectl