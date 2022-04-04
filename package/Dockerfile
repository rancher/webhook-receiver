FROM registry.suse.com/bci/bci-micro:15.3.16.3

ARG user=webhookreceiver

RUN echo "$user:x:1000:1000::/home/$user:/bin/bash" >> /etc/passwd && \
    echo "$user:x:1000:" >> /etc/group && \
    mkdir /home/$user && \
    chown -R $user:$user /home/$user

COPY bin/webhook-receiver /webhook-receiver

USER $user

CMD ["/webhook-receiver"]
