
# JWT Auth Service | Test Assignment
Тестовое задание - разработка части сервиса аутентификации

## Запуск
```bash
docker-compose -f docker-compose.yml up -d
```

При выполнении команды запускается сервер приложения на порте `:8080` и сервер `PostgreSQL` (`:5432`) – конфигурируется через `.env` файл.

## Описание API
### Генерация пары токенов
```bash
curl "/generate?123e4567-e89b-12d3-a456-426614174000" \
-H 'Accept: application/json'
```

Пример ответа:
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 14 Jul 2025 21:54:22 GMT
Content-Length: 844
Connection: close

{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI1MzAzNjEsImlhdCI6MTc1MjUzMDA2MSwianRpIjoiYzEwMjdmOTMtYzFkZC00MGFjLThkYjEtN2ZmOTcwYzZkOWMwIiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.nmckalIdM9T_iA6BDFw_k8YJExKsmfioCo92H2eOGHe4PT29BwTKcFcSeYzB-O1tsOvvcMl1apEmT-ZnN4T2Xg",
"refresh_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI3MDI4NjEsImlhdCI6MTc1MjUzMDA2MSwianRpIjoiYjFlY2ViZjEtMWZkOC00NzQ2LTgwNDctYzMwOGI2NDk0MjI1IiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.ZR84IpSW_kJZrqhmUCboda6ZV7hyeCm6k-uyOtAdyjrrG-vgiRrIAagIkv6vPr-4lRgLI3zZCAJslY4j7LZHyg"}
```

### Обновление пары токенов
```bash
curl -X "POST" "/refresh" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI1Mjg1NTQsImlhdCI6MTc1MjUyODI1NCwianRpIjoiMDk0YjRlZmUtMzc5Ny00ZDI2LTkwMDItMDMwYTE0ZTU1NGUxIiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.xQItoHZo0pmXGKQruPZnQUQk7FoHBTOgKlpki9zsrhjM6O6iGKf27SO0NHu_X_NsjZA40ihWj8NlRYU33a7xNw",
  "refresh_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI3MDEwNTQsImlhdCI6MTc1MjUyODI1NCwianRpIjoiOTEzYjVhNDktMWE2Zi00M2U1LWFmNGItMTIwN2NlNzY4YjNiIiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.UWFMJ7pAZpUxahKBwYVIoaDbMQY-yUBbbKsct0h73jM-ISIeJ6uv2OfZ78dvmy0WrETPkX4-YVsSkIznpYOp2Q"}'
```

Пример ответа:
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 14 Jul 2025 21:24:58 GMT
Content-Length: 844
Connection: close

{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI1Mjg1OTgsImlhdCI6MTc1MjUyODI5OCwianRpIjoiOWIyMGIxNjItNGY1Zi00MmNlLWI2ZjktYjgxNzU1NDM3NmY2IiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.FiMgu0U_WhrwsBN72eESKYYcrSAte57f59dO6uhsDOSeYa0xfH86d5uEVdTUKU2OSrGML4CLfh8nDypkNgRmIw",
"refresh_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI3MDEwOTgsImlhdCI6MTc1MjUyODI5OCwianRpIjoiYmEyY2NkMTktOTYxOS00ODVkLWFkZmMtOWJlZmQwOTBjYmQ3IiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.KxGMpl-tQCIpWJZuasrxVAXgTqrHlJ75GfjsUKWs2DsFl4zsrpFX1hRsJ7BOkDR3twFbncK8Zy5KMw--dI6gaA"}
```

### Получение GUID текущего пользователя
```bash
curl -X "POST" "/me" \
     -H 'Accept: application/json' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI1MzAzNjEsImlhdCI6MTc1MjUzMDA2MSwianRpIjoiYzEwMjdmOTMtYzFkZC00MGFjLThkYjEtN2ZmOTcwYzZkOWMwIiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.nmckalIdM9T_iA6BDFw_k8YJExKsmfioCo92H2eOGHe4PT29BwTKcFcSeYzB-O1tsOvvcMl1apEmT-ZnN4T2Xg"
}'
```

Пример ответа:
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 14 Jul 2025 21:54:39 GMT
Content-Length: 50
Connection: close

{"user_id":"123e4567-e89b-12d3-a456-426614174000"}
```

### Деавторизация пользователя
```bash
curl -X "POST" "/logout" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
     "user_id": "123e4567-e89b-12d3-a456-426614174000",
     "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3NTI1MzAzNjEsImlhdCI6MTc1MjUzMDA2MSwianRpIjoiYzEwMjdmOTMtYzFkZC00MGFjLThkYjEtN2ZmOTcwYzZkOWMwIiwiVXNlckFnZW50IjoiUmFwaWRBUEkvNC4zLjQgKE1hY2ludG9zaDsgT1MgWC8xNS41LjApIEdDREhUVFBSZXF1ZXN0IiwiVG9rZW5UeXBlIjoiIn0.nmckalIdM9T_iA6BDFw_k8YJExKsmfioCo92H2eOGHe4PT29BwTKcFcSeYzB-O1tsOvvcMl1apEmT-ZnN4T2Xg"
}'
```

Пример ответа:
```http
HTTP/1.1 200 OK
Date: Mon, 14 Jul 2025 21:54:46 GMT
Content-Length: 0
Connection: close
```