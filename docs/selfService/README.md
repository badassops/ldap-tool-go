## Self Service LDAP Password

### Background
This is an open source PHP web application that will allow an LDAP user to update their own password or SSH public key

### Installation
- install the required packages
```
apt install -y php-fpm php-gd smarty3
```
- download the sources
  The source can be found here [ Self Service Password](https://ltb-project.org/download.html)
- extra the source under /usr/share and then rename to /usr/share/self-service-password
- edit the configuration file  i/usr/share/self-service-password/conf/config.inc.php 
  you can find the one I use [here](https://github.com/badassops/ldap-tool-go/blob/main/docs/selfService/example-config.inc.php) (minus secrets 😈)
- set permission correct with 
```
cd  /usr/share/self-service-password
find . -type d | xargs chmod 0755
find . -type f | xargs chmod 0644
chown root:www-data conf/config.inc.php
chmod 0440 conf/config.inc.php
```
- copy the images you would like under the images directory
```
cp image image /usr/share/self-service-password/htdocs/images
chmod -444 /usr/share/self-service-password/htdocs/images/*
chown root:root /usr/share/self-service-password/htdocs/images/*
```
- setup php-fpm and nginx to serve the self service application (beyond this scope)
  you can find and example of my nginx configuration [here](https://github.com/badassops/ldap-tool-go/blob/main/docs/selfService/example-nginx-self-service.conf)


## The End
Congraculation you should be all set now : 🦄👏
 
