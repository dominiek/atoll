FROM    ubuntu

# Install Atoll DEB
RUN     apt-get install -y wget apt-transport-https
RUN     wget -qO - https://deb.packager.io/key | sudo apt-key add -
RUN     echo "deb https://deb.packager.io/gh/dominiek/atoll trusty master" | sudo tee /etc/apt/sources.list.d/atoll.list
RUN     apt-get update
RUN     apt-get install -y atoll
RUN     /opt/atoll/bin/atoll -c /opt/atoll/atoll.yml setup --publish-api-key "docker-web-app-large"

# Install Atoll Elasticsearch Plugin DEB
RUN     wget -qO - https://deb.packager.io/key | sudo apt-key add -
RUN     echo "deb https://deb.packager.io/gh/dominiek/atoll-elasticsearch trusty master" | sudo tee /etc/apt/sources.list.d/atoll-elasticsearch.list
RUN     apt-get update
RUN     apt-get install -y atoll-elasticsearch
RUN     echo "plugins:" >> /opt/atoll/atoll.yml
RUN     echo "  - /opt/atoll-elasticsearch/bin/atoll-elasticsearch" >> /opt/atoll/atoll.yml

# Install Elasticsearch
RUN     apt-get install -y software-properties-common python-software-properties
RUN     add-apt-repository ppa:webupd8team/java -y
RUN     wget -qO - https://packages.elastic.co/GPG-KEY-elasticsearch | sudo apt-key add -
RUN     echo "deb http://packages.elastic.co/elasticsearch/1.6/debian stable main" | sudo tee -a /etc/apt/sources.list.d/elk.list
RUN     apt-get update
RUN     echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | sudo debconf-set-selections
RUN     apt-get install -y oracle-java8-installer elasticsearch curl

EXPOSE  9200
CMD ./opt/atoll/bin/atoll -c /opt/atoll/atoll.yml -d && /usr/share/elasticsearch/bin/elasticsearch
