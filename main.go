package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings" // [추가] 이 줄을 import 문에 추가하세요.
)

// 요청 본문을 파싱하기 위한 구조체입니다.
type EmailRequest struct {
	To      string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	// "/email" 경로로 오는 POST 요청을 처리할 핸들러를 등록합니다.
	http.HandleFunc("/email", emailHandler)

	// 8080 포트에서 웹 서버를 시작합니다.
	fmt.Println("서버가 5525 포트에서 실행 중입니다...")
	log.Fatal(http.ListenAndServe(":5525", nil))
}

func emailHandler(w http.ResponseWriter, r *http.Request) {
	// POST 요청만 허용합니다.
	if r.Method != http.MethodPost {
		http.Error(w, "POST 요청만 지원합니다.", http.StatusMethodNotAllowed)
		return
	}

	// 요청 본문(JSON)을 디코딩하여 EmailRequest 구조체에 담습니다.
	var req EmailRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "잘못된 요청 본문입니다.", http.StatusBadRequest)
		return
	}

	// 이메일 전송에 필요한 정보를 설정합니다.
	from := "email"        // 보내는 사람 이메일 주소
	password := "password" // Gmail 앱 비밀번호

	// SMTP 서버 정보를 설정합니다. (Gmail 기준)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// [이 부분만 수정됩니다]
	// 1. 수신된 본문의 줄바꿈(\n)을 HTML 줄바꿈 태그(<br>)로 변경합니다.
	htmlBody := strings.ReplaceAll(req.Body, "\n", "<br>")

	// 2. 이메일 메시지를 구성할 때, 수정된 htmlBody를 사용합니다.
	msg := []byte("To: " + req.To + "\r\n" +
		"Subject: " + req.Subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody + "\r\n") // <-- req.Body 대신 htmlBody 사용

	// SMTP 인증 정보를 설정합니다.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 이메일을 전송합니다.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{req.To}, msg)
	if err != nil {
		log.Printf("이메일 전송 실패: %s", err)
		http.Error(w, "이메일 전송에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	// 성공 응답을 보냅니다.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "이메일이 성공적으로 전송되었습니다: %s", req.To)
	log.Printf("이메일 성공적으로 전송: %s", req.To)
}
