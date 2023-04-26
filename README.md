<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Version][version-shield]][version-url]
[![License](https://img.shields.io/github/license/serhiq/skye-trading-bot?style=for-the-badge)](https://github.com/serhiq/skye-trading-bot/blob/main/LICENSE)
![Last Commit](https://img.shields.io/github/last-commit/serhiq/skye-trading-bot.svg?style=for-the-badge)

<h3 align="center">Телеграмм бот для создания заказов</a></h3>

  <p align="center">
Телеграмм - бот для приложения "Мой ресторан" написанный на языке Go.
    <br />
    <!-- <a href="https://github.com/serhiq/skye-trading-bot"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/serhiq/skye-trading-bot">View Demo</a>
    · -->
    <a href="https://github.com/serhiq/skye-trading-bot/issues">Report Bug</a>
    ·
    <a href="https://github.com/serhiq/skye-trading-bot/issues">Request Feature</a>
  </p>


<!-- TABLE OF CONTENTS -->
* [About The Project](#about-the-project)
* [Requirements](#requirements)
* [Getting Started](#getting-started)
    * [Installation](#installation)
* [Usage](#usage)
* [License](#license)


<!-- ABOUT THE PROJECT -->
## About The Project

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->

Skye Trading Bot - это торговый бот для Telegram, который использует API Resto Evotor для создания заказов в ресторане. Бот позволяет пользователям просматривать меню, создавать и редактировать заказы, а также получать информацию о статусе и оплате заказа.

<div align="center">
    <img src="project/screenshots/sell_bot.gif" width="600px"/> 
</div>

<p align="right">(<a href="#readme-top">back to top</a>)</p>


|                                                                                           | Загрузка меню | Подтверждение заказа | Печать чека | комментарий                                           |
| ----------------------------------------------------------------------------------------- | ------------- | -------------------- | ----------- |-------------------------------------------------------|
| [Мой ресторан](https://market.evotor.ru/store/apps/06341a0a-a2d4-4d7f-a24f-fcc26531efb1)  | +             | +                    | +           | Рекомендуемый вариант                                 |
| [Чек по заказу](https://market.evotor.ru/store/apps/1ea7f5c6-bbfc-4379-9970-adde022ee6b8) |               | +                    | +           | загрузка меню через json файл, либо через ключ Эвотор |
| Файловый провайдер                                                                        | +             | -                    | -           | создан для тестирования                               |






<!-- GETTING STARTED -->
## Getting Started

### Installation


1. Создайте нового Telegram бота и получите его API токен. Инструкции по созданию Telegram бота доступны в [официальной документации Telegram](https://core.telegram.org/bots#3-how-do-i-create-a-bot).
2. Клонируйте репозиторий Skye Trading Bot на свой локальный компьютер:



```bash
git clone https://github.com/serhiq/skye-trading-bot.git
```
3. Перейдите в каталог с проектом:

```bash
cd skye-trading-bot/project
```

5. Создайте файл config.yaml в папке data на примере config_example.yaml.

Его структура такая:

* timezone - отвечает за форматирование даты в истории заказов
* telegram/token - токен для бота


* product_api - настройки внешнего API для загрузки номенклатуры
  * evo_api
  * resto_api
  * file_provider
* order_api - настройки внешнего API для создания заказов
- - order_api
- resto_api
- file_provider


В папке настроек есть отдельные примеры для настроек приложения Мой ресторан и Чек по заказу.

6. Запустите контейнеры приложения с помощью команды:
```bash
docker-compose up -d
```

## Usage

После успешной установки и запуска Skye Trading Bot можно начать просматривать меню и создавать заказы с помощью Telegram. .


<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the Creative Commons Licence. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


[version-shield]: https://img.shields.io/github/go-mod/go-version/serhiq/skye-trading-bot?filename=bot-service%2Fgo.mod&style=for-the-badge

[version-url]: https://github.com/serhiq/skye-trading-bot

[Golang]: https://img.shields.io/badge/Golang-000000?style=for-the-badge&logo=go&logoColor=white

[Golang-url]: https://go.dev/

