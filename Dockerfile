FROM golang:1.16.4-buster

RUN apt-get update
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
RUN apt-get install -y nodejs python-pip sudo
RUN pip install pipenv
RUN npm install -g yarn ts-node

# Install golangci-lint
RUN version=1.40.0 \
    && curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /go/bin v$version \
    && golangci-lint version

# Install pulumictl
RUN version=0.0.25 \
    && curl -fsSLO https://github.com/pulumi/pulumictl/releases/download/v$version/pulumictl-v$version-linux-amd64.tar.gz \
    && tar -xzf pulumictl-v$version-linux-amd64.tar.gz --directory /usr/local/bin --no-same-owner pulumictl \
    && rm -f pulumictl-v$version-linux-amd64.tar.gz \
    && pulumictl version

# Add non-root user
ARG USER_NAME=user
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USER_NAME \
    && useradd --uid $USER_UID --gid $USER_GID --shell /bin/bash -m $USER_NAME \
    && echo "$USER_NAME ALL=(ALL:ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER_NAME \
    && chmod 0440 /etc/sudoers.d/$USER_NAME

WORKDIR /var/pulumi
COPY . .

RUN mkdir -p /go/bin \
    && chown -R $USER_NAME: /go \
    && mkdir -p /opt/pulumi/bin \
    && chown -R $USER_NAME: /opt/pulumi \
    && chown -R $USER_NAME: /var/pulumi

ENV PULUMI_SKIP_UPDATE_CHECK=true
ENV USER=$USER_NAME

USER $USER_NAME

RUN echo "export PATH=/opt/pulumi:/opt/pulumi/bin:$GOPATH/bin:/usr/local/go/bin:$PATH" >> ~/.bashrc \
    && echo "unset XDG_CACHE_HOME XDG_CONFIG_HOME" >> ~/.bashrc \
    && echo "alias l='ls -aF'" >> ~/.bash_aliases \
    && echo "alias ll='ls -ahlF'" >> ~/.bash_aliases \
    && echo "alias ls='ls --color=auto --group-directories-first'" >> ~/.bash_aliases

RUN make install
RUN cd sdk/nodejs && yarn && make build && make install
