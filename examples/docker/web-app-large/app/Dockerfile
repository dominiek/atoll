FROM    centos:centos6
RUN     rpm -Uvh http://download.fedoraproject.org/pub/epel/6/i386/epel-release-6-8.noarch.rpm

# Install Atoll RPM
RUN     rpm --import https://rpm.packager.io/key
RUN     echo -e "[atoll]\nname=Repository for dominiek/atoll application.\nbaseurl=https://rpm.packager.io/gh/dominiek/atoll/centos6/master\nenabled=1" | tee /etc/yum.repos.d/atoll.repo
RUN     yum install -y atoll
RUN     /opt/atoll/bin/atoll -c /opt/atoll/atoll.yml setup --publish-api-key "docker-web-app-large"

RUN     yum install -y npm
COPY    . /deploy
RUN     cd /deploy; npm install
EXPOSE  80
CMD ./opt/atoll/bin/atoll -c /opt/atoll/atoll.yml -d && node /deploy/app.js
