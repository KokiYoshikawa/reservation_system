# API設計

## 1. API概要

本システムは、予約管理システムのバックエンドAPIとして、Goで実装したREST APIを提供する。  
利用者は予約枠の検索、予約登録、予約変更、キャンセルを行う。  
管理者は予約枠、サービス、利用者、予約状況を管理する。

---

## 2. 共通仕様

| 項目 | 内容 |
|---|---|
| Base URL | `/api/v1` |
| 通信方式 | HTTPS |
| データ形式 | JSON |
| 認証方式 | JWT Bearer Token |
| 文字コード | UTF-8 |
| 日時形式 | ISO 8601 |
| タイムゾーン | Asia/Tokyo |

### 認証ヘッダー

```http
Authorization: Bearer {access_token}
```

---

## 3. レスポンス共通形式

### 正常レスポンス

```json
{
  "data": {}
}
```

### エラーレスポンス

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力内容に誤りがあります。",
    "details": [
      {
        "field": "email",
        "message": "メールアドレスの形式が不正です。"
      }
    ]
  }
}
```

---

## 4. HTTPステータス

| ステータス | 内容 |
|---|---|
| 200 | 正常終了 |
| 201 | 登録成功 |
| 204 | 削除成功 |
| 400 | リクエスト不正 |
| 401 | 未認証 |
| 403 | 権限不足 |
| 404 | 対象データなし |
| 409 | 業務エラー・競合 |
| 500 | サーバーエラー |

---

## 5. API一覧

## 5.1 認証API

| No | メソッド | エンドポイント | 説明 | 認証 |
|---|---|---|---|---|
| A-001 | POST | `/auth/register` | 利用者登録 | 不要 |
| A-002 | POST | `/auth/login` | ログイン | 不要 |
| A-003 | POST | `/auth/logout` | ログアウト | 必要 |
| A-004 | GET | `/auth/me` | ログインユーザー情報取得 | 必要 |

---

## 5.2 利用者API

| No | メソッド | エンドポイント | 説明 | 認証 |
|---|---|---|---|---|
| U-001 | GET | `/services` | サービス一覧取得 | 不要 |
| U-002 | GET | `/reservation-slots` | 空き枠検索 | 不要 |
| U-003 | POST | `/reservations` | 予約登録 | 必要 |
| U-004 | GET | `/reservations` | 自分の予約一覧取得 | 必要 |
| U-005 | GET | `/reservations/{reservationId}` | 予約詳細取得 | 必要 |
| U-006 | PUT | `/reservations/{reservationId}` | 予約変更 | 必要 |
| U-007 | DELETE | `/reservations/{reservationId}` | 予約キャンセル | 必要 |

---

## 5.3 管理者API

| No | メソッド | エンドポイント | 説明 | 認証 |
|---|---|---|---|---|
| M-001 | GET | `/admin/reservations` | 全予約一覧取得 | 管理者 |
| M-002 | GET | `/admin/reservations/{reservationId}` | 予約詳細取得 | 管理者 |
| M-003 | PATCH | `/admin/reservations/{reservationId}/status` | 予約ステータス変更 | 管理者 |
| M-004 | GET | `/admin/users` | 利用者一覧取得 | 管理者 |
| M-005 | GET | `/admin/users/{userId}` | 利用者詳細取得 | 管理者 |
| M-006 | POST | `/admin/services` | サービス登録 | 管理者 |
| M-007 | PUT | `/admin/services/{serviceId}` | サービス更新 | 管理者 |
| M-008 | DELETE | `/admin/services/{serviceId}` | サービス削除 | 管理者 |
| M-009 | POST | `/admin/reservation-slots` | 予約枠登録 | 管理者 |
| M-010 | PUT | `/admin/reservation-slots/{slotId}` | 予約枠更新 | 管理者 |
| M-011 | DELETE | `/admin/reservation-slots/{slotId}` | 予約枠削除 | 管理者 |

---

## 5.4 集計・CSV API

| No | メソッド | エンドポイント | 説明 | 認証 |
|---|---|---|---|---|
| R-001 | GET | `/admin/reports/daily` | 日次予約集計取得 | 管理者 |
| R-002 | GET | `/admin/reports/daily/csv` | 日次予約集計CSV出力 | 管理者 |
| R-003 | GET | `/admin/reservations/csv` | 予約一覧CSV出力 | 管理者 |

---

# 6. API詳細

## A-002 ログイン

### Endpoint

```http
POST /api/v1/auth/login
```

### Request

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Response

```json
{
  "data": {
    "accessToken": "jwt-token",
    "user": {
      "id": 1,
      "name": "山田太郎",
      "email": "user@example.com",
      "role": "user"
    }
  }
}
```

---

## U-002 空き枠検索

### Endpoint

```http
GET /api/v1/reservation-slots?date=2026-07-01&serviceId=1
```

### Query Parameters

| パラメータ | 型 | 必須 | 説明 |
|---|---|---|---|
| date | string | ○ | 検索対象日 |
| serviceId | integer | ○ | サービスID |

### Response

```json
{
  "data": {
    "date": "2026-07-01",
    "slots": [
      {
        "slotId": 101,
        "startTime": "2026-07-01T10:00:00+09:00",
        "endTime": "2026-07-01T11:00:00+09:00",
        "capacity": 1,
        "reservedCount": 0,
        "available": true
      }
    ]
  }
}
```

---

## U-003 予約登録

### Endpoint

```http
POST /api/v1/reservations
```

### Request

```json
{
  "serviceId": 1,
  "slotId": 101,
  "note": "初回利用です"
}
```

### Response

```json
{
  "data": {
    "reservationId": 1001,
    "status": "reserved",
    "service": {
      "id": 1,
      "name": "カウンセリング"
    },
    "startTime": "2026-07-01T10:00:00+09:00",
    "endTime": "2026-07-01T11:00:00+09:00",
    "reservedAt": "2026-06-28T22:30:00+09:00"
  }
}
```

### 業務ルール

- 存在しない予約枠には予約できない
- 無効なサービスには予約できない
- 過去日時の予約枠には予約できない
- 同一予約枠に対する二重予約は不可
- 予約登録時に操作ログを登録する

---

## U-006 予約変更

### Endpoint

```http
PUT /api/v1/reservations/{reservationId}
```

### Path Parameters

| パラメータ | 型 | 必須 | 説明 |
|---|---|---|---|
| reservationId | integer | ○ | 予約ID |

### Request

```json
{
  "serviceId": 2,
  "slotId": 105,
  "note": "時間変更希望"
}
```

### Response

```json
{
  "data": {
    "reservationId": 1001,
    "status": "reserved",
    "service": {
      "id": 2,
      "name": "整体"
    },
    "startTime": "2026-07-01T15:00:00+09:00",
    "endTime": "2026-07-01T16:00:00+09:00",
    "updatedAt": "2026-06-28T22:45:00+09:00"
  }
}
```

---

## U-007 予約キャンセル

### Endpoint

```http
DELETE /api/v1/reservations/{reservationId}
```

### Response

```json
{
  "data": {
    "reservationId": 1001,
    "status": "cancelled",
    "cancelledAt": "2026-06-28T22:50:00+09:00"
  }
}
```

---

## M-003 予約ステータス変更

### Endpoint

```http
PATCH /api/v1/admin/reservations/{reservationId}/status
```

### Request

```json
{
  "status": "completed"
}
```

### Response

```json
{
  "data": {
    "reservationId": 1001,
    "status": "completed",
    "updatedAt": "2026-07-01T11:10:00+09:00"
  }
}
```

---

## M-009 予約枠登録

### Endpoint

```http
POST /api/v1/admin/reservation-slots
```

### Request

```json
{
  "startTime": "2026-07-01T10:00:00+09:00",
  "endTime": "2026-07-01T11:00:00+09:00",
  "capacity": 1
}
```

### Response

```json
{
  "data": {
    "slotId": 101,
    "startTime": "2026-07-01T10:00:00+09:00",
    "endTime": "2026-07-01T11:00:00+09:00",
    "capacity": 1
  }
}
```

---

# 7. ステータス定義

## 予約ステータス

| 値 | 内容 |
|---|---|
| reserved | 予約済み |
| cancelled | キャンセル済み |
| completed | 対応完了 |
| no_show | 無断キャンセル |

---

## ユーザー権限

| 値 | 内容 |
|---|---|
| user | 一般利用者 |
| admin | 管理者 |

---

# 8. エラーコード定義

| エラーコード | HTTPステータス | 内容 |
|---|---|---|
| VALIDATION_ERROR | 400 | 入力値不正 |
| UNAUTHORIZED | 401 | 未認証 |
| FORBIDDEN | 403 | 権限不足 |
| NOT_FOUND | 404 | 対象データなし |
| RESERVATION_CONFLICT | 409 | 予約枠の競合 |
| ALREADY_CANCELLED | 409 | 既にキャンセル済み |
| INTERNAL_SERVER_ERROR | 500 | サーバー内部エラー |

---

# 9. 設計上の補足

## REST API設計方針

本システムでは、リソース単位でエンドポイントを設計している。

例：

- `/reservations`
- `/reservation-slots`
- `/services`
- `/users`

登録・取得・更新・削除はHTTPメソッドで表現する。

| 操作 | HTTPメソッド |
|---|---|
| 登録 | POST |
| 取得 | GET |
| 更新 | PUT / PATCH |
| 削除 | DELETE |

---

## 認証・認可方針

JWTを利用してログイン状態を管理する。  
管理者APIは、JWT内のroleが`admin`の場合のみ実行可能とする。

---

## 二重予約防止方針

予約登録時は、Service層でトランザクションを開始し、対象の予約枠をロックする。  
その上で、同一slot_idに対して`reserved`状態の予約が存在しないことを確認する。

また、DB側にも部分ユニークインデックスを設定し、アプリケーション側のチェック漏れがあっても二重予約を防止する。

```sql
CREATE UNIQUE INDEX reservations_active_slot_uq
ON reservations (slot_id)
WHERE status = 'reserved';
```