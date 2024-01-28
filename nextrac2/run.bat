set NEXTRAC2_HOST=0.0.0.0
set NEXTRAC2_PORT=9078
set NEXTRAC2_RESOURCE_ID=trac
set NEXTRAC2_DB_CONNECTION=user=postgres password=paramadaksa dbname=nexSOFT sslmode=disable host=localhost port=5432
set NEXTRAC2_DB_SCHEMA=resource_nextrac_local
set NEXTRAC2_DB_VIEW_CONNECTION=user=postgres password=paramadaksa dbname=nexSOFT sslmode=disable host=localhost port=5432
set NEXTRAC2_DB_VIEW_SCHEMA=resource_NEXTRAC2
set NEXTRAC2_REDIS_HOST=localhost
set NEXTRAC2_REDIS_PORT=6379
set NEXTRAC2_REDIS_DB=0
set NEXTRAC2_REDIS_PASSWORD=
set NEXTRAC2_CLIENT_SECRET=9b4969403e984f9eacde8ad705d3137c
set NEXTRAC2_SIGNATURE_KEY=a478e88cab8f48cdaf58dce65a4df68d
set NEXTRAC2_ELASTIC_SEARCH=http://10.10.11.195:9201
go run main.go staging

