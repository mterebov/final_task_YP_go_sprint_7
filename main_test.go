package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, item := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=moscow", item.count), nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		if item.count == 0 {
			assert.Equal(t, item.want, len(strings.TrimSpace(response.Body.String())))
			continue
		}
		assert.Equal(t, item.want, len(strings.Split(strings.TrimSpace(response.Body.String()), ",")))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search string
		want int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	
	for _, test_req := range requests {
		var responseCafeCount int
		var counterForCheckSubstr int

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", test_req.search), nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")

		for _, str := range cafes {
			if strings.Contains(strings.ToUpper(str), strings.ToUpper(test_req.search)) {
				counterForCheckSubstr ++
			}
		}

		if response.Body.String() == "" {
			responseCafeCount = 0
		} else {
			responseCafeCount = len(cafes)
		}
		
		// Проверяем что все выданное сервером содержит искомую подстроку
		assert.Equal(t, counterForCheckSubstr, responseCafeCount)
		// Проверяем количество выданного сервером и ожидаемое
		assert.Equal(t, test_req.want, responseCafeCount)
	}
}
