package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
    "net/http/cookiejar"
    "time"
    "spbear/snippetbox/pkg/models/mock"
    "regexp"
    "html"
    "github.com/golangcollege/sessions"
    "net/url"
)

type testServer struct {
    *httptest.Server
}

// We create a subexpression with (), in this subexpression we define our regex
// match
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='([a-zA-Z0-9/\+&=;'#]*)'`)

func extractCSRFToken(t *testing.T, body []byte) string {
    matches := csrfTokenRX.FindSubmatch(body)

    if len(matches) < 2 {
        t.Fatal("no csrf token found in body")
    }

    return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {
    templateCache, err := newTemplateCache("./../../ui/html/")
    if err != nil {
        t.Fatal(err)
    }

    session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
    session.Lifetime = 12 * time.Hour
    session.Secure = true

    return &application{
        errorLog: log.New(ioutil.Discard, "", 0),
        infoLog: log.New(ioutil.Discard, "", 0),
        session: session,
        snippets: &mock.SnippetModel{},
        templateCache: templateCache,
        users: &mock.UserModel{},
    }
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
    ts := httptest.NewTLSServer(h)

    jar, err := cookiejar.New(nil)
    if err != nil {
        t.Fatal(err)
    }

    ts.Client().Jar = jar

    ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }

    return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
    rs, err := ts.Client().Get(ts.URL + urlPath)
    if err != nil {
        t.Fatal(err)
    }
    defer rs.Body.Close()

    body, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }

    return rs.StatusCode, rs.Header, body
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) { 
    rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
    if err != nil {
        t.Fatal(err)
    }

    defer rs.Body.Close()
    body, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }

    return rs.StatusCode, rs.Header, body
}
