package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type PasswordConfig struct {
	Length         int  `json:"length"`
	UseUppercase   bool `json:"useUppercase"`
	UseLowercase   bool `json:"useLowercase"`
	UseNumbers     bool `json:"useNumbers"`
	UseSpecialChar bool `json:"useSpecialChar"`
}

const (
	uppercaseChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars   = "abcdefghijklmnopqrstuvwxyz"
	numberChars      = "0123456789"
	specialCharChars = "!@#$%^&*()_-+=[]{}|:;<>,.?/~"
)

func generatePassword(config PasswordConfig) string {
	var charsToUse string

	if config.UseUppercase {
		charsToUse += uppercaseChars
	}

	if config.UseLowercase {
		charsToUse += lowercaseChars
	}

	if config.UseNumbers {
		charsToUse += numberChars
	}

	if config.UseSpecialChar {
		charsToUse += specialCharChars
	}

	rand.Seed(time.Now().UnixNano())
	password := make([]byte, config.Length)

	for i := range password {
		password[i] = charsToUse[rand.Intn(len(charsToUse))]
	}

	return string(password)
}

func getQueryParam(r *http.Request, param string) (string, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return "", httpError("O parâmetro '"+param+"' não foi fornecido", http.StatusBadRequest)
	}
	return value, nil
}

func httpError(message string, statusCode int) error {
	return &httpErrorWrapper{message, statusCode}
}

type httpErrorWrapper struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *httpErrorWrapper) Error() string {
	return e.Message
}

func generatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	lengthParam, err := getQueryParam(r, "length")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	useUppercaseParam, err := getQueryParam(r, "useUppercase")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	useLowercaseParam, err := getQueryParam(r, "useLowercase")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	useNumbersParam, err := getQueryParam(r, "useNumbers")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	useSpecialCharParam, err := getQueryParam(r, "useSpecialChar")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	length, err := strconv.Atoi(lengthParam)

	if err != nil || length <= 0 {
		http.Error(w, "O parâmetro 'length' deve ser um número inteiro positivo", http.StatusBadRequest)
		return
	}

	useUppercase, err := strconv.ParseBool(useUppercaseParam)

	if err != nil {
		http.Error(w, "O parâmetro 'useUppercase' deve ser um valor booleano ('true' ou 'false')", http.StatusBadRequest)
		return
	}

	useLowercase, err := strconv.ParseBool(useLowercaseParam)

	if err != nil {
		http.Error(w, "O parâmetro 'useLowercase' deve ser um valor booleano ('true' ou 'false')", http.StatusBadRequest)
		return
	}

	useNumbers, err := strconv.ParseBool(useNumbersParam)

	if err != nil {
		http.Error(w, "O parâmetro 'useNumbers' deve ser um valor booleano ('true' ou 'false')", http.StatusBadRequest)
		return
	}

	useSpecialChar, err := strconv.ParseBool(useSpecialCharParam)

	if err != nil {
		http.Error(w, "O parâmetro 'useSpecialChar' deve ser um valor booleano ('true' ou 'false')", http.StatusBadRequest)
		return
	}

	config := PasswordConfig{
		Length:         length,
		UseUppercase:   useUppercase,
		UseLowercase:   useLowercase,
		UseNumbers:     useNumbers,
		UseSpecialChar: useSpecialChar,
	}

	password := generatePassword(config)

	response := map[string]string{ "password": password }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/generate-password", generatePasswordHandler)

	port := ":8080"
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}