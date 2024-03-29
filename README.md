# Объем данных

- Основные данные: незначительный объем (менее 100 MB):
    - до 100 докладов на одну конференцию
    - до 3000 пользователей на одну конференцию
    - это даёт до 300 000 оценок на одну конференцию, где каждая оценка — это пара чисел от 1 до 5 плюс текстовый комментарий
- Вспомогательные данные (текущее состояние чатов): зависит от реализации

# Важно

Для упрощения реализации нет требования сохранять целостность данных в БД бота: в БД могут быть доклады, не относящиеся к текущей конференции (по дате начала), могут быть оценки, не относящиеся к существующим докладам (по URL доклада), в избранных докладах могут быть отсутствующие в БД, etc. Поэтому все данные по докладам и оценкам в БД считаются либо актуальными (если относятся к докладам с датой начала, попадающей в интервал работы текущей конференции) либо неактуальными (такие просто игнорируются)
    - Бот даёт возможность пользователям (включая админов) работать только с актуальными данными


# Чеклист

`!!!` - опциональные

## Админы

- [ ] Залить файл с расписанием докладов
    - [ ] Если данные в файле не проходят валидацию (описанную в разделе Данные),
    то весь файл игнорируется, и админ получает уведомление об ошибке(ах)
    - [ ] Доклады в БД с URL, отсутствующим в этом файле, удаляются из БД
        - [ ] Все пользователи получают уведомление о таких докладах
        - [ ] Из БД удаляются только доклады, остальные данные (наличие URL
        этого доклада в “Избранном”, оценки этого доклада)
        сохраняются (но не выводятся пользователю).
        Если доклад с этим URL будет снова добавлен, то связанная
        с ним информация снова должна показываться пользователю
    - [ ] Доклады в файле с URL, существующим в БД, обновляются в БД
        - [ ] Если у доклада изменилось время начала, то все пользователи получают
        уведомление о таких докладах
    - [ ] Доклады в файле с URL, не существующим в БД, добавляются в БД
        - [ ] Все пользователи получают уведомление о таких докладах
    - [ ] Админ получает уведомление о докладах в файле, у которых время начала
    не попадает в интервал работы конференции (такие доклады хранятся в БД,
    но нигде кроме этого уведомления, не показываются)
- [ ] Скачать файл с оценками
- [ ] !!! Задать/изменить/удалить текст общего уведомления
    - Этот текст (если он непустой/удален) отправляется:
        - всем текущим пользователям в момент задания/изменения этого текста
        - новым пользователям в момент подключения бота

## Обычные юзеры

- [ ] Команда просмотра информации о конференции
- [ ] Просмотреть расписание докладов текущей конференции
- [ ] Редактировать список избранных докладов
- [ ] Редактировать свою идентификацию (номер билета/ФИО или др идентификатор)
- [ ] Редактировать/Удалять свои оценки докладов
    - [ ] функционал доступен начиная со времени начала доклада и заканчивая временем завершения оценок конфы
    - [ ] если юзер не задал идентификацию, после добавления/изменения оценки ему нужно присылать уведомление, что его оценка не будет учтена, пока он не идентифицирует себя
- [ ] Просмотр своих оценок на все доклады
- [ ] Просмотреть все доклады без оценок
- Напоминания для обычных юзеров:
    - [ ] За 10 минут до начала доклада: о начале доклада
        - для всех докладов, если нет "избранных" докладов, иначе только для "избранных"
    - [ ] После завершения доклада запрашивать оценку доклада, если она еще не задана
    - [ ] В конце дня (через 1 час после завершения последнего доклада): запросить оценку по всем докладам этого дня, по которым она не задана
    - [ ] Через 2 дня после завершения конфы: запросить оценку по всем  докладам конфы, по которым она не задана

## Другое

- [ ] Настроить Линтер и запуск линтера в CI
- [x] Билд и запуск тестов в CI
- [x] Конфигурация через переменные окружения
- [ ] !!! Добавить рейт лимит в конфиг (https://github.com/tucnak/telebot/pull/487)
- [x] !!! Структурное логирование
- [x] !!! Поддержка выбора формата логов (text, json)
- [ ] !!! Служебные метрики (4 Golden Signals для API/БД/ресурсов)
- [ ] !!! Бизнес-метрики (количество пользователей/докладов/оценок/…)
- [ ] !!! Корректное восстановление состояния текущих диалогов после креша/перезапуска (скорее всего, потребует получения истории всех чатов через API Telegram, синхронизации её с данными в БД бота и продолжения текущего диалога при необходимости)

# Формат файлов

## Расписание

```csv
Start (MSK Time Zone),Duration (min),Title,Speakers,URL
21/07/2024 10:00:00,45,"Прагматичная архитектура: как проектировать Go-приложение, руководствуясь чисто практическими соображениями","Даниил Подольский, YADRO",<https://conf.ontico.ru/lectures/5088349/discuss>
22/07/2024 10:00:00,45,"Алиса? Салют? Алекса? Нет, <Название ассистента в разработке> !",Эдгар Сипки,<https://conf.ontico.ru/lectures/5389176/discuss>
```

## Оценки

Этот файл содержит только записи, для которых истинно всё перечисленное:
- Они относятся к актуальным докладам
- Тип ответа задан как “Оценка”
- У пользователя задана Идентификация
```json
[
	{
		"user": "номер билета или другой текст (email/ФИО/...)",
		"url": "<https://conf.ontico.ru/lectures/5088349/discuss>",
		"content": 5,        // Целое число от 1 до 5.
		"performance": 3,    // Целое число от 1 до 5.
		"comment": "опциональный комментарий"
	},
	...
]
```
