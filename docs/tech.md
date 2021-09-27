Что по стэку
---
- база [neo4j](https://neo4j.com/download-center/#community) - как самая взрослая из графовых
- апишка на [golang + echo](https://github.com/labstack/echo) - требований мало, а экспертиза есть
- тула для заполнения базы golang - зачем расширять стек
- фронт на голом html + [vanilla.js](http://vanilla-js.com/) - фронтендера в команде нет, делаем максимально простой UI
- граф рисуем через [Ogma.js](https://doc.linkurio.us/ogma/latest/download.html) если дадут лицензию для опенсорса
- или страдаем и рисуем на [vis.js](https://visjs.github.io/vis-network/examples/network/nodeStyles/customGroups.html) 


Что по архитектуре
---

- Статику хостим на гитхабе (gh-pages)
- Базу пробуем вместить на самый [простой дроплет DO](https://www.digitalocean.com/pricing/)
- Апишку - рядом с БД или отдельно на бесплатный хост [heroku](https://www.heroku.com/)
- CI в репозитории на гитхабе через gh-Actions
