# シーケンス図

## 1. ログイン

```mermaid
sequenceDiagram
    actor User as 利用者
    participant Browser as Webブラウザ
    participant Handler as AuthHandler
    participant Service as AuthService
    participant Repo as UserRepository
    participant DB as PostgreSQL
    participant JWT as JWTService

    User->>Browser: メールアドレス・パスワード入力
    Browser->>Handler: POST /api/v1/auth/login
    Handler->>Service: Login(email, password)
    Service->>Repo: FindByEmail(email)
    Repo->>DB: SELECT * FROM users WHERE email = ?
    DB-->>Repo: ユーザー情報
    Repo-->>Service: User
    Service->>Service: パスワード検証
    Service->>JWT: GenerateToken(userId, role)
    JWT-->>Service: accessToken
    Service-->>Handler: LoginResult
    Handler-->>Browser: 200 OK + JWT
    Browser-->>User: マイページ表示
```

---

## 2. 空き枠検索

```mermaid
sequenceDiagram
    actor User as 利用者
    participant Browser as Webブラウザ
    participant Handler as ReservationSlotHandler
    participant Service as ReservationSlotService
    participant ServiceRepo as ServiceRepository
    participant SlotRepo as ReservationSlotRepository
    participant ReservationRepo as ReservationRepository
    participant DB as PostgreSQL

    User->>Browser: サービス・日付を選択
    Browser->>Handler: GET /api/v1/reservation-slots?date=...&serviceId=...
    Handler->>Service: SearchAvailableSlots(date, serviceId)
    Service->>ServiceRepo: FindByID(serviceId)
    ServiceRepo->>DB: SELECT * FROM services WHERE id = ?
    DB-->>ServiceRepo: サービス情報
    ServiceRepo-->>Service: Service

    Service->>SlotRepo: FindByDate(date)
    SlotRepo->>DB: SELECT * FROM reservation_slots WHERE start_time BETWEEN ...
    DB-->>SlotRepo: 予約枠一覧
    SlotRepo-->>Service: Slots

    Service->>ReservationRepo: CountReservedBySlotIds(slotIds)
    ReservationRepo->>DB: SELECT slot_id, COUNT(*) FROM reservations WHERE status = 'reserved' GROUP BY slot_id
    DB-->>ReservationRepo: 予約済み件数
    ReservationRepo-->>Service: ReservedCounts

    Service->>Service: 空き枠判定
    Service-->>Handler: 空き枠一覧
    Handler-->>Browser: 200 OK
    Browser-->>User: 空き枠一覧表示
```

---

## 3. 予約登録

```mermaid
sequenceDiagram
    actor User as 利用者
    participant Browser as Webブラウザ
    participant Handler as ReservationHandler
    participant Service as ReservationService
    participant Tx as TransactionManager
    participant SlotRepo as ReservationSlotRepository
    participant ReservationRepo as ReservationRepository
    participant LogRepo as OperationLogRepository
    participant DB as PostgreSQL

    User->>Browser: 予約内容を確定
    Browser->>Handler: POST /api/v1/reservations
    Handler->>Service: CreateReservation(userId, request)

    Service->>Tx: Begin
    Tx->>DB: BEGIN

    Service->>SlotRepo: FindByIDForUpdate(slotId)
    SlotRepo->>DB: SELECT * FROM reservation_slots WHERE id = ? FOR UPDATE
    DB-->>SlotRepo: 予約枠情報
    SlotRepo-->>Service: Slot

    Service->>ReservationRepo: ExistsActiveBySlotID(slotId)
    ReservationRepo->>DB: SELECT EXISTS ... WHERE slot_id = ? AND status = 'reserved'
    DB-->>ReservationRepo: true / false
    ReservationRepo-->>Service: 予約済み判定

    alt 既に予約済み
        Service->>Tx: Rollback
        Tx->>DB: ROLLBACK
        Service-->>Handler: RESERVATION_CONFLICT
        Handler-->>Browser: 409 Conflict
        Browser-->>User: 予約済みエラー表示
    else 予約可能
        Service->>ReservationRepo: Create(reservation)
        ReservationRepo->>DB: INSERT INTO reservations ...
        DB-->>ReservationRepo: reservationId

        Service->>LogRepo: Create(operationLog)
        LogRepo->>DB: INSERT INTO operation_logs ...
        DB-->>LogRepo: OK

        Service->>Tx: Commit
        Tx->>DB: COMMIT

        Service-->>Handler: ReservationResult
        Handler-->>Browser: 201 Created
        Browser-->>User: 予約完了画面表示
    end
```

---

## 4. 予約変更

```mermaid
sequenceDiagram
    actor User as 利用者
    participant Browser as Webブラウザ
    participant Handler as ReservationHandler
    participant Service as ReservationService
    participant Tx as TransactionManager
    participant ReservationRepo as ReservationRepository
    participant SlotRepo as ReservationSlotRepository
    participant LogRepo as OperationLogRepository
    participant DB as PostgreSQL

    User->>Browser: 変更内容を確定
    Browser->>Handler: PUT /api/v1/reservations/{reservationId}
    Handler->>Service: UpdateReservation(userId, reservationId, request)

    Service->>Tx: Begin
    Tx->>DB: BEGIN

    Service->>ReservationRepo: FindByID(reservationId)
    ReservationRepo->>DB: SELECT * FROM reservations WHERE id = ?
    DB-->>ReservationRepo: 予約情報
    ReservationRepo-->>Service: Reservation

    Service->>Service: 予約所有者チェック
    Service->>Service: 予約状態チェック

    Service->>SlotRepo: FindByIDForUpdate(newSlotId)
    SlotRepo->>DB: SELECT * FROM reservation_slots WHERE id = ? FOR UPDATE
    DB-->>SlotRepo: 新予約枠情報
    SlotRepo-->>Service: Slot

    Service->>ReservationRepo: ExistsActiveBySlotID(newSlotId)
    ReservationRepo->>DB: SELECT EXISTS ... WHERE slot_id = ? AND status = 'reserved'
    DB-->>ReservationRepo: true / false
    ReservationRepo-->>Service: 予約済み判定

    alt 新予約枠が予約済み
        Service->>Tx: Rollback
        Tx->>DB: ROLLBACK
        Service-->>Handler: RESERVATION_CONFLICT
        Handler-->>Browser: 409 Conflict
        Browser-->>User: 予約済みエラー表示
    else 変更可能
        Service->>ReservationRepo: Update(reservation)
        ReservationRepo->>DB: UPDATE reservations SET service_id=?, slot_id=?, note=?, updated_at=?
        DB-->>ReservationRepo: OK

        Service->>LogRepo: Create(operationLog)
        LogRepo->>DB: INSERT INTO operation_logs ...
        DB-->>LogRepo: OK

        Service->>Tx: Commit
        Tx->>DB: COMMIT

        Service-->>Handler: UpdateResult
        Handler-->>Browser: 200 OK
        Browser-->>User: 変更完了表示
    end
```

---

## 5. 予約キャンセル

```mermaid
sequenceDiagram
    actor User as 利用者
    participant Browser as Webブラウザ
    participant Handler as ReservationHandler
    participant Service as ReservationService
    participant Tx as TransactionManager
    participant ReservationRepo as ReservationRepository
    participant LogRepo as OperationLogRepository
    participant DB as PostgreSQL

    User->>Browser: キャンセル実行
    Browser->>Handler: DELETE /api/v1/reservations/{reservationId}
    Handler->>Service: CancelReservation(userId, reservationId)

    Service->>Tx: Begin
    Tx->>DB: BEGIN

    Service->>ReservationRepo: FindByID(reservationId)
    ReservationRepo->>DB: SELECT * FROM reservations WHERE id = ?
    DB-->>ReservationRepo: 予約情報
    ReservationRepo-->>Service: Reservation

    Service->>Service: 予約所有者チェック
    Service->>Service: キャンセル可能状態チェック

    alt 既にキャンセル済み
        Service->>Tx: Rollback
        Tx->>DB: ROLLBACK
        Service-->>Handler: ALREADY_CANCELLED
        Handler-->>Browser: 409 Conflict
        Browser-->>User: キャンセル済みエラー表示
    else キャンセル可能
        Service->>ReservationRepo: UpdateStatus(cancelled)
        ReservationRepo->>DB: UPDATE reservations SET status='cancelled', cancelled_at=?, updated_at=?
        DB-->>ReservationRepo: OK

        Service->>LogRepo: Create(operationLog)
        LogRepo->>DB: INSERT INTO operation_logs ...
        DB-->>LogRepo: OK

        Service->>Tx: Commit
        Tx->>DB: COMMIT

        Service-->>Handler: CancelResult
        Handler-->>Browser: 200 OK
        Browser-->>User: キャンセル完了画面表示
    end
```

---

## 6. 管理者：予約一覧検索

```mermaid
sequenceDiagram
    actor Admin as 管理者
    participant Browser as 管理画面
    participant Handler as AdminReservationHandler
    participant Service as AdminReservationService
    participant Repo as ReservationRepository
    participant DB as PostgreSQL

    Admin->>Browser: 検索条件を入力
    Browser->>Handler: GET /api/v1/admin/reservations?date=...&status=...
    Handler->>Service: SearchReservations(condition)
    Service->>Service: 管理者権限チェック
    Service->>Repo: Search(condition)
    Repo->>DB: SELECT reservations ... WHERE ...
    DB-->>Repo: 予約一覧
    Repo-->>Service: Reservations
    Service-->>Handler: ReservationList
    Handler-->>Browser: 200 OK
    Browser-->>Admin: 予約一覧表示
```

---

## 7. 管理者：予約ステータス変更

```mermaid
sequenceDiagram
    actor Admin as 管理者
    participant Browser as 管理画面
    participant Handler as AdminReservationHandler
    participant Service as AdminReservationService
    participant Tx as TransactionManager
    participant ReservationRepo as ReservationRepository
    participant LogRepo as OperationLogRepository
    participant DB as PostgreSQL

    Admin->>Browser: ステータス変更
    Browser->>Handler: PATCH /api/v1/admin/reservations/{reservationId}/status
    Handler->>Service: UpdateReservationStatus(adminId, reservationId, status)

    Service->>Service: 管理者権限チェック

    Service->>Tx: Begin
    Tx->>DB: BEGIN

    Service->>ReservationRepo: FindByID(reservationId)
    ReservationRepo->>DB: SELECT * FROM reservations WHERE id = ?
    DB-->>ReservationRepo: 予約情報
    ReservationRepo-->>Service: Reservation

    Service->>Service: ステータス遷移チェック

    Service->>ReservationRepo: UpdateStatus(status)
    ReservationRepo->>DB: UPDATE reservations SET status=?, updated_at=?
    DB-->>ReservationRepo: OK

    Service->>LogRepo: Create(operationLog)
    LogRepo->>DB: INSERT INTO operation_logs ...
    DB-->>LogRepo: OK

    Service->>Tx: Commit
    Tx->>DB: COMMIT

    Service-->>Handler: UpdateStatusResult
    Handler-->>Browser: 200 OK
    Browser-->>Admin: 更新結果表示
```

---

## 8. 管理者：サービス登録・更新

```mermaid
sequenceDiagram
    actor Admin as 管理者
    participant Browser as 管理画面
    participant Handler as ServiceAdminHandler
    participant Service as ServiceAdminService
    participant Repo as ServiceRepository
    participant DB as PostgreSQL

    Admin->>Browser: サービス情報を入力
    Browser->>Handler: POST /api/v1/admin/services
    Handler->>Service: CreateService(adminId, request)

    Service->>Service: 管理者権限チェック
    Service->>Service: 入力値検証
    Service->>Repo: Create(service)
    Repo->>DB: INSERT INTO services ...
    DB-->>Repo: serviceId
    Repo-->>Service: Service

    Service-->>Handler: ServiceResult
    Handler-->>Browser: 201 Created
    Browser-->>Admin: サービス一覧を更新表示
```

---

## 9. 管理者：予約枠登録

```mermaid
sequenceDiagram
    actor Admin as 管理者
    participant Browser as 管理画面
    participant Handler as ReservationSlotAdminHandler
    participant Service as ReservationSlotAdminService
    participant Repo as ReservationSlotRepository
    participant DB as PostgreSQL

    Admin->>Browser: 日時・定員を入力
    Browser->>Handler: POST /api/v1/admin/reservation-slots
    Handler->>Service: CreateReservationSlot(adminId, request)

    Service->>Service: 管理者権限チェック
    Service->>Service: 日時・定員チェック
    Service->>Repo: ExistsSameTimeRange(startTime, endTime)
    Repo->>DB: SELECT EXISTS ... WHERE start_time=? AND end_time=?
    DB-->>Repo: true / false

    alt 同一時間帯の予約枠あり
        Service-->>Handler: VALIDATION_ERROR
        Handler-->>Browser: 400 Bad Request
        Browser-->>Admin: 入力エラー表示
    else 登録可能
        Service->>Repo: Create(slot)
        Repo->>DB: INSERT INTO reservation_slots ...
        DB-->>Repo: slotId
        Repo-->>Service: ReservationSlot

        Service-->>Handler: SlotResult
        Handler-->>Browser: 201 Created
        Browser-->>Admin: 予約枠一覧を更新表示
    end
```

---

## 10. 日次集計バッチ

```mermaid
sequenceDiagram
    participant Scheduler as Scheduler
    participant Batch as DailyReportBatch
    participant ReservationRepo as ReservationRepository
    participant ReportRepo as ReportRepository
    participant FileStorage as FileStorage
    participant DB as PostgreSQL

    Scheduler->>Batch: 日次集計バッチ起動
    Batch->>ReservationRepo: FindByDate(targetDate)
    ReservationRepo->>DB: SELECT * FROM reservations WHERE reserved_at BETWEEN ...
    DB-->>ReservationRepo: 対象予約一覧
    ReservationRepo-->>Batch: Reservations

    Batch->>Batch: 予約件数・キャンセル数・売上見込み集計
    Batch->>ReportRepo: SaveDailyReport(report)
    ReportRepo->>DB: INSERT INTO daily_reports ...
    DB-->>ReportRepo: OK

    Batch->>Batch: CSV生成
    Batch->>FileStorage: Save(csvFile)
    FileStorage-->>Batch: 保存完了

    Batch-->>Scheduler: 正常終了
```