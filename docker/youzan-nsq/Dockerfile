FROM amd64/alpine

WORKDIR /youzan/nsq

COPY nsq-0.3.7-HA.1.9.4.linux-amd64.go1.10.8.tar.gz .
COPY  nsqadmin.conf .
COPY  nsqd.conf . 
COPY  nsqlookupd.conf .

RUN tar -zxvf nsq-0.3.7-HA.1.9.4.linux-amd64.go1.10.8.tar.gz \ 
    && mv nsq-0.3.7-HA.1.9.4.linux-amd64.go1.10.8/bin/* /usr/bin && \
    rm nsq-0.3.7-HA.1.9.4.linux-amd64.go1.10.8.tar.gz && \
    rm -fr nsq-0.3.7-HA.1.9.4.linux-amd64.go1.10.8

# RUN wget https://raw.githubusercontent.com/youzan/nsq/master/contrib/nsqadmin.cfg.example -O nsqadmin.conf && \
    # wget https://raw.githubusercontent.com/youzan/nsq/master/contrib/nsqd.cfg.example -O nsqd.conf && \
    # wget https://raw.githubusercontent.com/youzan/nsq/master/contrib/nsqlookupd.cfg.example -O nsqlookupd.conf

EXPOSE 4150 4151 4160 4161 4170 4171