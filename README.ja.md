# **JSON to Table (json-to-table)**

このプロジェクトは **nlink-jp** と **Google の Gemini** の協同開発です。

[English README is here](README.md)

`json-to-table`は、[`nlink-jp/splunk-cli`](https://github.com/nlink-jp/splunk-cli) のコンパニオンツールとして開発された、Go言語製の汎用的なコマンドライン補助ツールです。JSON配列を整形されたテーブルとして出力します。標準入力からJSONデータを受け取るため、`splunk-cli ... | jq .results`のようなコマンドの出力を直接パイプして、人間に読みやすい形式や、レポートに貼り付けやすい画像形式に変換することを主な目的としています。

変更点の詳細については、[CHANGELOG](CHANGELOG.md)をご覧ください。

## **特徴**

*   **汎用的な入力**: 標準入力から、オブジェクトのJSON配列を受け取ります。
*   **多彩な出力形式**:
    *   `text`: ターミナル表示に適した、罫線付きのプレーンテキスト形式。
    *   `md`: GitHub Flavored Markdown形式のテーブル。
    *   `csv`: コンマ区切り形式。スプレッドシートでの利用に適しています。
    *   `png`: **日本語対応の画像形式**。レポートやチャットでの共有に最適です。
    *   `html`: 基本的なスタイルが適用された自己完結型のHTMLファイル。
    *   `slack-block-kit`: SlackのBlock Kit形式のJSON出力。Slackメッセージで直接利用するのに最適です。
*   **柔軟なカラムの選択と順序付け**:
    *   **カラムの包含**: `--columns` (`-c`) フラグで、表示するカラムとその順序を自由に指定できます。
    *   **カラムの除外**: `--exclude-columns` (`-e`) フラグで、出力から除外するカラムを指定できます。
    *   包含と除外の両方で、`*`（残りすべて）や`prefix*`（前方一致）といった強力なワイルドカードをサポートします。
*   **画像カスタマイズ**:
    *   `--title`で画像にタイトルを追加できます。
    *   `--font-size`で文字の大きさを調整できます。
*   **自己完結型**: 日本語フォントをバイナリに埋め込んでいるため、外部ファイルへの依存がなく、単一の実行可能ファイルとして動作します。

## **インストール**

macOS、Windows、Linux向けのコンパイル済みバイナリは[リリースページ](https://github.com/nlink-jp/json-to-table/releases)から入手できます。

## **使い方**

### **基本的なパイプライン**

splunk-cliの出力をjqで絞り込み、その結果をjson-to-tableに渡すのが基本的な使い方です。

```bash
# splunk-cliの結果をテキスト形式のテーブルで表示
splunk-cli run --silent -spl "..." | jq .results | json-to-table
```

### **サンプルデータでの使用**

テスト用に提供されている`testdata/test_data.json`を以下のように使用できます：

```bash
cat testdata/test_data.json | json-to-table
```

### **出力形式の指定**

`--format`フラグで出力形式を変更できます。

*   **Markdown形式でファイルに出力:**
    ```bash
    splunk-cli run ... | jq .results | json-to-table --format md -o report.md
    ```

*   **PNG画像形式でファイルに出力:**
    ```bash
    splunk-cli run ... | jq .results | json-to-table --format png --title "DNS Query Ranking" -o report.png
    ```

*   **HTML形式でファイルに出力:**
    ```bash
    splunk-cli run ... | jq .results | json-to-table --format html -o report.html
    ```

*   **Slack Block Kit形式で出力:**
    ```bash
    splunk-cli run ... | jq .results | json-to-table --format slack-block-kit
    ```

*   **CSV形式でファイルに出力:**
    ```bash
    splunk-cli run ... | jq .results | json-to-table --format csv -o report.csv
    ```

### **カラムの選択と順序付け**

`json-to-table`は、カラムの選択を2つの段階で処理します。まず除外、次いで包含です。

#### **1. カラムの除外 (`--exclude-columns` または `-e`)**

利用可能なカラムの初期セットから削除するカラム名またはパターンを指定します。ワイルドカードは`--columns`と同様に動作します。

*   **特定のカラムを除外:**
    ```bash
    ... | json-to-table -e "id,timestamp"
    ```
    （出力から`id`と`timestamp`を除外します。）

*   **プレフィックスでカラムを除外:**
    ```bash
    ... | json-to-table -e "http_*,_internal*"
    ```
    （`http_`または`_internal`で始まるすべてのカラムを除外します。）

*   **すべてのカラムを除外（注意して使用してください。空のテーブルになります）:**
    ```bash
    ... | json-to-table -e "*"
    ```

#### **2. カラムの包含と順序付け (`--columns` または `-c`)**

除外が適用された後、このフラグを使用して、*残りの*カラムのうちどれを表示するか、そしてその順序を指定します。ワイルドカードは柔軟な順序付けを可能にします。

*   **特定のカラムを先頭に、残りをその後に表示:**
    ```bash
    ... | json-to-table -c "user,*"
    ```

*   **特定のカラムを先頭と末尾に配置:**
    ```bash
    ... | json-to-table -c "user,*,count,total"
    ```

*   **プレフィックスでカラムをグループ化:**
    `http_`で始まるすべてのカラムをまとめて表示します。
    ```bash
    ... | json-to-table -c "user,http_*,*"
    ```

*   **定義された順序で特定のカラムセットのみを表示:**
    ```bash
    ... | json-to-table -c "user,action,status"
    ```

#### **組み合わせた使用例**

最初に`_internal_id`と`timestamp`を除外し、次に`user`、`action`、および残りのすべてのカラムを表示する場合：

```bash
... | json-to-table -e "_internal_id,timestamp" -c "user,action,*"
```

### **フラグ一覧**

*   `--format`: 出力形式 (`text`, `md`, `csv`, `png`, `html`, `slack-block-kit`, `blocks`)。デフォルトは`text`。
*   `-o <file>`: 出力先のファイルパス。デフォルトは標準出力。
*   `--columns, -c <order>`: 包含するカラムと希望する順序をカンマ区切りで指定。`*`は残りのカラムのワイルドカードとして使用。
*   `--exclude-columns, -e <order>`: 除外するカラムをカンマ区切りで指定。`*`はワイルドカードとして使用。
*   `--title <text>`: PNG出力時のタイトル。
*   `--font-size <number>`: PNG出力時のフォントサイズ。デフォルトは12。
*   `--version`: バージョン情報を表示して終了します。

## **ソースからのビルド**

ソースからビルドするには、Goと`make`がインストールされている必要があります。

1.  **リポジトリをクローン:**
    ```bash
    git clone https://github.com/nlink-jp/json-to-table.git
    cd json-to-table
    ```

2.  **バイナリのビルド:**
    ```bash
    make build
    ```
    コンパイルされたバイナリは`dist`ディレクトリに配置されます。

3.  **リリース用パッケージ（ZIP）の作成:**
    ```bash
    make package
    ```
    各OS向けのZIPアーカイブが`dist`ディレクトリに作成され、GitHubリリースにそのまま添付できます。

## **謝辞**

このツールは **Mplus 1 Code** フォントを使用しています。このフォントは、SIL Open Font License, Version 1.1 のもとでライセンスされています。素晴らしいフォントを提供してくださった M+ FONTS Project に感謝します。

## **ライセンス**

このプロジェクトはMITライセンスのもとで公開されています。詳細は[LICENSE](LICENSE)ファイルをご覧ください。
