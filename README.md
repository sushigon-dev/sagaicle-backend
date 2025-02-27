# sagaicle-backend

- sagaicle のうちバックエンド

## レイヤードアーキテクチャ

1. ドメイン層 domain
2. リポジトリ層 repository
3. サービス層 service
4. ハンドラー層 handler
5. エントリーポイント cmd/main.go

## gen/

- ルーター／サーバー関連
  - `oas_router_gen.go`
  - `oas_server_gen.go`
  - HTTP リクエストを受け取り、対応するハンドラに振り分ける仕組みが実装されている
- ハンドラ／インターフェース関連
  - `oas_handlers_gen.go`
  - `oas_interfaces_gen.go`
    0 `oas_unimplemented_gen.go`
  - 各エンドポイントに対応する関数やインターフェースが定義されている
  - 初期状態では未実装（Not Implemented）状態になっている
  - ここにビジネスロジック層への委譲処理を実装する
- リクエスト／レスポンスのエンコーダ・デコーダ
  - `oas_request_decoders_gen.go`
  - `oas_response_encoders_gen.go`
  - HTTP のリクエストボディやパラメータ、レスポンスのシリアライズ・デシリアライズの処理が自動生成されている
- その他
  - `oas_middleware_gen.go`
  - `oas_security_gen.go`
  - `oas_validators_gen.go`
  - ミドルウェアやセキュリティ、バリデーション関連のコードも生成され、適切な位置で呼び出される設計になっています。
