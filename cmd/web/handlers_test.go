package main

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "net/url"
    "log"
    "testing"
    "bytes"
)

func TestPingA(t *testing.T) {
    rr := httptest.NewRecorder()
    r, err := http.NewRequest("GET", "/", nil)

    if err != nil {
        t.Fatal(err)
    }

    ping(rr, r)
    rs := rr.Result()
    if rs.StatusCode != http.StatusOK {
        t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
    }
    defer rs.Body.Close()

    body, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }

    if string(body) != "OK" {
        t.Errorf("want body to equal %q", "OK")
    }
}

func TestPingB(t *testing.T) {
    app := &application{
        errorLog: log.New(ioutil.Discard, "", 0),
        infoLog: log.New(ioutil.Discard, "", 0),
    }

    ts := httptest.NewTLSServer(app.routes())
    defer ts.Close()

    rs, err := ts.Client().Get(ts.URL + "/ping")
    if err != nil {
        t.Fatal(err)
    }

    if rs.StatusCode != http.StatusOK {
        t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
    }
    defer rs.Body.Close()

    body, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }
    if string(body) != "OK" {
        t.Errorf("want body to equal %q", "OK")
    }
}

func TestPingC(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    code, _, body := ts.get(t, "/ping")
    if code != http.StatusOK {
        t.Errorf("want %d; got %d", http.StatusOK, code)
    }
    if string(body) != "OK" {
        t.Errorf("want body to equal %q", "OK")
    }
}

func TestShowSnippet(t *testing.T) {
    app := newTestApplication(t)

    ts := newTestServer(t, app.routes())
    defer ts.Close()

    tests := []struct {
        name string
        urlPath string
        wantCode int
        wantBody []byte
    }{
        {"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
        {"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
        {"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
        {"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
        {"String ID", "/snippet/foo", http.StatusNotFound, nil},
        {"Empty ID", "/snippet/", http.StatusNotFound, nil},
        {"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            code, _, body := ts.get(t, tt.urlPath)
            if code != tt.wantCode {
                t.Errorf("want %d; got %d", tt.wantCode, code)
            }
            if !bytes.Contains(body, tt.wantBody) {
                t.Errorf("want body to contain %q", tt.wantBody)
            }
        })
    }
}

func TestSignupUserA(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    _, _, body := ts.get(t, "/user/signup")
    csrfToken := extractCSRFToken(t, body)

    t.Log(csrfToken)
}

func TestSignupUserB(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    _, _, body := ts.get(t, "/user/signup")
    csrfToken := extractCSRFToken(t, body)
    tests := []struct {
        name string
        userName string
        userEmail string
        userPassword string
        csrfToken string
        wantCode int
        wantBody []byte
    }{
        // These responses are not valid...
        {"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, 200, []byte("")},
        {"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, 200, []byte("")}, 
        {"Empty email", "Bob", "", "validPa$$word", csrfToken, 200, []byte("")}, 
        {"Empty password", "Bob", "bob@example.com", "", csrfToken, 200, []byte("")}, 
        {"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, 200, []byte("")}, 
        {"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, 200, []byte("")},
        {"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, 200, []byte("")},
        {"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, 200, []byte("")},
        {"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, 200, []byte("")},
        {"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, []byte("")},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            form := url.Values{}
            form.Add("name", tt.userName)
            form.Add("email", tt.userEmail)
            form.Add("password", tt.userPassword)
            form.Add("csrf_token", tt.csrfToken)

            code, _, body := ts.postForm(t, "/user/signup", form)
            if code != tt.wantCode {
                t.Errorf("want %d; got %d", tt.wantCode, code)
            }

            if !bytes.Contains(body, tt.wantBody) {
                t.Errorf("want body %s to contain %q", body, tt.wantBody)
            }
        })
    }
}
