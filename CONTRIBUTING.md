# Contributing to Syswatch Daemon

Öncelikle projeye katkı yapmak istediğiniz için teşekkürler! 🙌  
Bu doküman, katkı sürecinde izlemeniz gereken adımları ve kuralları özetler.

---

## 🚀 Nasıl Katkı Sağlayabilirim?

1. **Fork ve Clone**
   - Reponun bir fork’unu alın
   - Kendi bilgisayarınıza klonlayın

2. **Yeni Branch Açın**
   ```bash
   git checkout -b feature/ozellik-adi
   ```

* Branch isimleri `feature/`, `fix/`, `docs/` gibi öneklerle başlamalı.

3. **Kodlama Standartları**

   * Kodunuz **temiz**, **okunabilir** ve **yorumlarla desteklenmiş** olmalı.
   * Commit mesajları açıklayıcı yazılmalı:

     * ✅ `fix: memory leak in process watcher`
     * ✅ `feat: add config file parser`
     * ❌ `update stuff`

4. **Testler**

   * Yeni özellik eklerken mümkünse test yazın.
   * Tüm testlerinizi çalıştırın ve geçtiklerinden emin olun.

5. **Pull Request**

   * PR açıklamasında yapılan değişiklikleri net olarak yazın.
   * PR’ınızı ilgili **issue** veya **milestone** ile ilişkilendirin.

---

## 📌 Issue Açmadan Önce

* Aynı sorun daha önce açılmış mı kontrol edin.
* Açıklamanız **net**, **teknik** ve gerekirse loglarla desteklenmiş olsun.
* Gereksiz veya yinelenen issue’lar kapatılacaktır.

---

## 🔑 Genel Kurallar

* `main` branch **her zaman stabil** kalmalı.
* Büyük değişiklikler öncesinde **tartışma başlatın** (discussion veya issue).
* Katkılarınız [Code of Conduct](CODE_OF_CONDUCT.md) belgesine uygun olmalı.
