# LRU-Cache

"Задание: ""Написать сервис (класс/структура) кэширования"".

Основные условия:
- Кэш ограниченной емкости, метод вытеснения ключей LRU;
- Сервис должен быть потокобезопасный;
- Сервис должен принимать любые значения;
- Реализовать unit-тесты.

Сервис должен реализовывать следующий интерфейс:
````
type ICache interface {
    Cap() int
    Len() int
    Clear() // удаляет все ключи
    Add(key, value any)
    AddWithTTL(key, value any, ttl time.Duration) // добавляет ключ со сроком жизни ttl
    Get(key any) (value any, ok bool)
    Remove(key any)
}
````
TTL - через сколько должен удалиться ключ"