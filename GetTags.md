# sagaicle-backend API(GetTags)設計書

## 検索可能なタグを全て取得
- GetTags
- エンドポイント: `/api/tags`
- HTTP メソッド: `GET`

###成功

- HTTP/1.1 200 OK
  リクエスト例
  ```bash
  curl http://localhost:8080/api/tags
  ```

  レスポンス例
  ```json
  {
  "tags": ["温泉", "ファミリー", "エンジョイ”]
  }
  ```

### 失敗
- HTTP/1.1 500 Internal Server Error
 - 取得可能なタグが見つからない
  リクエスト例
  ```bash
  curl http://localhost:8080/api/tags
  ```

  レスポンス例
  ```json
  {
  "error": "tag not found"
  }
  ```