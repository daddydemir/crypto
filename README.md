# CRYPTO
---
![crypto](https://github.com/daddydemir/crypto/assets/42716195/14f475ca-e6bb-4c42-8205-1b90bc70da18)


---
Crypto projesi coingecko'nun API'lerini kullanarak günlük, haftalık ve aylık olarak kripto borsadaki değişimleri raporlamak için oluşturduğum bir repodur. Bu raporları telegram'ın API'ini kullanarak telefonuma bildirim gönderecek şekilde geliştirdim.

### ÖZET
----
Coingecko'nun ücretsiz olarak sunmuş olduğu API'alert günde iki kez istek atıp verileri kendi veritabanıma kaydediyorum. Günde iki kez çalışması için farklı endpointler var ve bu endpointlere cron job sayesinde otomatik olarak istek atıyorum. İlk istek attığım zaman gece 00.15
ikinci istek ise 23.45. İkinci istek attığımda gelen veriler ile ilk istekteki veriler arasında bazı hesaplamalar yaparak veritabanına kayıt işlemini gerçekleştiriyorum. Aynı zamanda elde ettiğim hesaplamalar sonucunda bazı coinleri telefonuma göndermesi için kuyruğa atıyorum.
Aynı zamanda küçük bir python scripti sayesinde başka bir web sitesinde RSI grafiğini çekip resim olarak da teleonuma göndermesi için bir çalışma yaptım ancak o bu repo içerinde dahil değil. Henüz onu github'da paylaşmadım, eğer paylaşırsam buraya bir referans verebilirim.

### TEKNOLOJI
------

- ![GO](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white)
- ![RABBITMQ](https://img.shields.io/badge/RabbitMQ-FF6600.svg?style=for-the-badge&logo=RabbitMQ&logoColor=white)
- ![MYSQL](https://img.shields.io/badge/MySQL-4479A1.svg?style=for-the-badge&logo=MySQL&logoColor=white)
- ![DOCKER](https://img.shields.io/badge/Docker-2496ED.svg?style=for-the-badge&logo=Docker&logoColor=white)

### RUN | DEPLOY
---------
Projede veritabanı bağlantısı ve rabbitmq bağlantısı için bir .env dosyası bulunuyor, onları kendinize uygun bir şekilde değiştirdiğinizde docker container oluşturabilirsiniz. Bunun için aşağıdaki komutu kullanın.

```sh

docker build --tag crypto . 

```

> bu komutun sağlıklı bir şekilde çalışabilmesi için terminalde projenin Dockerfile dosyasının bulunduğu dizine gitmeniz ve orada çalıştırmanız gerekiyor.

```sh

docker run -p 8080:8080 crypto

```

bu komutu kullanarak oluşturulan containeri çalıştırabilirsiniz.

Eğer bana projeyle alakalı bir geliştirme veya bug yada herhangi bir sebepten ötürü ulaşmak isterseniz Discord'u kullanabilirsiniz. ( DEMİRON#1218 )
 
