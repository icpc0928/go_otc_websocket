package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"os"
)

func main() {

	initConfig()
	_ = context.Background()

	//---------------
	m := melody.New() //new方法中可以配置一些參數

	m.HandleConnect(func(s *melody.Session) {
		fmt.Println("建立連線")
		fmt.Println("session.request: ", s.Request)
		fmt.Println("session.Keys: ", s.Keys)
		//fmt.Println("LocalAddr: ", s.LocalAddr(), "RemoteAddr: ", s.Remo)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		fmt.Println("連線斷開")
	})
	m.HandleError(func(s *melody.Session, err error) {
		fmt.Println("出現錯誤")
		fmt.Println(err)
	})
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println("收到消息: ", string(msg))
		s.Write([]byte("return msg")) //向客戶發送消息
	})

	http.HandleFunc("/ws-m", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r) //http訪問 /ws-m 時轉交給melody處理 實現ws功能
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("index incoming")
		http.ServeFile(w, r, "index.html") //訪問 / 給靜態檔
	})
	http.ListenAndServe(viper.GetString("service"), nil) //啟動服務氣

}

func initConfig() {
	//box := main - otc
	var env = os.Getenv("env")
	fmt.Println("env: " + env)
	// export CONFIG_PATH = prod
	//func TestGetenv(t *testing.T){
	//        path := os.Getenv("CONFIG_PATH")
	//        fmt.Println(path)
	//}
	if env != "" {
		env = "." + env
	}
	viper.SetConfigName("config" + env) // name of config file (without extension)
	viper.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
