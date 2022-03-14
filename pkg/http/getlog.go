package http

import (
    "net/http"

    "github.com/gorilla/mux"
)

func GetLog(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    date := vars["date"]

    logFilePath := "logs/app-"+date+".log"

    http.ServeFile(res, req, logFilePath)
}
