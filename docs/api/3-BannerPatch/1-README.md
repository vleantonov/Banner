# Обновление баннера

## Описание

Ручка `PATCH /banner/:id` также доступна только по токену админа. 
Тело запроса может содержать от одного до всех параметров `content`, `tag_ids`, `feature_id`, `is_active`.
Возвращает только стату код. Валидация полей аналогична, как в ручке `POST /banner`

## Примеры запросов

Исходный баннер

![First.png](First.png)

1) Включение баннера
![UpdateIsActive.png](UpdateIsActive.png)
![Second.png](Second.png)
2) Обновление контента
![UpdateContent.png](UpdateContent.png)
![Third.png](Third.png)
3) Обновление тегов
![UpdateTagIds.png](UpdateTagIds.png)
![Fourth.png](Fourth.png)
4) Обновление фичи
![UpdateFeatureID400.png](UpdateFeatureID400.png)
![UpdateFeatureID.png](UpdateFeatureID.png)
![Fifth.png](Fifth.png)
5) Обновление тега и фичи
![UpdateTagFeatureID.png](UpdateTagFeatureID.png)
![Six.png](Six.png)
6) Обновление всего 
![UpdateAll.png](UpdateAll.png)
![Last.png](Last.png)