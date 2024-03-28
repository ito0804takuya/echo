FROM golang:1.20.14-alpine3.19

# AlpineLinuxのパッケージ管理コマンドapkをupdate, gitをインストール
RUN apk update && apk add git

ENV TZ /usr/share/zoneinfo/Asia/Tokyo

# 起動したいサーバーのディレクトリ
WORKDIR /hello-world

COPY . .

# go.modに記載されているパッケージをインストール
RUN go mod download

# ホットリローダーairをインストール
RUN go install github.com/cosmtrek/air@latest

EXPOSE 5050

# airを使って起動
CMD ["air", "-c", ".air.toml"]
