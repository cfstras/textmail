// +build goguerilla

package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

const TestMessage = `HELO test
MAIL FROM: P.A.S.S+caf_=gateway+test-+bla=foo.bar.com@MYDOMAIN.com
RCPT TO: test+test@localhost
DATA
Received: by server-bar-01.MYDOMAIN.com (Postfix, from userid 0)
	id 64ED61F9B8; Thu, 14 Jan 2016 17:05:03 -0700 (MST)
To: SOMEBODY@MYDOMAIN.com
Subject: (primary)Test: 'SSH - server-foo-01.MYDOMAIN.com' is DOWN
X-PHP-Originating-Script: 48:phpmailer.class.php
Received: from phpmailer ([]) by  with HTTP;
	 Thu, 14 Jan 2016 17:05:03 -0700
Date: Thu, 14 Jan 2016 17:05:03 -0700
From: P.A.S.S+caf_=gateway+test-+bla=foo.bar.com@MYDOMAIN.com
X-Priority: 1
X-Mailer: phpmailer [version 1.54]
MIME-Version: 1.0
Content-Type: multipart/alternative;
	boundary="b1_755f9caf25cbc2267e0bd344b189291a"
Message-Id: <20160115000503.64ED61F9B8@server-bar-01.MYDOMAIN.com>

--b1_755f9caf25cbc2267e0bd344b189291a
Content-Type: text/plain; charset = "iso-8859-1"

PASS

--b1_755f9caf25cbc2267e0bd344b189291a
Content-Type: text/html; charset = "iso-8859-1"

FAIL

--b1_755f9caf25cbc2267e0bd344b189291a--

.
QUIT
`

func TestListen(t *testing.T) {
	os.Args = []string{os.Args[0], "-v", "y"}
	gConfig["GM_ALLOWED_HOSTS"] = "localhost"
	gConfig["GSTMP_LISTEN_INTERFACE"] = "0.0.0.0:2525"
	gConfig["GSMTP_PRV_KEY"] = ""
	gConfig["GSMTP_PUB_KEY"] = ""

	wait := sync.WaitGroup{}
	wait.Add(1)
	SetCallback(func(msg Message) {
		body := getClearBody(msg)
		t.Logf("Got: %+v\n", msg)
		if body != "PASS" {
			t.Log("ClearBody should be PASS, but is", body)
			t.Fail()
		}
		addr := getClearAddress(msg.From)
		if addr != "P.A.S.S@MYDOMAIN.com" {
			t.Log("ClearBody should be P.A.S.S@MYDOMAIN.com, but is", addr)
			t.Fail()
		}
		wait.Done()
	})
	go mainLoop()

	go func() {
		time.Sleep(20 * time.Second)
		wait.Done()
		t.Fatal("Timeout")
	}()

	time.Sleep(2 * time.Second)
	conn, err := net.Dial("tcp", "localhost:2525")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(conn, TestMessage)

	conn.Close()

	wait.Wait()
}
