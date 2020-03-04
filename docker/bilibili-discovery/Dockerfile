from alpine

WORKDIR /app/bilibili

# download install sh
RUN wget https://raw.githubusercontent.com/bilibili/discovery/master/install.sh
# example config file
RUN wget https://raw.githubusercontent.com/bilibili/discovery/master/cmd/discovery/discovery-example.toml -O discovery.toml
# download binary here
RUN sh ./install.sh -b .

# export 7171 port
EXPOSE 7171

CMD ["./discovery", "-conf", "discovery.toml", "-alsologtostderr"]