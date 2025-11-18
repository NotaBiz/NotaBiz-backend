Project to help business owner maintain & track their accounting record

**Prerequisites :**
1. go language compiler https://go.dev/doc/install
2. redis server https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/
3. rabbitmq https://www.rabbitmq.com/docs/download
4. postgresql https://www.postgresql.org/download/
5. gmail account https://support.google.com/accounts/answer/185833
6. whatsapp API as message sender (example usage Fonnte https://docs.fonnte.com/)

**Configure Installation :**
1. clone the project
2. copy the buildscript, env, and dataseed.go in folder database. for env file, you must configure
   - database connection to postgre
   - app config
   - jwt config
   - google client config (https://support.google.com/googleapi/answer/6158849?hl=en)
   - logger config
   - redis config
   - rabbitmq config
   - email & smtp config (**use your password from app password you register above in your email**)
   - Fonnte token (**token from device connected to Fonnte**)
3. update dependencies with 'go mod tidy'
4. run redis & rabbitmq server
5. run project with 'go run .'
6. open swagger documentation in 'http://ip:port/api/v1/swagger/index.html'
7. opsional, you can run documentation from golang (https://github.com/amalmadhu06/godoc-example) with 'godoc -http=:6060' and open 'http://ip:6060'
