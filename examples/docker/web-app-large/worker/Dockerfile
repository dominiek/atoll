FROM    ubuntu

# Install Atoll DEB
RUN     apt-get install -y wget apt-transport-https
RUN     wget -qO - https://deb.packager.io/key | sudo apt-key add -
RUN     echo "deb https://deb.packager.io/gh/dominiek/atoll trusty master" | sudo tee /etc/apt/sources.list.d/atoll.list
RUN     apt-get update
RUN     apt-get install -y atoll
RUN     /opt/atoll/bin/atoll -c /opt/atoll/atoll.yml setup --publish-api-key "docker-web-app-large"

# Install Atoll RabbitMQ Plugin DEB
RUN     wget -qO - https://deb.packager.io/key | sudo apt-key add -
RUN     echo "deb https://deb.packager.io/gh/dominiek/atoll-rabbitmq trusty master" | sudo tee /etc/apt/sources.list.d/atoll-rabbitmq.list
RUN     apt-get update
RUN     apt-get install -y atoll-rabbitmq
RUN     echo "plugins:" >> /opt/atoll/atoll.yml
RUN     echo "  - /opt/atoll-rabbitmq/bin/atoll-rabbitmq" >> /opt/atoll/atoll.yml

# Install RabbitMQ
RUN     apt-get install -y rabbitmq-server
RUN     rabbitmq-plugins enable rabbitmq_management

RUN     apt-get install -y npm
COPY    . /deploy
RUN     cd /deploy; npm install
RUN     ln -s /usr/bin/nodejs /usr/bin/node

EXPOSE  5672
CMD ./opt/atoll/bin/atoll -c /opt/atoll/atoll.yml -d && /usr/sbin/rabbitmq-server & node /deploy/app.js
