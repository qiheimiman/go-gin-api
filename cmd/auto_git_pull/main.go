// webhook.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
)

// webhook 接收 GitHub 推送事件 测试
func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		// 读取请求体
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// 解析 JSON
		var data map[string]interface{}
		if err := json.Unmarshal(payload, &data); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// 获取 ref 字段，判断是否为 dev 分支
		ref, ok := data["ref"].(string)
		if !ok || ref != "refs/heads/dev" {
			w.WriteHeader(200)
			w.Write([]byte("Not dev branch, ignored."))
			return
		}

		log.Println("Received push to dev branch, deploying...")

		// 执行部署命令
		cmd := exec.Command("sh", "-c", `
			../../scripts/start.sh
		`)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Deployment failed: %v\nOutput: %s", err, output)
			http.Error(w, string(output), 500)
			return
		}

		log.Println("Deployment successful.")
		w.WriteHeader(200)
		w.Write([]byte("Dev branch updated and app restarted."))
	})

	log.Println("Webhook server listening on :83")
	log.Fatal(http.ListenAndServe(":83", nil))
}
