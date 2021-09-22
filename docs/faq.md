FAQ
---

---
> Четвёртый сценарий предполагает под собой что-то вроде "теории шести рукопожатий"? 
> Просто в одной ситуации там 1 ребро где есть искомый фильм и актер, который в нем снимается. 
> В другой ситуации есть N ребер и K вершин, которые связывают актера через кучу фактов и условно друзей с искомым фильмом. 
> Какое направление верное?

Мысли о теории рукопожатий верные - это примерно оно и есть. ""
В простейшем случае в этом сценарии будет граф из двух вершин вида "фильм" -> "актёр" (если мы искали связь от Джонни Депа до Пиратов карибского моря).

В более интересном случае это десяток фильмов и людей со связями вида "сценарист того же фильма где снималась искомая персона играл эпизодическую роль в фильме, в котором режисёром был тот же челочек, который играл в ..."

--- 

> Понимаю отношение фильм – актер, но осталось непонятным роль актера. 
> В виде ребра оно скучно выглядит, т.е если будет показывать у актера роль. 
> В ином случае мне в голове рисуется кластер (роль), в которой обединено N актеров. Пикча

Над визуализацией предлагаю думать уже на реальных данных. Я бы начал с максимально простой модели где:
- рёбра {актёр+роль, сценарист, режисёр, etc}
- вершины {фильм, человек}

---

> Непонятно еще отношения фильм – фильм. По какому принципу они связываются? По актерам?

Смотри прошлый вопрос. Фильм-фильм возможен только через человека, который так или иначе связан с обоими фильмами.
Связи напрямую фильмов я бы предложил оставить на будущее, потому что пока в данных ничего похожего не видно.

---

> Для отрисовки графа будет использоваться готовая либа, какая? 
> у них есть свои ограничения на кастомизацию, понимать бы. 
> Или ручками сами на канвасе или в свг рисуем? (последний вариант кажется гемором т.к я с таким не работал :D)

На первые итерации точно возьмём готовую либу и чутка подкрасим. Примеры готовые есть в доке рядом (@see graph-ui.md)

---

> Что будет делать suggestion-popup? Он будет работать по клику на вершину или что-то другое?

Это просто подсказка для формы поиска, поскольку нам нужен конкретный идишник фильма/человека, чтобы строить от него граф.
Вроде бы это называется саджесшон попап?)

---

> Как работает фильтр movie-only?
 
В первой итерации даём искать только по именам фильмов, вот и весь фильтр.

---