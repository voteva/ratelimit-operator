Для написания тестов используется фреймворк ginkgo(https://github.com/onsi/ginkgo) и библиотека для asserts 
gomega (https://github.com/onsi/gomega). Для работы  необходимо установить их в локальный кэш:

go get -u github.com/onsi/ginkgo/ginkgo

go get -u github.com/onsi/gomega/...

Водная статья для работы с ginkgo - https://medium.com/boldly-going/unit-testing-in-go-with-ginkgo-part-1-ce6ff06eb17f

Запуск тестов производится из корневой директории: go test ./... или ginkgo -r

В директории test лежит схема работы StubController и StubClient. XML файл (stubs_scheme.xml) открывается 
на сайте draw.io. Пример структуры теста: test/controller/some_test.go