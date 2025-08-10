package mw

import (
	"fmt"
	"net/http"
)

func CheckTaskIndex(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		fmt.Println("MW handler")

		
	})
}