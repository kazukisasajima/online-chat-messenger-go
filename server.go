package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
	"sync"
)

// ChatRoom構造体：チャットルームの情報を保持
type ChatRoom struct {
	Name      string                 // チャットルーム名
	Password  string                 // チャットルームのパスワード
	// Mutexは複数のごルーチンが同時に同じデータにアクセスするのを防ぐために使用
	// lock/unlockを使って排他制御を行う
	Members   map[string]*net.UDPAddr // トークン: クライアントアドレス
	Mutex     sync.Mutex             // メンバー操作のためのミューテックス
}

var (
	chatRooms = make(map[string]*ChatRoom) // チャットルームのマップ（ルーム名: ChatRoom）
	roomMutex sync.Mutex                   // ルーム操作のためのミューテックス
)

const serverAddress = "0.0.0.0:9001"

func main() {
	// TCPとUDPのリスナーを初期化
	tcpAddr, _ := net.ResolveTCPAddr("tcp", serverAddress)
	udpAddr, _ := net.ResolveUDPAddr("udp", serverAddress)

	// TCPリスナーを開始
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("Failed to start TCP listener: %v", err)
	}
	defer tcpListener.Close()

	// UDPリスナーを開始
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to start UDP listener: %v", err)
	}
	defer udpConn.Close()

	log.Println("Server is running...")

	// 並列でTCP接続を処理
	go handleTCPConnections(tcpListener)

	// UDP接続を処理（メインスレッド）
	handleUDPConnection(udpConn)
}

func handleTCPConnections(tcpListener *net.TCPListener) {
	for {
		// TCP接続を受け入れる
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}
		// 新しいTCP接続を別のゴルーチンで処理
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	// ヘッダー情報を読み取る（固定長4バイト）
	header := make([]byte, 4)
	_, err := conn.Read(header)
	if err != nil {
		log.Printf("Failed to read TCP header: %v", err)
		return
	}

	// ヘッダーから各フィールドを解析
	roomNameSize := int(header[0])   // ルーム名のサイズ
	action := header[1]             // アクション（1: 作成, 2: 参加）
	passwordSize := int(header[2])  // パスワードのサイズ
	usernameSize := int(header[3])  // ユーザー名のサイズ

	// 各フィールドのデータを読み取る
	roomNameBytes := make([]byte, roomNameSize)
	conn.Read(roomNameBytes)
	roomName := string(roomNameBytes)

	passwordBytes := make([]byte, passwordSize)
	conn.Read(passwordBytes)
	password := string(passwordBytes)

	usernameBytes := make([]byte, usernameSize)
	conn.Read(usernameBytes)
	username := string(usernameBytes)

	// ログにリクエスト内容を記録
	log.Printf("TCP Request - Action: %d, Room: '%s', Password: '%s', Username: '%s'", action, roomName, password, username)

	// 入力検証
	if roomName == "" || password == "" || username == "" {
		log.Println("Invalid request: Room name, password, or username is empty.")
		conn.Write([]byte("Error: Invalid request"))
		return
	}

	// チャットルーム操作を排他的に処理
	roomMutex.Lock()
	defer roomMutex.Unlock()

	switch action {
	case 1: // ルーム作成
		if _, exists := chatRooms[roomName]; exists {
			log.Printf("Room creation failed: Room '%s' already exists", roomName)
			conn.Write([]byte("Error: Room already exists"))
			return
		}
		// 新しいルームを作成
		chatRooms[roomName] = &ChatRoom{
			Name:     roomName,
			Password: password,
			Members:  make(map[string]*net.UDPAddr),
		}
		log.Printf("Room '%s' created successfully by user '%s'", roomName, username)
		conn.Write([]byte("Room created successfully"))
	case 2: // ルーム参加
		room, exists := chatRooms[roomName]
		if !exists {
			log.Printf("Room join failed: Room '%s' does not exist", roomName)
			conn.Write([]byte("Error: Room does not exist"))
			return
		}
		if room.Password != password {
			log.Printf("Room join failed: Incorrect password for room '%s'", roomName)
			conn.Write([]byte("Error: Incorrect password"))
			return
		}
		token := generateToken() // ユニークなトークンを生成
		room.Members[token] = nil
		log.Printf("User '%s' joined room '%s' with token '%s'", username, roomName, token)
		conn.Write([]byte(token))
	default: // 不明なアクション
		log.Printf("Invalid action received: %d", action)
		conn.Write([]byte("Error: Invalid action"))
	}
}

func handleUDPConnection(conn *net.UDPConn) {
	buf := make([]byte, 4096)
	for {
		// UDPデータを受信
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading UDP data: %v", err)
			continue
		}

		// ヘッダー解析
		roomNameSize := int(buf[0])
		tokenSize := int(buf[1])
		roomName := string(buf[2 : 2+roomNameSize])
		token := string(buf[2+roomNameSize : 2+roomNameSize+tokenSize])
		message := string(buf[2+roomNameSize+tokenSize : n])

		log.Printf("UDP Message - Room: '%s', Token: '%s', From: %s, Message: '%s'", roomName, token, clientAddr, message)

		// チャットルーム存在確認
		roomMutex.Lock()
		room, exists := chatRooms[roomName]
		roomMutex.Unlock()

		if !exists {
			log.Printf("Message discarded: Room '%s' does not exist", roomName)
			continue
		}

		// メッセージを他のクライアントにリレー
		room.Mutex.Lock()
		clientTokenAddr, valid := room.Members[token]
		if !valid || clientTokenAddr == nil {
			room.Members[token] = clientAddr
		}
		for tok, addr := range room.Members {
			if tok != token && addr != nil {
				log.Printf("Relaying message from %s to %s", clientAddr, addr)
				conn.WriteToUDP([]byte(message), addr)
			}
		}
		room.Mutex.Unlock()
	}
}

// トークンを生成する関数
func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
