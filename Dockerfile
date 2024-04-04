# syntax=docker/dockerfile:1

FROM golang:1.22.1-alpine3.19

# AlpineLinuxのパッケージ管理コマンドapkをupdate, gitをインストール
RUN apk update && apk add git

ENV TZ /usr/share/zoneinfo/Asia/Tokyo

# 起動したいサーバーのディレクトリ
WORKDIR /hello-world

# 第二引数の./は、WORKDIRで指定したディレクトリから見て./なので、すなわち/hello-world
# なので、ホストマシンの./hello-world/*を、コンテナ内の/hello-worldにコピーしているということ
COPY ./hello-world/* ./

# go.modに記載されているパッケージをインストール
RUN go mod download

# ホットリローダーairをインストール
RUN go install github.com/cosmtrek/air@latest

EXPOSE 5050

# airを使って起動
CMD ["air", "-c", ".air.toml"]

# 検証用 すぐコンテナ終了しないように
# CMD ["tail", "-f", "/dev/null"]