## テーブル定義

### users

| カラム名 | データ型 | NULL | PK | FK | UNIQUE | 説明 |
|----------|----------|------|----|----|--------|------|
| id | BIGSERIAL | × | ○ | - | ○ | ユーザーID |
| name | VARCHAR(100) | × | - | - | - | 氏名 |
| email | VARCHAR(255) | × | - | - | ○ | メールアドレス |
| password_hash | VARCHAR(255) | × | - | - | - | パスワードハッシュ |
| role | VARCHAR(20) | × | - | - | - | 権限（user/admin） |
| created_at | TIMESTAMP | × | - | - | - | 作成日時 |
| updated_at | TIMESTAMP | × | - | - | - | 更新日時 |

**インデックス**

| インデックス名 | カラム | 種別 | 用途 |
|---------------|--------|------|------|
| pk_users | id | PRIMARY KEY | 主キー |
| uk_users_email | email | UNIQUE | メールアドレス重複防止 |
| idx_users_role | role | INDEX | 権限検索 |

---

### services

| カラム名 | データ型 | NULL | PK | FK | UNIQUE | 説明 |
|----------|----------|------|----|----|--------|------|
| id | BIGSERIAL | × | ○ | - | ○ | サービスID |
| name | VARCHAR(100) | × | - | - | - | サービス名 |
| duration_minutes | INTEGER | × | - | - | - | 所要時間 |
| price | INTEGER | × | - | - | - | 料金 |
| is_active | BOOLEAN | × | - | - | - | 利用可否 |
| created_at | TIMESTAMP | × | - | - | - | 作成日時 |
| updated_at | TIMESTAMP | × | - | - | - | 更新日時 |

**インデックス**

| インデックス名 | カラム | 種別 | 用途 |
|---------------|--------|------|------|
| pk_services | id | PRIMARY KEY | 主キー |
| idx_services_active | is_active | INDEX | 有効サービス検索 |

---

### reservation_slots

| カラム名 | データ型 | NULL | PK | FK | UNIQUE | 説明 |
|----------|----------|------|----|----|--------|------|
| id | BIGSERIAL | × | ○ | - | ○ | 予約枠ID |
| start_time | TIMESTAMP | × | - | - | - | 開始日時 |
| end_time | TIMESTAMP | × | - | - | - | 終了日時 |
| capacity | INTEGER | × | - | - | - | 予約可能人数 |
| created_at | TIMESTAMP | × | - | - | - | 作成日時 |
| updated_at | TIMESTAMP | × | - | - | - | 更新日時 |

**インデックス**

| インデックス名 | カラム | 種別 | 用途 |
|---------------|--------|------|------|
| pk_reservation_slots | id | PRIMARY KEY | 主キー |
| idx_slots_start_time | start_time | INDEX | 日時検索 |
| idx_slots_time_range | start_time,end_time | INDEX | 期間検索 |
| uk_slots_time | start_time,end_time | UNIQUE | 同一予約枠防止 |

---

### reservations

| カラム名 | データ型 | NULL | PK | FK | UNIQUE | 説明 |
|----------|----------|------|----|----|--------|------|
| id | BIGSERIAL | × | ○ | - | ○ | 予約ID |
| user_id | BIGINT | × | - | users.id | - | 利用者 |
| service_id | BIGINT | × | - | services.id | - | サービス |
| slot_id | BIGINT | × | - | reservation_slots.id | - | 予約枠 |
| status | VARCHAR(20) | × | - | - | - | 予約状態 |
| note | TEXT | ○ | - | - | - | 備考 |
| reserved_at | TIMESTAMP | × | - | - | - | 予約日時 |
| cancelled_at | TIMESTAMP | ○ | - | - | - | キャンセル日時 |
| created_at | TIMESTAMP | × | - | - | - | 作成日時 |
| updated_at | TIMESTAMP | × | - | - | - | 更新日時 |

**インデックス**

| インデックス名 | カラム | 種別 | 用途 |
|---------------|--------|------|------|
| pk_reservations | id | PRIMARY KEY | 主キー |
| idx_reservations_user | user_id | INDEX | 利用者検索 |
| idx_reservations_slot | slot_id | INDEX | 予約枠検索 |
| idx_reservations_status | status | INDEX | ステータス検索 |
| idx_reservations_reserved_at | reserved_at | INDEX | 日次集計 |
| idx_reservations_slot_status | slot_id,status | INDEX | 空き枠検索・二重予約判定 |

---

### operation_logs

| カラム名 | データ型 | NULL | PK | FK | UNIQUE | 説明 |
|----------|----------|------|----|----|--------|------|
| id | BIGSERIAL | × | ○ | - | ○ | ログID |
| user_id | BIGINT | × | - | users.id | - | 操作ユーザー |
| reservation_id | BIGINT | ○ | - | reservations.id | - | 対象予約 |
| action | VARCHAR(50) | × | - | - | - | 操作内容 |
| detail | TEXT | ○ | - | - | - | 詳細 |
| created_at | TIMESTAMP | × | - | - | - | 作成日時 |

**インデックス**

| インデックス名 | カラム | 種別 | 用途 |
|---------------|--------|------|------|
| pk_operation_logs | id | PRIMARY KEY | 主キー |
| idx_logs_user | user_id | INDEX | 操作者検索 |
| idx_logs_reservation | reservation_id | INDEX | 予約検索 |
| idx_logs_created_at | created_at | INDEX | 時系列検索 |

---

## 外部キー一覧

| テーブル | カラム | 参照テーブル | 参照カラム | ON DELETE |
|----------|--------|-------------|------------|-----------|
| reservations | user_id | users | id | RESTRICT |
| reservations | service_id | services | id | RESTRICT |
| reservations | slot_id | reservation_slots | id | RESTRICT |
| operation_logs | user_id | users | id | RESTRICT |
| operation_logs | reservation_id | reservations | id | SET NULL |