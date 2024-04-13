# Получение всех баннеров

## Описание

Ручка `GET /banner` нужна админам для просмотра баннеров. Доступна только по токену админа.
Опциональные параметры запроса `limit` `offset` `tag_id` `feature_id`.

## Примеры

1) Получение всего списка
![GetAll.png](GetAll.png)
2) Получение первой записи
![AllLimit1.png](AllLimit1.png)
3) Получение записей по фиче
![FeatureIDFilter.png](FeatureIDFilter.png)
4) Получение записей с фичей и лимитом
![FeatureIDLimit.png](FeatureIDLimit.png)
5) Получение одной записи со смещением 2 по списку
![Limit1Offset2.png](Limit1Offset2.png)
6) Получение записей по тегу
![TagIDFilter.png](TagIDFilter.png)