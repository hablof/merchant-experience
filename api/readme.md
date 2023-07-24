JSON схема для передачи таблицы с товарами:

``` json
{
    "tableURL": "example.com/path",
    "sellerId": 42
}
```
Ответ в формате:
``` json
{
    "added": 10,
    "updated": 20,
    "deleted": 30,
    "errors": [
        {
            "row": 10,
            "field": "name",
            "errMsg": "too long name"
        }
    ]
}
```

URL схема для получения списока товаров из базы:

``` url
    host:port/?seller_id=1,2,3&offer_id=1,2,3&substring="substring"
```