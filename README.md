# HTTP-сервер магазина одежды

Go 1.24 + SQLite (modernc.org/sqlite, без CGO).

## Сетап

1. Установить Go 1.24+.
2. Скопировать `.env.example` в `.env`.
3. Установить зависимости:

go mod tidy


4. Запустить сервер:

go run cmd/server/main.go


Сервер на http://localhost:8080.

## Команды для показа (PowerShell)

Использовать `curl.exe` — в PowerShell `curl` это алиас на Invoke-WebRequest.

1. Список ролей (создаются автоматически при старте):

curl.exe http://localhost:8080/role


2. Добавление товаров (админ):

curl.exe -X POST http://localhost:8080/product -H "Content-Type: application/json" -d '{"name":"Футболка","type":"Одежда","gender":"М","size":"L","color":"Черный","brand":"Test","season":"Лето","price":1500}'

curl.exe -X POST http://localhost:8080/product -H "Content-Type: application/json" -d '{"name":"Джинсы","type":"Одежда","gender":"М","size":"M","color":"Синий","brand":"Test","season":"Демисезон","price":3500}'

curl.exe -X POST http://localhost:8080/product -H "Content-Type: application/json" -d '{"name":"Платье","type":"Одежда","gender":"Ж","size":"S","color":"Белый","brand":"Other","season":"Лето","price":2500}'


3. Каталог с фильтрацией и сортировкой (гость):

curl.exe "http://localhost:8080/product?gender=М"

curl.exe "http://localhost:8080/product?sort_by=price&order=desc"

curl.exe "http://localhost:8080/product?brand=Test&limit=2&offset=0"


4. Детальная карточка товара:

curl.exe http://localhost:8080/product/<PRODUCT_ID>


5. Регистрация клиента:

curl.exe -X POST http://localhost:8080/user -H "Content-Type: application/json" -d '{"name":"Иван","email":"ivan@test.com","phone":"+79991234567","role_id":"22222222-2222-2222-2222-222222222222"}'


6. Добавление в корзину:

curl.exe -X POST http://localhost:8080/cart -H "Content-Type: application/json" -d '{"user_id":"<USER_ID>","product_id":"<PRODUCT_ID>","count":2}'


7. Оформление заказа:

curl.exe -X POST http://localhost:8080/order -H "Content-Type: application/json" -d '{"user_id":"<USER_ID>","status":"Принят"}'


8. Добавление элементов заказа:

curl.exe -X POST http://localhost:8080/order_item -H "Content-Type: application/json" -d '{"order_id":"<ORDER_ID>","product_id":"<PRODUCT_ID>","price":1500}'


9. Менеджер меняет статус заказа:

curl.exe "http://localhost:8080/order?status=Принят"

curl.exe -X PUT http://localhost:8080/order/<ORDER_ID> -H "Content-Type: application/json" -d '{"status":"В работе"}'

curl.exe -X PUT http://localhost:8080/order/<ORDER_ID> -H "Content-Type: application/json" -d '{"status":"Собран"}'

curl.exe -X PUT http://localhost:8080/order/<ORDER_ID> -H "Content-Type: application/json" -d '{"status":"Готов к выдаче"}'


10. Загрузка и получение фото:

curl.exe -X POST http://localhost:8080/photo -F "product_id=<PRODUCT_ID>" -F "photo=@C:\путь\к\фото.jpg" #или относительный путь

curl.exe http://localhost:8080/photo/<PHOTO_ID> --output photo.jpg


11. Фильтрация заказов по времени:

curl.exe "http://localhost:8080/order?from=2026-01-01T00:00:00Z&to=2026-12-31T23:59:59Z"


## Реализовано

- Go 1.24, net/http с маршрутизацией по методам и path-параметрам.
- SQLite через modernc.org/sqlite (чистый Go, без CGO).
- Конфигурация через .env и godotenv.
- Структурированные JSON-логи через log/slog.
- Все 7 сущностей и все эндпоинты из ТЗ.
- Пагинация (limit, offset) для товаров, пользователей, заказов.
- Фильтрация товаров по типу, полу, размеру, цвету, бренду, сезону.
- Сортировка товаров по цене (ASC/DESC).
- Фильтрация заказов по статусу и интервалу времени.
- Валидация статусов заказа.
- SQL-миграция с таблицами и 4 базовыми ролями.
- Загрузка/отдача/удаление фото через multipart/form-data.
- Структура: /cmd, /internal, /migrations.

## Не реализовано

- Авторизация/аутентификация (middleware).
- Транзакции при оформлении заказа.

## Сброс проекта

Остановить сервер (Ctrl+C), затем в PowerShell:

Remove-Item store.db -ErrorAction SilentlyContinue
Remove-Item uploads\* -Recurse -ErrorAction SilentlyContinue


БД и роли создадутся заново.