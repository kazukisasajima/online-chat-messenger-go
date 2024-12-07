package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// サーバーアドレスを取得
	serverAddress := getServerAddress(reader)

	// ルームに関するアクション（作成または参加）を選択
	action := getRoomAction(reader)

	// ルーム名、パスワード、ユーザー名を取得
	roomName := getRoomName(reader)
	password := getPassword(reader)
	username := getUsername(reader)

	// TCP接続を確立
	tcpConn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Failed to connect to TCP server: %v", err)
	}
	defer tcpConn.Close()

	// TCPアクションを実行し、サーバーからの応答（トークンまたはエラーメッセージ）を受信
	token := handleTCPAction(tcpConn, action, roomName, password, username)
	if strings.HasPrefix(token, "Error:") {
		// エラーが発生した場合は表示して終了
		fmt.Println(token)
		return
	}

	// UDP接続を確立
	udpConn := connectToUDP(serverAddress)
	defer udpConn.Close()

	// チャットを開始
	chat(udpConn, reader, roomName, token, username)
}

func getServerAddress(reader *bufio.Reader) string {
	// サーバーアドレスを入力させる
	fmt.Print("Enter server address: ")
	address, _ := reader.ReadString('\n')
	return strings.TrimSpace(address) + ":9001" // デフォルトポートを追加
}

func getRoomAction(reader *bufio.Reader) int {
	// ルームのアクションを選択させる（createまたはjoin）
	for {
		fmt.Print("Do you want to create or join a room? (create/join): ")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)
		if action == "create" {
			return 1 // アクション1: 作成
		} else if action == "join" {
			return 2 // アクション2: 参加
		}
		fmt.Println("Invalid action. Please type 'create' or 'join'.")
	}
}

func getRoomName(reader *bufio.Reader) string {
	// ルーム名を入力させる
	fmt.Print("Enter room name: ")
	roomName, _ := reader.ReadString('\n')
	return strings.TrimSpace(roomName)
}

func getPassword(reader *bufio.Reader) string {
	// パスワードを入力させる
	fmt.Print("Enter room password: ")
	password, _ := reader.ReadString('\n')
	return strings.TrimSpace(password)
}

func getUsername(reader *bufio.Reader) string {
	// ユーザー名を入力させる
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	return strings.TrimSpace(username)
}

func handleTCPAction(conn net.Conn, action int, roomName, password, username string) string {
	// TCPリクエストのヘッダーを作成
	header := make([]byte, 4)
	header[0] = byte(len(roomName))  // ルーム名の長さ
	header[1] = byte(action)        // アクション（作成または参加）
	header[2] = byte(len(password)) // パスワードの長さ
	header[3] = byte(len(username)) // ユーザー名の長さ

	// ヘッダーとデータをサーバーに送信
	conn.Write(header)
	conn.Write([]byte(roomName))
	conn.Write([]byte(password))
	conn.Write([]byte(username))

	// サーバーからの応答を受信
	response := make([]byte, 256)
	n, _ := conn.Read(response)

	// 応答を文字列に変換して返す
	return strings.TrimSpace(string(response[:n]))
}

func connectToUDP(serverAddress string) *net.UDPConn {
	// UDP接続を確立
	// :0でポート番号を0にすることで、OSに自動的に空いているポートを割り当てさせる
	localAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	serverAddr, _ := net.ResolveUDPAddr("udp", serverAddress)
	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to UDP server: %v", err)
	}
	return conn
}

func chat(conn *net.UDPConn, reader *bufio.Reader, roomName, token, username string) {
	// メッセージ送信部分を別のゴルーチンで非同期に実行
	go func() {
		for {
			fmt.Print("> ") // プロンプトを表示
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			// メッセージデータを構築（ヘッダー + ボディ）
			data := append([]byte{byte(len(roomName)), byte(len(token))}, []byte(roomName)...)
			data = append(data, []byte(token)...)
			data = append(data, []byte(username+": "+message)...)

			// メッセージをUDPでサーバーに送信
			conn.Write(data)
		}
	}()

	// メッセージ受信部分
	buf := make([]byte, 4096)
	for {
		// サーバーからのメッセージを受信
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}
		// メッセージを出力
		fmt.Println(string(buf[:n]))
		fmt.Print("> ") // 再度プロンプトを表示
	}
}
