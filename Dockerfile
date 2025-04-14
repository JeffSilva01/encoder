FROM golang:1.24.2-alpine3.21

# Configuração das variáveis de ambiente em um único comando
ENV BENTO4_VERSION="1-6-0-641" \
  BENTO4_PATH="/opt/bento4" \
  BENTO4_BIN="/opt/bento4/bin" \
  PATH="$PATH:/opt/bento4/bin:/bin/bash"

WORKDIR /tmp/bento4

# Instalação de dependências em uma única camada com versões fixadas
RUN apk add --no-cache \
  ffmpeg=6.0.0-r3 \
  bash=5.2.21-r0 \
  make=4.4.1-r1 \
  cmake=3.27.8-r0 \
  python3=3.11.8-r0 \
  gcc=13.2.1_git20231014-r0 \
  g++=13.2.1_git20231014-r0 \
  wget=1.21.4-r0 \
  unzip=6.0-r14 && \
  # Download e descompactação do Bento4
  wget -q "http://zebulon.bok.net/Bento4/source/Bento4-SRC-${BENTO4_VERSION}.zip" && \
  unzip "Bento4-SRC-${BENTO4_VERSION}.zip" -d ${BENTO4_PATH} && \
  rm -f "Bento4-SRC-${BENTO4_VERSION}.zip"

# Compilação do Bento4
WORKDIR ${BENTO4_PATH}/cmakebuild
RUN cmake -DCMAKE_BUILD_TYPE=Release .. && \
  make

# Criação do diretório bin e cópia dos binários compilados
WORKDIR ${BENTO4_PATH}
RUN mkdir -p ${BENTO4_PATH}/bin && \
  find ${BENTO4_PATH}/cmakebuild -type f -executable -name "mp4*" -exec cp {} ${BENTO4_PATH}/bin/ \; && \
  find ${BENTO4_PATH}/cmakebuild -type f -executable -name "aac*" -exec cp {} ${BENTO4_PATH}/bin/ \; && \
  find ${BENTO4_PATH}/cmakebuild -type f -executable -name "*.py" -exec cp {} ${BENTO4_PATH}/bin/ \; && \
  # Cópia dos scripts Python
  cp -R ${BENTO4_PATH}/Source/Python/utils ${BENTO4_PATH}/ 2>/dev/null || true && \
  cp -a ${BENTO4_PATH}/Source/Python/wrappers/. ${BENTO4_PATH}/bin/ 2>/dev/null || true && \
  # Garantia de execução dos binários 
  chmod +x ${BENTO4_PATH}/bin/* && \
  # Limpeza para reduzir o tamanho da imagem
  apk del unzip wget && \
  rm -rf /var/cache/apk/* /tmp/*

# Configuração do diretório de trabalho da aplicação
WORKDIR /go/src

# Comando para manter o container rodando
CMD ["tail", "-f", "/dev/null"]
