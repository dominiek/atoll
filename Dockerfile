FROM    ubuntu
RUN     apt-get update
RUN     apt-get install -y golang-go make
COPY    . /src
CMD     cd /src; $BUILD_COMMAND
