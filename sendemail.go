package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// Credentials holds info of mail
type Credentials struct {
	serverAddr, password, emailAddr, delimeter string
	tos, cc                                    *[]string
	portNumber                                 string
}

func genCred() *Credentials {
	return &Credentials{
		serverAddr: os.Getenv("SERVERADDRESS"),
		password:   os.Getenv("PASSWORD"),
		emailAddr:  os.Getenv("FROM"),
		tos:        &[]string{os.Getenv("TO")},
		cc:         &[]string{os.Getenv("CC")},
		portNumber: os.Getenv("PORT"),
		delimeter:  os.Getenv("DELIMETER"),
	}
}

func (cred *Credentials) tlsconfigs() *tls.Config {
	return &tls.Config{
		ServerName:         cred.serverAddr,
		InsecureSkipVerify: true,
	}
}

func (cred *Credentials) tlsConnect() *tls.Conn {
	// log.Println("Establish TLS connection")
	conn, connErr := tls.Dial("tcp", fmt.Sprintf("%s:%s", cred.serverAddr, cred.portNumber), cred.tlsconfigs())
	if connErr != nil {
		log.Printf("error to connect %v\n", connErr)
		return nil
	}
	return conn
}

func (cred *Credentials) createClient(conn *tls.Conn) *smtp.Client {
	// log.Println("create new email client")
	client, clientErr := smtp.NewClient(conn, cred.serverAddr)
	if clientErr != nil {
		log.Printf("error to create client %v\n", clientErr)
		return nil
	}
	return client
}

func (cred *Credentials) authenticate() smtp.Auth {
	return smtp.PlainAuth("", cred.emailAddr, cred.password, cred.serverAddr)
}
func setupAuth(client *smtp.Client, auth smtp.Auth) {
	err := client.Auth(auth)
	if err != nil {
		log.Printf("error to connect client with auth %v\n", err)
		log.Panic(err)
	}
}

func (cred *Credentials) setFrom(client *smtp.Client) {
	err := client.Mail(cred.emailAddr)
	if err != nil {
		// log.Panic(err)
		log.Printf("error to setup from %v\n", err)
		return
	}
}

func (cred *Credentials) setupTo(client *smtp.Client) {
	for _, to := range *cred.tos {
		err := client.Rcpt(to)
		if err != nil {
			log.Printf("error to setup to emails %v\n", err)
			return
		}
	}
}

func (cred *Credentials) setupEmailHeader(attachmentFilePath, filename string) string {
	//basic email headers
	sampleMsg := fmt.Sprintf("From: %s\r\n", cred.emailAddr)
	sampleMsg += fmt.Sprintf("To: %s\r\n", strings.Join(*cred.tos, ";"))
	if len(*cred.cc) > 0 {
		sampleMsg += fmt.Sprintf("Cc: %s\r\n", strings.Join(*cred.cc, ";"))
	}
	sampleMsg += "Subject: Captured image\r\n"

	sampleMsg += "MIME-Version: 1.0\r\n"
	sampleMsg += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", cred.delimeter)

	//place HTML message
	// log.Println("Put HTML message")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", cred.delimeter)
	sampleMsg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: 7bit\r\n"
	sampleMsg += fmt.Sprintf("\r\n%s", "<html><body><h1>Hi Sir</h1>"+
		"<p>this is activity checker result</p></body></html>\r\n")

	//place file
	// log.Println("Put file attachment")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", cred.delimeter)
	sampleMsg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: base64\r\n"
	sampleMsg += "Content-Disposition: attachment;filename=\"" + filename + "\"\r\n"
	//read file
	rawFile := readFile(attachmentFilePath)

	sampleMsg += "\r\n" + base64.StdEncoding.EncodeToString(rawFile)
	return sampleMsg
}

func readFile(attachmentFilePath string) []byte {
	rawFile, fileErr := ioutil.ReadFile(attachmentFilePath)
	if fileErr != nil {
		log.Printf("error to read file %v\n", fileErr)
		return nil
	}
	return rawFile
}

func mailsPrepared(folderName, fileName string) {
	cred := genCred()
	attachmentFilePath := fmt.Sprintf("./%s/%s", folderName, fileName)
	filename := fileName

	log.Println("NOTE: user need to turn on 'less secure apps' options")
	log.Println("URL:  https://myaccount.google.com/lesssecureapps\n\r")

	conn := cred.tlsConnect()
	defer conn.Close()

	client := cred.createClient(conn)
	defer client.Close()

	// log.Println("setup authenticate credential")
	auth := cred.authenticate()

	setupAuth(client, auth)

	// log.Println("Start write mail content")
	cred.setFrom(client)

	cred.setupTo(client)

	writer, writerErr := client.Data()
	if writerErr != nil {
		log.Printf("error to write data to client %v\n", writerErr)
		return
	}

	sampleMsg := cred.setupEmailHeader(attachmentFilePath, filename)

	_, err := writer.Write([]byte(sampleMsg))
	if err != nil {
		log.Printf("error to write header into email %v\n", err)
		return
	}

	closeErr := writer.Close()
	if closeErr != nil {
		log.Printf("error to close writer %v\n", closeErr)
		return
	}

	client.Quit()
	// log.Println(sampleMsg)
	log.Print("sent :))")

}
