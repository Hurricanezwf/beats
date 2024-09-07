FROM --platform=linux/amd64 elastic/filebeat:8.14.2

ADD ./filebeat      /usr/share/filebeat/filebeat-my
ADD ./module/cls    /usr/share/filebeat/module/

CMD ["/usr/share/filebeat/filebeat"]
