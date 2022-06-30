# Как запустить Борду локально

1. Сделать копию файла `.env.example` и назвать его `.env`. Задать в нём нужные переменные окружения.
   - `SERVER_ADDR` - адрес на катором будет запущена Борда
   - `DATABASE_URL` - адрес для подключения к базе данных
   - `CONTENT_URL` - ссылка до репозитория с задачами на GtiHub для работы с Git
   - `GITHUB_TOKEN` - токен для работы с GraphQL GitHub (персональный токен можно сгенерировать, как описано в [инструкции](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)).
2. Экспортировать переменные окрружения.

      set -o allexport; source .env; set +o allexport 

3. Запустить локальный веб-сервер командой `go run cmd/borda/main.go`