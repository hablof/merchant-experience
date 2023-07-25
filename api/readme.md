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
    host:port/?seller_id=15&offer_id=1,2,3&substring="substring"
```
Ответ в формате:
``` json
[
    {
        "sellerId": 15,
        "offerId": 1,
        "name": "name1",
        "price": 100500,
        "quantity": 150
    },
    {
        "sellerId": 15,
        "offerId": 2,
        "name": "name2",
        "price": 1000,
        "quantity": 10
    }
]
```