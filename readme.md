Можно еще добавить тест кейсы на остальные хендлеры, graceful shutdown, привести в порядок код, добавить документацию,
но желания заниматься этим в тестовом задании как-то не особо.




Задание:
Web-Сервис сотрудников, сделанный на Golang Сервис должен уметь:
1. Добавлять сотрудников, в ответ должен приходить Id добавленного сотрудника.
```
curl localhost:8080/employee -X POST -d '{"Name":"employee1Name","Surname":"employee1Surname",
"Phone":"employee1Phone","CompanyId":3,"PassportType":"employee1PassportType",
"PassportNumber":"employee1PassportNumber","DepartmentName":"employee1DepartmentName",
"DepartmentPhone":"employee1DepartmentPhone"}'
```
2. Удалять сотрудников по Id.
```
curl localhost:8080/employee/3 -X DELETE
```

3. Выводить список сотрудников для указанной компании. Все доступные поля.
```
curl localhost:8080/company/3 -X GET -v
```


4. Выводить список сотрудников для указанного отдела компании. Все доступные поля.
```
curl localhost:8080/company/555/department/employee1DepartmentName123 -X GET
```
5. Изменять сотрудника по его Id. Изменения должно быть только тех
   полей, которые указаны в запросе.
```
curl localhost:8080/employee/5 -X PATCH -d '{"Name":"NewName123;drop table employees", 
"test": 123,"CompanyId": 555}'
```



Модель сотрудника:
```
{
   Id int
   Name string 
   Surname string 
   Phone string 
   CompanyId int 
   Passport {
        Type string
        Number string 
   }
   Department { 
        Name string
        Phone string
   }
}
```
Все методы должны быть реализованы в виде HTTP запросов в формате JSON. БД: любая.



```
mockgen -destination=mocks/mock_storage.go -package=mocks  github.com/nktau/employee-service/internal/storagelayer Storage

```

Запуск базы:
```
docker run -d --name employee_service_postgres -e POSTGRES_PASSWORD=mysecretpassword -v ./pg_data:/var/lib/postgresql/data -p 5433:5432 --rm postgres:14
```
запуск приложения:
```
go run cmd/main.go
```
запуск тестов:
```
go test ./...
```