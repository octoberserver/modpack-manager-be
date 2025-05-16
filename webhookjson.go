package main

import "fmt"

func getWebhookJSON(name string, ver string, url string, cf string) string {
	return fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": ":loudspeaker: 最新模組包(%s) :loudspeaker:",
      "description": "名稱：[**%s**](%s)",
      "color": 14614528,
      "fields": [
        {
          "name": "**<:curseforge2:1110086664455471125>  Curseforge 版 (CF01.012.00)**",
          "value": "%s"
        }
      ],
      "footer": {
        "text": "＊2025/05/14 更新＊"
      }
    }
  ],
  "attachments": []
}`, ver, name, url, cf)
}
