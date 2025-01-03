curl -X GET http://localhost:8080/tasks \
-H "Content-Type: application/json"



curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
    "id": "3",
    "description": "Изучить Docker",
    "note": "Начать с основ контейнеризации",
    "applications": [
        "Docker Desktop",
        "Terminal",
        "VS Code"
    ]
}'



# Получение существующей задачи (ID = 1)
curl -X GET http://localhost:8080/tasks/1 \
-H "Content-Type: application/json"

# Получение несуществующей задачи (для проверки обработки ошибок)
curl -X GET http://localhost:8080/tasks/999 \
-H "Content-Type: application/json"



# Удаление существующей задачи
curl -X DELETE http://localhost:8080/tasks/1 \
-H "Content-Type: application/json"

# Удаление несуществующей задачи (для проверки обработки ошибок)
curl -X DELETE http://localhost:8080/tasks/999 \
-H "Content-Type: application/json"


# Создание задачи с минимальными данными
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
    "id": "4",
    "description": "Простая задача",
    "note": "",
    "applications": []
}'

# Создание задачи с существующим ID (для проверки обработки ошибок)
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
    "id": "1",
    "description": "Тестовая задача",
    "note": "Эта задача не должна быть создана",
    "applications": ["test"]
}'

# Создание задачи с некорректным JSON (для проверки обработки ошибок)
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
    "id": "5"
    "description": "Некорректный JSON"
}'