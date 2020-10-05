packer-builder-sakuracloud
===

[![Build Status](https://travis-ci.org/sacloud/packer-builder-sakuracloud.svg?branch=master)](https://travis-ci.org/sacloud/packer-builder-sakuracloud)
[![Slack](https://slack.usacloud.jp/badge.svg)](https://slack.usacloud.jp/)  

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
 - (macOS + Homebrew)`homebrew`でインストール

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

    # APIキーを環境変数で指定、カレントディレクトリのexample.jsonをテンプレートとして指定してビルドする例
    $ docker run -it --rm \
             -e SAKURACLOUD_ACCESS_TOKEN \
             -e SAKURACLOUD_ACCESS_TOKEN_SECRET \
             -v $PWD:/work \
             sacloud/packer:latest build example.json

### `homebrew`を利用する場合(macOSでHomebrewをご利用の場合)

以下のコマンドでインストール可能です。

    brew tap sacloud/homebrew-packer-builder-sakuracloud; brew install packer-builder-sakuracloud

インストール後は画面に表示される指示にしたがってコマンドを実行し、`~/.packer.d/plugins`配下にシンボリックリンクを作成します。


## 使い方(アーカイブ作成)

#### APIキーの設定

APIキーを環境変数に設定しておきます。
(APIキーは以下で作成するjsonファイルに記載することも可能です)

    $ export SAKURACLOUD_ACCESS_TOKEN=[APIトークン]
    $ export SAKURACLOUD_ACCESS_TOKEN_SECRET=[APIシークレット]

#### テンプレートファイル(JSON)の作成

Packerでのビルド用に以下のようなJSONファイルを作成します。
以下の例は、石狩第2ゾーン(`is1b`)に、CentOSパブリックアーカイブをベースとしたアーカイブを作成します。
プロビジョニングとして、`shell`にてdockerのインストールを行なっています。

#### packer用jsonファイルの例
```sh
$ cat <<EOF > example.json
{
    "builders": [{
        "type": "sakuracloud",
        "zone": "is1b",
        "os_type": "centos7",
        "password": "TestUserPassword01"
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

さくらのクラウド上のパブリックアーカイブだけでなく、ISOイメージからの構築も可能です。
詳細は[テンプレートのサンプル](#テンプレートサンプル)を参照してください。

#### テンプレートファイル(HCL)の作成

Packer 1.5以降で利用できるHCLテンプレートに対応しています。 
以下の様なファイルを作成し拡張子を`.pkr.hcl`とすることで`packer build`が実行できます。

```hcl
source "sakuracloud" "example" {
  zone = "is1b"
  zones = ["is1a", "is1b", "tk1a", "tk1v"]

  os_type   = "centos7"
  password  = "TestUserPassword01"
  disk_size = 20
  disk_plan = "ssd"

  core        = 2
  memory_size = 4

  archive_name        = "packer-example-centos"
  archive_description = "description of archive"
}

build {
  sources = [
    "source.sakuracloud.example"
  ]
  provisioner "shell" {
    inline = [
      "yum update -y",
      "curl -fsSL https://get.docker.com/ | sh",
      "systemctl enable docker.service",
    ]
  }
}
```


## オプション一覧

jsonファイルで指定できるオプションの一覧は以下の通りです。

### 必須項目

- `access_token`(string): APIトークン。この値は環境変数`SAKURACLOUD_ACCESS_TOKEN`で指定することも可能です。

- `access_token_secret`(string): APIシークレット。この値は環境変数`SAKURACLOUD_ACCESS_TOKEN_SECRET`で指定することも可能です。

- `zone`(string): 対象ゾーン。`tk1a`(東京)又は`is1b`(石狩第2)を指定

- `os_type`(string): ベースとするアーカイブの種別。以下の値が指定可能です。

| 値                              | 説明                                 |
|-- -------------------------    | ------------------               --|
| `centos`                       | CentOS(最新安定板)                      |
| `centos8`                      | CentOS 8                           |
| `centos7`                      | CentOS 7                           |
| `centos6`                      | CentOS 6                           |
| `ubuntu`                       | Ubuntu(最新安定板)                      |
| `ubuntu2004`                   | Ubuntu 20.04                       |
| `ubuntu1804`                   | Ubuntu 18.04                       |
| `ubuntu1604`                   | Ubuntu 16.04                       |
| `debian`                       | Debian(最新安定板)                      |
| `debian10`                     | Debian10                           |
| `debian9`                      | Debian9                            |
| `coreos`                       | CoreOS                             |
| `rancheros`                    | RancherOS                          |
| `k3os`                         | k3OS                               |
| `kusanagi`                     | Kusanagi(CentOS7)                  |
| `freebsd`                      | FreeBSD                            |
| `windows2016`                  | Windows 2016                       |
| `windows2016-rds`              | Windows 2016(RDS)                  |
| `windows2016-rds-office`       | Windows 2016(RDS + Office)         |
| `windows2016-sql-web`          | Windows 2016 SQLServer(Web)        |
| `windows2016-sql-standard`     | Windows 2016 SQLServer(Standard)   |
| `windows2016-sql-standard-all` | Windows 2016 SQLServer(RDS+Office) |
| `windows2019`                  | Windows 2019 Datacenter Edition    |
| `custom`                       | 任意のアーカイブID/ディスクIDを指定する場合           |
| `iso`                          | さくらのクラウド上のISOイメージ、<br />またはURLを指定してISOイメージをダウンロードする場合 |


`os_type`が`custom`の場合、`source_archive` 又は `source_disk`の何れかの指定が必須です。

`os_type`が`iso`の場合、`iso_id`、または以下のISOイメージ関連の値の指定が必須です。
  - `iso_url`または`iso_urls` : ISOイメージのURL
  - `iso_checksum`/`iso_checksum_url`: ISOイメージのチェックサム、またはチェックサムを記載したファイルのURL
  - `iso_checksum_type`: チェックサムの形式(`none`, `md5`, `sha1`, `sha256`, or `sha512`のいずれか)

`iso`の場合の詳細は[ISOイメージ関連項目の指定について](#isoイメージ関連項目の指定について)を参照ください。

### オプション項目

- `zones`: 作成したアーカイブを転送する宛先ゾーン名のリスト。以下のいずれかの値を指定(複数指定可)
    - `is1a`: 石狩第1ゾーン
    - `is1b`: 石狩第2ゾーン
    - `tk1a`: 東京第1ゾーン
    - `tk1v`: サンドボックス

- `user_name`(string): SSH/WinRM接続時のユーザー名

`os_type`にWindows系を指定した場合のデフォルト値:`Administrator`

`os_type`がWindows系以外の場合のデフォルト値:`root`

- `password`(string): SSH/WinRM接続時のパスワード

- `us_keyboard`(bool): ISOイメージからのインストール時、コマンド送信にUSキーボードレイアウトを利用するか、デフォルト値:`false`

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

- `iso_size`(int): アップロードするISOファイルのサイズ、`5`または`10`が指定可能(GB単位)、デフォルト値: `5`

- `iso_name`(string): アップロードするISOファイルの名前、デフォルト値: `iso_checksum`が空でない場合は`iso_checksum`、以外はタイムスタンプから自動生成

- `archive_name`(string): 作成されるアーカイブの名前

- `archive_tags`([]string): 作成されるアーカイブに付与するタグ、デフォルト値:なし(空の場合、かつ祖先にパブリックアーカイブが存在する場合は祖先となったパブリックアーカイブからタグを引き継ぐ)

- `archive_description`(string): 作成されるアーカイブに付与する説明

- `boot_wait`(duration): 仮想マシンのプロビジョニング開始までの起動からの待ち時間。`10s`、`1m`のように指定する。デフォルト値:`0s`

- `boot_command`([]string): 仮想マシン作成後にVNC経由で仮想マシンに送信するキーボード入力。デフォルト値:なし

`boot_command`はPackerの[VMWare Builder(from ISO)](https://www.packer.io/docs/builders/vmware-iso.html)や[Qemu Builder](https://www.packer.io/docs/builders/qemu.html)と
互換性があります。

- `api_client_timeout`(duration): ディスクのコピーやアーカイブ作成待ちなどの、さくらのクラウドAPI呼び出しで利用するタイムアウト時間。`10s`、`1m`のように指定する。デフォルト値:`20m`

- SSH関連オプション: SSHキー関連の挙動を指定します。

  - `disable_generate_public_key`: trueの場合、秘密鍵に対応する公開鍵の生成/さくらのクラウドAPIを通じての公開鍵アップロードを行いません。
  - `ssh_private_key_file`: 秘密鍵ファイルのパス

Note: `ssh_private_key_file`が未指定の場合、`packer build`実行時に秘密鍵/公開鍵が生成されます。  
`ssh_private_key_file`を指定した場合、指定された秘密鍵から公開鍵を生成します。

生成されたSSH公開鍵は作成するサーバが[ディスクの修正 API](https://manual.sakura.ad.jp/cloud/storage/modifydisk/about.html)に対応している場合はディスクの修正APIを通じてサーバに登録されます。  
この挙動は`disable_generate_public_key`オプションで制御可能です。	


### ISOイメージ関連項目の指定について

ISOイメージからの構築を行う場合、以下の項目を指定してください。

#### さくらのクラウド上のISOイメージを利用する場合

  - `os_type`に`iso`を指定
  - `iso_id`に利用したいISOイメージのIDを指定

例:
```json
        "os_type": "iso",
        "iso_id": 123456789012,
```

#### ISOイメージをダウンロードして利用する場合

ダウンロード元のURLなどを以下のように指定します。
PackerがISOイメージのダウンロードを行い、さくらのクラウド上へアップロードを行います。

  - `os_type`に`iso`を指定
  - `iso_url`または`iso_urls` にISOイメージのURLを指定
  - `iso_checksum`/`iso_checksum_url`: ISOイメージのチェックサム、またはチェックサムを記載したファイルのURLを指定
  - `iso_checksum_type`: チェックサムの形式(`none`, `md5`, `sha1`, `sha256`, or `sha512`のいずれか)を指定

例:
```json
        "os_type": "iso",
        "iso_url": "http://ftp.tsukuba.wide.ad.jp/software/vyos/iso/release/1.1.7/vyos-1.1.7-amd64.iso",
        "iso_checksum": "c40a889469e0eea43d92c73149f1058e3650863b",
        "iso_checksum_type": "sha1",
```

### boot_commandについて

`boot_command`には通常の文字に加え、以下の特殊キーが指定可能です。

- `<bs>` - バックスペース
- `<del>` - デリート
- `<enter>` `<return>` - エンター(リターン)キー
- `<esc>` - エスケープ
- `<tab>` - タブ
- `<f1>` - `<f12>` - ファンクションキー(F1〜F12)
- `<up>` `<down>` `<left>` `<right>` - 矢印キー
- `<spacebar>` - スペース
- `<insert>` - インサート
- `<home>` `<end>` - ホーム、エンド
- `<pageUp>` `<pageDown>` - ページアップ、ページダウン
- `<leftAlt>` `<rightAlt>` - 左右それぞれのオルト
- `<leftCtrl>` `<rightCtrl>` - 左右それぞれのコントロール
- `<leftShift>` `<rightShift>` - 左右それぞれのシフト
- `<leftAltOn>` `<rightAltOn>` - オルトを押下している状態にする
- `<leftCtrlOn>` `<rightCtrlOn>` - コントロールを押下している状態にする
- `<leftShiftOn>` `<rightShiftOn>` - シフトを押下している状態にする
- `<leftAltOff>` `<rightAltOff>` - オルトを押下解除する
- `<leftCtrlOff>` `<rightCtrlOff>` - コントロールを押下解除する
- `<leftShiftOff>` `<rightShiftOff>` - シフトを押下解除する
- `<wait>` `<wait5>` `<wait10>` - 指定秒数待機
- `<waitXX>` - 指定時間待機する。待機数量 + 単位で指定する。例: `<wait10s>` `<wait1m>`

#### 使用例

この例は以下ののキー入力を行うものです。
 - 1) Ctrl+Alt+Delを送信
 - 2) 10秒待機
 - 3) パスワード文字列("put-your-password") + Enterキーを送信

```json
"boot_command": [
    "<leftAltOn><leftCtrlOn><del><leftAltOff><leftCtrlOff>",
    "<wait10>",
    "put-your-password<enter>"
]
```
####

その他の利用例や詳細は以下のテンプレートサンプルを参照してください。

## テンプレートサンプル

以下のサンプルを用意しています。

### パブリックアーカイブからの構築サンプル

 - [[CentOS]](examples/centos): CentOSパブリックアーカイブからの構築
 - [[CoreOS]](examples/coreos): CoreOSパブリックアーカイブからの構築
 - [[Ubuntu]](examples/ubuntu): Ubuntuパブリックアーカイブからの構築
 - [[VyOS]](examples/vyos): VyOSパブリックアーカイブから`boot_command`でSSHを有効にする構成
 
    VyOSパブリックアーカイブからサーバーを作成した場合、デフォルトの状態ではSSH接続ができないため、`boot_command`にてSSHを有効化しています。
    
 - [[Windows]](examples/windows): Windows Server 2012パブリックアーカイブから、`boot_command`でWinRMを有効にする構成
 - [[Windows2016]](examples/windows2016): Windows Server 2016パブリックアーカイブから、`boot_command`でWinRMを有効にする構成
 
    Windowsパブリックアーカイブから構築したサーバーに対し、`boot_command`にて初回ログイン〜WinRMの有効化までを行います。

    WinRMの有効化にはAnsibleが公開しているPowerShellスクリプトを用いています。
    
    WinRMでの接続後にPowerShellでのプロビジョニングを行います。
 

### ISOイメージからの構築サンプル

 - [[CentOS Atomic Host]](examples/atomic_host_centos) : CentOS Atomic Host ISOイメージのダウンロード〜最小構成での構築
 - [[CentOS]](examples/centos_iso): さくらのクラウド上のCentOS ISOイメージからの構築(ISOイメージをダウンロードするサンプルもあります)
 - [[Scientific Linux]](examples/scientific_linux): さくらのクラウド上のScienfitic Linux ISOイメージからの構築
 - [[Ubuntu]](examples/ubuntu_iso): さくらのクラウド上のUbuntu ISOイメージからの構築
 - [[VyOS]](examples/vyos_iso): VyOSISOイメージのダウンロード〜構築(参考元:https://github.com/higebu/packer-templates)


## License

  `packer-builder-sakuracloud` Copyright (C) 2016-2020 Kazumichi Yamamoto.

  This project is published under [MPL-2.0](LICENSE).
  
## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))

