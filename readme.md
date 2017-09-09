golang で Twitterライクなサービス
- 使用言語: Golang
 
- How to
gotwitter DB作成 : create database if not exists gotwitter;
Users と Tweets 用のテーブルを作成
 * Users
~~~
create table if not exists users (id INT  AUTO_INCREMENT NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL, display_name VARCHAR(255) NOT NULL, created_at TIMEST    AMP, updated_at TIMESTAMP, email VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL), salt VARCHAR(255) NOT NULL;
~~~
~~~
TIMESTAMPで怒られたら,こっち
~~~
create table if not exists users (id INT  AUTO_INCREMENT NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL, display_name VARCHAR(255) NOT NULL, created_at timest    amp not null default current_timestamp, updated_at timestamp not null default current_timestamp on update current_timestamp, email VARCHAR(255) NOT NULL, passwo    rd VARCHAR(255) NOT NULL,salt VARCHAR(255) NOT NULL);
~~~
* Tweets
~~~
create table if not exists tweets (id INT  AUTO_INCREMENT NOT NULL PRIMARY KEY, user_id INT NOT NULL, created_at TIMESTAMP,  text VARCHAR(255), mention INT);
~~~
 
go (vertion 1.9 で動作確認済み) を install
go get
go run で main.go を叩く (サーバ起動)
ブラウザ等で localhost8080/login にアクセス

- contents
HTMLのテンプレたち

- server_test  
TL 名前表示検証用
