package main

import (
	"testing"
)

var bodies map[string]string = map[string]string{
	`MIME-Version: 1.0
Message-ID: <CANm3GnZD8Tnz5h89OD60P=yt-fBTC77AHUzxW2kV7vdEmQ1rLg@mail.gmail.com>
Date: Mon, 05 Sep 2016 10:02:33 +0000
Subject: tgest 2
From: Me Me <me@me.com>
To: gateway+adsfdasdg2efsg-+123123@txt.me.im
Content-Type: multipart/alternative; boundary=001a114032de63e68a053bbfc726

--001a114032de63e68a053bbfc726
Content-Type: text/plain; charset=UTF-8; format=flowed; delsp=yes

yup

--001a114032de63e68a053bbfc726
Content-Type: text/html; charset=UTF-8

<div dir="ltr">yup</div>
--001a114032de63e68a053bbfc726--`: "yup",
	`Header: nope

 yup `: "yup",
	`Delivered-To: me.me@gmail.com
Received: by 10.25.210.73 with SMTP id j70csp718235lfg;
        Wed, 31 Aug 2016 23:44:12 -0700 (PDT)
X-Received: by 10.194.20.65 with SMTP id l1mr14162036wje.71.1472712252202;
        Wed, 31 Aug 2016 23:44:12 -0700 (PDT)
Return-Path: <owner-discuss-medot-org@medot.org>
Received: from mail.me.com (mail.me.com. [1.1.1.1])
        by mx.google.com with ESMTPS id v11si4178335wjr.95.2016.08.31.23.44.12
        for <me.me@gmail.com>
        (version=TLS1_2 cipher=AES128-SHA bits=128/128);
        Wed, 31 Aug 2016 23:44:12 -0700 (PDT)
Received-SPF: neutral (google.com: 1.1.1.1 is neither permitted nor denied by best guess record for domain of owner-discuss-medot-org@medot.org) client-ip=1.1.1.1;
Authentication-Results: mx.google.com;
       spf=neutral (google.com: 1.1.1.1 is neither permitted nor denied by best guess record for domain of owner-discuss-medot-org@medot.org) smtp.mailfrom=owner-discuss-medot-org@medot.org
Received: from asdg.asdf.com ([85.13.151.217])
	by mail.me.com with esmtps (TLS1.2:DHE_RSA_AES_256_CBC_SHA256:256)
	(Exim 4.80)
	(envelope-from <owner-discuss-medot-org@medot.org>)
	id 1bfLjK-00034g-GX
	for c+medot@me.im; Thu, 01 Sep 2016 08:44:12 +0200
Received: by asdg.asdf.com (Postfix)
	id 0A9AA8009A8; Thu,  1 Sep 2016 08:44:05 +0200 (CEST)
Delivered-To: discuss-medot-org-list@asdg.asdf.com
Received: by asdg.asdf.com (Postfix, from userid 65534)
	id ED183802B8C; Thu,  1 Sep 2016 08:44:04 +0200 (CEST)
Delivered-To: discuss-medot-org@asdg.asdf.com
Received: from dd6600.asdf.com (dd6600.asdf.com [85.13.131.77])
	by asdg.asdf.com (Postfix) with ESMTPS id 1EBB88009A8
	for <discuss-medot-org@medot.org>; Thu,  1 Sep 2016 08:44:04 +0200 (CEST)
Received: from [10.42.217.52] (unknown [89.204.135.52])
	by dd6600.asdf.com (Postfix) with ESMTPSA id A72DA16802B2
	for <discuss-medot-org@medot.org>; Thu,  1 Sep 2016 08:44:03 +0200 (CEST)
References: <780337603.209302.1472629632038.JavaMail.ngmail@nope.net> <003801d20362$9d1ede90$d75c9bb0$@yup.de> <2c2ba3cd-4e59-a5ab-7a05-1d1bc6f8dfa6@foo.bar> <C3B13662-9399-4A41-BD00-16C901587136@bar.com> <7e469c72-d927-5286-a01f-02f83cc95f3b@foo.bar> <BEC0C3FE-47A8-4BF0-B9F0-2212C4516488@bar.com> <662e0718-07f2-3ab1-2ef6-2feba6851beb@foo.bar>
From: nope <nope@jup.net>
X-Mailer: iPhone Mail (13G36)
In-Reply-To: <662e0718-07f2-3ab1-2ef6-2feba6851beb@foo.bar>
Message-Id: <64858C8C-6927-44BA-99A4-ADE25C9D26C0@jup.net>
Date: Thu, 1 Sep 2016 08:43:57 +0200
To: discuss-medot-org@medot.org
Mime-Version: 1.0 (1.0)
Sender: owner-discuss-medot-org@medot.org
Precedence: bulk
Reply-To: discuss-medot-org@medot.org
X-SA-Exim-Connect-IP: 1.1.1.1
X-SA-Exim-Mail-From: owner-discuss-medot-org@medot.org
X-Spam-Checker-Version: SpamAssassin 3.3.2 (2011-06-06) on vm1.me.com
X-Spam-Level: 
X-Spam-Status: No, score=-1.9 required=5.0 tests=BAYES_00 autolearn=ham
	version=3.3.2
Subject: =?utf-8?Q?Re:_[d-fd]_Re:_[d-fd]_AW:_[d-fd]_Forum,_was_f=C3=BCr_e?=
 =?utf-8?Q?in_Forum=3F_Wo_gibt_es_Infos_dazu=3F?=
Content-Type: text/plain;
	charset=utf-8
Content-Transfer-Encoding: quoted-printable
X-SA-Exim-Version: 4.2.1 (built Mon, 26 Dec 2011 16:24:06 +0000)
X-SA-Exim-Scanned: Yes (on mail.me.com)

This is a test. just a test.`: "This is a test. just a test.",
}

func TestSendBody(t *testing.T) {

	for source, exp := range bodies {
		header, body := splitHeaders(source)
		t.Log("--------- Source: ----------\n" + source)
		t.Log("--------- Headers: ----------\n", header)
		t.Log("--------- Parsed Body: ----------\n", body)
		got := getClearBody(Message{Body: body, Headers: header})
		t.Log("--------- Got: ----------\n", got)
		if exp != got {
			t.Errorf("Got '%s', expected '%s' on source '%s'", got, exp, source)
		} else {
			t.Log("---- OK -----")
		}
	}

}
