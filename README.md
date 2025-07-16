# Discord 機器人模板

go 1.24.2

## 使用

1. 克隆代碼 `git clone https://github.com/Wuchieh/discord-bot-template.git`
2. 設置配置文件 `config.yaml` (可以將範例文件 `config.example.yaml` 複製成 `config.yaml`)
3. 編譯程序 `go build . -o discord-bot`
4. 設定執行權限 `chmod +x discord-bot` (這一步不一定需要執行)
5. 執行 `./discord-bot`

```bash  
    git clone https://github.com/Wuchieh/discord-bot-template.git discord-bot
    cd discord-bot
    mv config.example.yaml config.yaml
    go build . -o discord-bot
    chmod +x discord-bot
    ./discord-bot
```

## 設定
```yaml
bot_token: [機器人Token]

db:
    file: [資料庫檔案名稱] # 結尾必須是 .db
    log_level: [日誌等級] # info, warn, error, silent (預設為 warn)

reaction_role: # 反應角色功能
    enabled: [是否啟用] # true, false (預設為 false)
    guild_id: # 伺服器ID 可以填寫多個
        - [伺服器ID-1]
        - [伺服器ID-2]
        - [伺服器ID-3]
```