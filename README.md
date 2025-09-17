# Syswatch Daemon

**Linux sistemleri için kernel seviyesinde çalışan, gerçek zamanlı sistem izleme ve uyarı aracı. Cluster desteği ve anomaly detection içerir.**  

---

## 🚀 Proje Tanıtımı

Syswatch Daemon, Linux sistemlerini izlemek ve performans sorunlarını önceden tespit etmek için geliştirilmiş bir araçtır.  
Gerçek zamanlı CPU, RAM, Disk, GPU takibi sağlar ve birden fazla sunucu üzerinde cluster desteği ile çalışabilir.  

**Öne çıkan özellikler:**
- Kernel seviyesinde CPU, RAM, Disk, GPU izleme
- Gerçek zamanlı anomaly detection ve uyarılar
- Cluster destekli çok sunuculu yapı
- Opsiyonel web dashboard veya grafiksel arayüz
- Docker ve Kubernetes uyumlu deployment

---

## 📈 Yol Haritası

**Başlangıç aşaması (1-2 hafta):**
- Go ile temel daemon geliştirme
- Temel metric toplama ve loglama
- Basit dashboard (terminal veya web) hazırlama

**Orta aşama (2-6 hafta):**
- Anomaly detection algoritmaları ekleme
- Cluster-wide metric toplama ve konsolidasyon
- Grafana / Plotly dashboard entegrasyonu

**İleri aşama (6+ hafta):**
- Prometheus + Alertmanager entegrasyonu
- Telegram / Slack bot ile gerçek zamanlı uyarı sistemi
- Kubernetes deployment ve ölçeklendirme

---

## 🛠 Proje Adımları

1. Repo’yu klonla:
```bash
git clone https://github.com/kullaniciadi/syswatch-daemon.git
````

2. Go ortamını kur ve daemon’u çalıştır:

```bash
cd syswatch-daemon
go build -o syswatch-daemon
sudo ./syswatch-daemon
```

3. Dashboard veya görselleştirme araçlarını kur:

* Grafana / Plotly ile metric görselleştirme
* Docker container deployment

4. Cluster ve anomaly detection ekle:

* Çok sunuculu metric toplama
* Anomaly detection algoritmaları
* Real-time alert sistemi

---

## ✅ Yapılacak İşler (To-Do List)

### Başlangıç Aşaması (1-2 hafta)

* [ ] Temel Go daemon’u oluştur
* [ ] Metric toplama ve loglama fonksiyonları ekle
* [ ] Basit terminal/web dashboard

### Orta Aşama (2-6 hafta)

* [ ] Cluster-wide metric toplama
* [ ] Anomaly detection algoritmaları
* [ ] Grafana / Plotly entegrasyonu

### İleri Aşama (6+ hafta)

* [ ] Prometheus + Alertmanager entegrasyonu
* [ ] Real-time alert bot (Telegram/Slack)
* [ ] Kubernetes deployment ve ölçeklendirme
* [ ] Performans testleri ve optimizasyon

---

## ⚡ Not

Proje geliştirmeye **kısa süre içinde başlanacaktır** ve roadmap’deki adımlar sırasıyla ilerleyecektir. Her aşamada yeni özellikler eklenerek, **production-ready Linux sistem izleme daemon’u** hedeflenmektedir.

---

## 📄 Lisans

MIT License
