package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HelloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func TestWeb(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	recorder := httptest.NewRecorder()
	HelloHandle(recorder, req)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)
	fmt.Println(bodyString)
}

func ParamHandle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		fmt.Fprintf(w, "Hello %s!", name)
	}
	email := r.URL.Query().Get("email")
	if email != "" {
		fmt.Fprintf(w, ", Email: %s", email)
	}
}

func TestParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080?name=eko&email=eko@gmail.com", nil)
	rec := httptest.NewRecorder()

	ParamHandle(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)

	assert.Equal(t, "Hello eko!, Email: eko@gmail.com", bodyString)
}

func TestParamValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080?name=eko&name=budi", nil)
	rec := httptest.NewRecorder()

	ParamHandleValue(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)

	assert.Equal(t, "eko budi", bodyString)

}

func ParamHandleValue(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	names := query["name"]
	fmt.Fprint(w, strings.Join(names, " "))
}

func ParamHandleFormPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	name := r.PostForm.Get("name")
	fmt.Fprint(w, name)
}

func TestFormPost(t *testing.T) {
	// values := url.Values{
	// 	"name": {"John"},
	// }
	values := "name=John"

	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(values))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	ParamHandleFormPost(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)

	assert.Equal(t, "John", bodyString)
}

func SetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := new(http.Cookie)
	cookie.Name = "x-name"
	cookie.Value = r.URL.Query().Get("name")
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "Cookie set")
}

func TestSetCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/?name=John", nil)
	recorder := httptest.NewRecorder()

	SetCookie(recorder, req)

	response := recorder.Result()
	cookie := response.Cookies()[0]
	assert.Equal(t, "John", cookie.Value)
}

func GetCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("x-name")
	if err != nil {
		fmt.Fprintln(w, "no cookie")
		return
	}
	fmt.Fprintln(w, c.Value)
}

func TestGetCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "x-name",
		Value: "John",
	})

	rec := httptest.NewRecorder()
	GetCookie(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)

	assert.Equal(t, "John\n", bodyString)
}
