packer-builder-sakuracloud
===

A builder plugin of packer for SakuraCloud

さくらのクラウド用Packerプラグイン

## 概要

`packer-builder-sakuracloud`は[さくらのクラウド](http://cloud.sakura.ad.jp)での
アーカイブ(構築済みOSのテンプレート)を[Packer](https://packer.io)で作成するためのPackerプラグインです。

## 使い方(セットアップ)

`packer-builder-sakuracloud`を実行するには、さくらのクラウドAPIキーが必要です。

あらかじめコントロールパネルからAPIキーを発行しておいてください。

セットアップは以下のような方法があります。お好きな方をご利用ください。

 - ローカルマシンにインストール
 - Dockerを利用する

### ローカルマシンへのインストール

[リリースページ](https://github.com/sacloud/packer-builder-sakuracloud/releases/latest)から各プラットフォーム用のバイナリをダウンロードし、
以下の何れかのディレクトリへ展開、実行権を付与してください。

- 1) `packer`コマンドが格納されているディレクトリ
- 2) Unix系OSの場合 `~/.packer.d/plugins` , Windows系OSの場合`%APPDATA%/packer.d/plugins`

### Dockerを利用する場合

Dockerを利用する場合、`packer-builder-sakuracloud`の事前のインストール作業は不要です。
DockerHubでイメージを公開していますので、以下のように実行できます。

    docker run -it --rm sacloud/packer:latest [packerサブコマンド] [packerオプション]

APIキーを環境変数で指定する場合、`-e`オプションなどを適切に指定して実行してください。

    # APIキーを環境変数で指定して実行する例
    $ docker run -it --rm \
             -e SAKURACLOUD_ACCESS_TOKEN \
             -e SAKURACLOUD_ACCESS_TOKEN_SECRET \
             sacloud/packer:latest [packerサブコマンド] [packerオプション]

## 使い方(アーカイブ作成)

#### APIキーの設定

APIキーを環境変数に設定しておきます。
(APIキーは以下で作成するjsonファイルに記載することも可能です)

    $ export SAKURACLOUD_ACCESS_TOKEN=[APIトークン]
    $ export SAKURACLOUD_ACCESS_TOKEN_SECRET=[APIシークレット]

#### アーカイブ定義ファイル(json)の作成

Packerでのビルド用に以下のようなjsonファイルを作成します。
以下の例は、石狩第2ゾーン(`is1b`)に、CentOSパブリックアーカイブをベースとしたアーカイブを作成します。
プロビジョニングとして、`shell`にてdockerのインストールを行なっています。

#### packer用jsonファイルの例
```sh
$ cat <<EOF > example.json
{
    "builders": [{
        "type": "sakuracloud",
        "zone": "is1b",
        "os_type": "centos",
        "password": "TestUserPassword01",
    }],
    "provisioners":[
    {
        "type": "shell",
        "inline": [
            "yum update -y",
            "curl -fsSL https://get.docker.com/ | sh",
            "systemctl enable docker.service"
        ]
    }]
}
EOF
```

作成したら以下のように`packer build`を実行すると、さくらのクラウド上にアーカイブが作成されます。

    $ packer build example.json

## オプション一覧

jsonファイルで指定できるオプションの一覧は以下の通りです。

### 必須項目

- `access_token`(string): APIトークン。この値は環境変数`SAKURACLOUD_ACCESS_TOKEN`で指定することも可能です。

- `access_token_secret`(string): APIシークレット。この値は環境変数`SAKURACLOUD_ACCESS_TOKEN_SECRET`で指定することも可能です。

- `zone`(string): 対象ゾーン。`tk1a`(東京)又は`is1b`(石狩第2)を指定

- `os_type`(string): ベースとするアーカイブの種別。以下の値が指定可能です。
    - `centos`: CentOSパブリックアーカイブ(さくらのクラウドで提供されている最新安定版)
    - `ubuntu`: Ubuntuパブリックアーカイブ
    - `debian`: Debianパブリックアーカイブ
    - `coreos`: CoreOSパブリックアーカイブ
    - `kusanagi`: Kusanagi(CentOSベース)パブリックアーカイブ
    - `windows`: Windows系パブリックアーカイブ
    - `custom`: その他アーカイブ、又はディスク

`os_type`が`windows`の場合、`source_archive`の指定が必須です。

`os_type`が`custom`の場合、`source_archive` 又は `source_disk`の何れかの指定が必須です。

### オプション項目

- `user_name`(string): SSH/WinRM接続時のユーザー名、デフォルト値:`root`

- `password`(string): SSH/WinRM接続時のユーザー名

- `disk_size`(int): 作成するディスクのサイズ(GB単位)、デフォルト値:`20`

- `disk_connection`(string): ディスク接続方法、以下の値が指定可能です。デフォルト値:`virtio`
    - `ide`: IDE接続
    - `virtio`: 準仮想モード(virtio)

- `disk_plan`(string): ディスクプラン、以下の値が指定可能です。デフォルト値:`ssd`
    - `ssd`: SSDプラン
    - `hdd`: 通常プラン(ハードディスク)

- `core`(int): CPUコア数。デフォルト値:`1`

- `memory_size`(int): メモリサイズ(GB単位)、デフォルト値:`1`

- `disable_virtio_net`(bool): `true`の場合、NICでの仮想化ドライバ利用を無効化します。デフォルト値:`false`

- `source_archive`(int64): 元となるアーカイブのID、`os_type`が`windows`、又は`custom`の場合に指定可能です。

- `source_disk`(int64): 元となるディスクのID、`os_type`が`custom`の場合のみ指定可能です。

- `archive_name`(string): 作成されるアーカイブの名前

- `archive_tags`([]string): 作成されるアーカイブに付与するタグ、デフォルト値:`["@size-extendable"]`

- `archive_description`(string): 作成されるアーカイブに付与する説明

## License

  `packer-builder-sakuracloud` Copyright (C) 2016 Kazumichi Yamamoto.

  This project is published under [MPL-2.0](LICENSE).
  
## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))

