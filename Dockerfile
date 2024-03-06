FROM smalls0098/go-alpine:3.19.1

LABEL maintainer="smalls0098@gmail.com"

ENV PROT=13822
ENV KEY=smalls0098

ENV WORKDIR /app
ENV TZ=Asia/Shanghai

ADD ./bin/app $WORKDIR/app
RUN chmod +x $WORKDIR/app

WORKDIR $WORKDIR

EXPOSE 13822
CMD ["./app"]