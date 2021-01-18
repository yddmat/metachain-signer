FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y \
        wget \
        curl \
        git \
        vim \
        unzip \
        xz-utils \
        software-properties-common \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN wget -O - https://apt.kitware.com/keys/kitware-archive-latest.asc | apt-key add - \
    && apt-add-repository 'deb https://apt.kitware.com/ubuntu/ bionic main' \
    && apt-add-repository -y ppa:mhier/libboost-latest \
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys DE19EB17684BA42D

RUN apt-get update \
    && apt-get install -y \
        build-essential \
        libtool autoconf pkg-config \
        ninja-build \
        ruby-full \
        clang-10 \
        llvm-10 \
        libc++-dev libc++abi-dev \
        cmake \        
        libboost1.74-dev \
        ccache \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

ENV CC=/usr/bin/clang-10
ENV CXX=/usr/bin/clang++-10

COPY . /app
WORKDIR /app

RUN git clone --depth 1 --branch 2.5.3 https://github.com/trustwallet/wallet-core.git

WORKDIR /app/wallet-core
RUN tools/install-dependencies

RUN tools/generate-files \
    && cmake -H. -Bbuild -DCMAKE_BUILD_TYPE=Debug \
    && make -Cbuild -j12

WORKDIR /app
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update && apt-get install -y golang-go
RUN go build .

CMD ["./multichain-signer"]