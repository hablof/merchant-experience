## Важное замечание по структуре models.Product
Хотя структура и содержит поле `SellerId`, оно не используется в методе (s *Service) UpdateProducts.


## метод (s *Service) UpdateProducts

Должен обновить данные в репозитории

- получает на вход ID продовца `sellerId` и информацию об обновлениях продуктов `productUpdates`
- должен вернуть информацию об обновлении репозитория `UpdateResults`: количество созданных товаров, обновлённых, удалённых, а также ошибки валидации данных; ошибку сервиса/репозитория

Для успешной работы делает следующее:
<!-- 1. Валидирует входящую информацию:
    создаём слайс `validatedUpdates` с ёмкостью равной длине `productUpdates`, добавляем туда все элементы, прошедшие валидацию. -->
1. Вызывает метод репозитория `SellerProductIDs` чтобы получить все айдишники продавца `sellerId`. Сортируем айдишники (далее будем использовать бинарный поиск).
2. Разбирает входящие `productUpdates` на три категории: 
- продукты которые необходимо удалить (имеют значение `false` в поле `Available`)
- продукты которые необходимо изменить (имеют совпадение поля `OfferId` с одним из имеющихся товаров (бинарный поиск))
- продукты которые необходимо добавить (НЕ имеют совпадение поля `OfferId` с одним из имеющихся товаров (бинарный поиск))
3. Пробегается валидацией по продуктам, которые необходимо изменить/добавить. Невалидные выкидываются.
    - условия валидации: *количество символов* в строке `Name` не больше 100. 
     <!--SQL defines two primary character types: character varying(n) and character(n), where n is a positive integer. Both of these types can store strings up to n characters (not bytes) in length.
     https://www.postgresql.org/docs/15/datatype-character.html   -->
4. вызывает метод репозитория `ManageProducts`
5. собирает длины слайсов в соответствующие поля структуры `UpdateResults`. Ошибки валидации и ошибки репозитория в поле `Errors`.

к сожалению, из-за необходимости отдельно подсчитать количество продуктов которые обновляются и добавляются сложность алгоритма -- O(n log n)

## метод ProductsByFilter
не содержит логики домена, а только не пропускает конкретную ошибку репозитория наружу