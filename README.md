# OnlineChatMessenger

## 概要

CLIで簡単なグループチャットができるアプリケーションです。
<br>
コンピュータサイエンス学習サービス[Recursion](https://recursion.example.com)の課題でPythonで作成したものをGoでも作成してみました。


## 機能

- クライアントがサーバに接続し、メッセージを送信すると他の全クライアントに配信されます。
- クライアントは新しいチャットルームを作成するか、既存のチャットルームに名前を指定して参加することができます。
- チャットルームはパスワードで保護されており、正しいパスワードを入力しないと参加できません。
- チャットルームの作成と参加にはTCP接続を使用し、メッセージのやり取りにはUDPを使用しています。

## 目的
基本的なネットワーク通信（TCP/UDP、ファイル操作などのOS機能の理解を目的としています。


## 実行方法

1.サーバを起動
```sh
go run server.go
```  
<img width="600" alt="online_chat_messenger1" src="https://github.com/user-attachments/assets/c8296e97-9807-4473-9372-711f9a10b30a">

2.クライアントを起動
```sh
go run client.go
```  
<img width="600" alt="online_chat_messenger2" src="https://github.com/user-attachments/assets/6b29d583-e44a-4453-8992-3238e18c0a61">
  
3.アドレスを入力<br>
今回はローカル環境なので、0.0.0.0を入力<br>
<img width="600" alt="online_chat_messenger3" src="https://github.com/user-attachments/assets/0efdb8ac-356a-4b08-8f3d-acd3de64e370">
<br>

4.「create」と入力<br>
<img width="600" alt="online_chat_messenger4" src="https://github.com/user-attachments/assets/aab4ac4d-e630-4fd0-a245-42f7b29651f7">
<br>

5.作成したいルーム名を入力<br> 
<img width="600" alt="online_chat_messenger5" src="https://github.com/user-attachments/assets/a1eeb829-5bd2-4e3a-bc6a-58c708683824">
<br>

6.作成するルームのパスワードを設定<br> 
<img width="600" alt="online_chat_messenger6" src="https://github.com/user-attachments/assets/3f7a3521-0749-48c7-acd5-b9a380b5f8eb">
<br>

7.作成するユーザー名を入力<br>
<img width="600" alt="online_chat_messenger7" src="https://github.com/user-attachments/assets/8b867191-9c63-4636-a890-5b60e92b054e">
<br>

8.サーバのターミナルに作成が成功した旨のメッセージが表示される<br>
<img width="600" alt="online_chat_messenger8" src="https://github.com/user-attachments/assets/646adcc9-3781-4c6c-9d30-241091d7ddfa">
<br>

9.別のターミナルでクライアントをもう一つ起動させ、新しいユーザーで作成したルームに参加する<br>
<img width="600" alt="online_chat_messenger9" src="https://github.com/user-attachments/assets/b4a53aea-9b4f-476c-9a2a-efccce5e5bf7">
<br>

10.user1とuser2でチャットを行う<br>
<img width="600" alt="online_chat_messenger10" src="https://github.com/user-attachments/assets/1e160858-f6b8-40cf-8f1a-f80ed2cbb96e">
<br>
<img width="600" alt="online_chat_messenger11" src="https://github.com/user-attachments/assets/8bacd9cd-828b-4e7f-83f5-cde8f72c1ed4">
