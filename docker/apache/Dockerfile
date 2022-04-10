FROM httpd:2.4

RUN /bin/bash -c 'echo "Include /usr/local/apache2/conf/include/*.conf"' | \ 
    tee -a "/usr/local/apache2/conf/httpd.conf"
