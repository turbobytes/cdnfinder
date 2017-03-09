FROM debian

RUN apt-get update && apt-get install -y fontconfig && rm -rf /var/lib/apt/lists/*

ADD bin/* /bin/