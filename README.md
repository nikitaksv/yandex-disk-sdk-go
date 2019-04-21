<p style="background: white; padding: 25px 0px 0px 0px" align="center">
    <a href="https://disk.yandex.com/" target="_blank" rel="noopener">
        <img src="https://yastatic.net/s3/auth2/_/logo-red_en.1d255bcb.svg" alt="Yandex.Disk" width="400" height="200"/> 
    </a>
     <a href="https://disk.yandex.com/" target="_blank" rel="noopener">
            <img src="https://golang.org/doc/gopher/run.png" alt="Yandex.Disk" width="194" height="180"/> 
    </a>
</p>

Yandex Disk SDK on GoLand

It is a fast, safe and efficient tool that works immediately after installation.

[![Build Status](https://travis-ci.com/nikitaksv/yandex-disk-sdk-go.svg?branch=master)](https://travis-ci.com/nikitaksv/yandex-disk-sdk-go)
[![Coverage Status](https://coveralls.io/repos/github/nikitaksv/yandex-disk-sdk-go/badge.svg)](https://coveralls.io/github/nikitaksv/yandex-disk-sdk-go)
[![CodeFactor](https://www.codefactor.io/repository/github/nikitaksv/yandex-disk-sdk-go/badge)](https://www.codefactor.io/repository/github/nikitaksv/yandex-disk-sdk-go)
-

Installation
------------

Use module (recommended)
```go
import "github.com/nikitaksv/yandex-disk-sdk-go"
```

Use vendor
```sh
go get github.com/nikitaksv/yandex-disk-sdk-go
```

Documentation
-------------

**Useful links on official docs:**

* [Rest API Disk](https://tech.yandex.com/disk/rest/)
* [API Documentation](https://tech.yandex.com/disk/api/concepts/about-docpage/)
* [Try API](https://tech.yandex.com/disk/poligon/)
* [Get Token](https://tech.yandex.com/oauth/) 


Create new instance Yandex.Disk

```go
yaDisk,err := yadisk.NewYaDisk(ctx.Background(),http.DefaultClient, &yadisk.Token{AccessToken: "YOUR_TOKEN"})
if err != nil {
    panic(err.Error())
}
disk,err := yaDisk.GetDisk([]string{})
if err != nil {
    // If response get error
    e, ok := err.(*yadisk.Error)
    if !ok {
        panic(err.Error())
    }
    // e.ErrorID
    // e.Message
}
```

---

<p align="center">
    <b>If the Yandex.Disk SDK helped you, consider donating to the author of this project, Nikita Krasnikov, to show your support. Thanks you!</b>
</p>
<p align="center">
    <a href="https://www.patreon.com/bePatron?u=19197324" data-patreon-widget-type="become-patron-button">
        <img src="https://c5.patreon.com/external/logo/become_a_patron_button@2x.png" width="160" title="Become a Patron!">
    </a>
</p>
