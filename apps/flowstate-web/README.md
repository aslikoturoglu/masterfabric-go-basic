# FlowState AI — MasterFabric üzerine inşa

Haftalık Optimizasyon Motoru. MasterFabric Go Basic'in GraphQL API'si, auth ve altyapısı kullanılır.

## Çalıştırma

1. MasterFabric backend'i başlatın:
   ```bash
   make docker-infra
   set -a && source .env && set +a && make run
   ```

2. FlowState frontend:
   ```bash
   make flowstate-web
   ```
   veya
   ```bash
   cd apps/flowstate-web && npm install && npm run dev
   ```

- Backend: http://localhost:8080
- FlowState: http://localhost:3001

## GraphQL Endpoints (FlowState)

- `flowstateFixedEvents` — Sabit program listesi
- `flowstateCreateFixedEvent` — Yeni sabit etkinlik
- `flowstateDeleteFixedEvent` — Etkinlik sil
- `flowstateFlexibleTasks` — Esnek görev listesi
- `flowstateCreateFlexibleTask` — Yeni görev
- `flowstateDeleteFlexibleTask` — Görev sil
- `flowstateGenerateSchedule` — AI ile haftalık takvim oluştur
- `flowstateSchedule(week)` — Oluşturulmuş takvimi getir
