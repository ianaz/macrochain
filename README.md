# 🌍 Macrochain — Global Macro + Ethereum Intelligence

**Macrochain** is a data-rich, modern web application that brings together global macroeconomic indicators and Ethereum ecosystem metrics in a clean, intuitive platform designed for curious retail users.

---

## 📌 Overview

**Macrochain** provides:
- Real-time macroeconomic indicators (e.g. inflation, central bank decisions, interest rates, commodities)
- Ethereum-specific metrics (price, gas, staking, TVL, L2 fees, DEX volume)
- DeFi and NFT analytics
- A powerful backend infrastructure to collect and unify data from both traditional finance and Ethereum
- Other helpful data

---

## 🎯 Goals

- Make Ethereum and macro data readable and accessible for everyday users
- Show central bank events and macroeconomic movements that matter (e.g. FED decisions, SNB rate changes, ECB meetings, commodity prices)
- Provide an intuitive, responsive platform with a rich set of visual tools and insights
- Deliver both live data and context, with light educational content for deeper understanding

---

## 🤩 Features

- 🌐 **Dashboard** with key macro + Ethereum metrics
- 📊 **Deep Dive** sections for Macro and Ethereum data
- 🔁 **Compare** section to visualize correlations
- 🧠 **Learn**: Interactive content to explain key concepts
- 📚 **Glossary**: Crypto + Macro terms explained
- ☁️ Real-time or near real-time data
- 🌓 Light/Dark mode
- 🌍 Multilingual
---

## 🏗️ Tech Stack

### 💻 Frontend
- **Vue.js** 
- **Tailwind CSS** (or similar utility-first framework)
- **Chart.js / Recharts / TradingView Embeds** for visualizations
- **i18next** for internationalization

### ⚙️ Backend & Scrapers

The most critical part of the system is a **scraper-based backend infrastructure**.

- **Language**: **Golang**
- **Architecture**: Multiple small and independent scrapers, each handling a specific data domain (e.g. FED, SNB, Ethereum on-chain, DeFi protocols,...)
- **Data Storage**:
  - **PostgreSQL** for structured + relational data
  - Optionally: time-series extension (TimescaleDB) for long-term trends
- **Queue System**:
  - **Recommended (Simple)**: **Redis** (easy to integrate, works well for home setup, can be used for pub/sub or simple job queue)
- **Orchestration & Scheduling**:
  - Simple **cron-like scheduling** inside each container
  - Scrapers run as services inside `docker-compose`
- **API Layer**:
  - RESTful API in Go (or lightweight framework like `chi`, `echo`, or `fiber`)
- **Migration**:
  - For migration we'll use golang-migrate lib
- **Deployment**:
  - Everything runs locally on a private server via **Docker Compose**

---

## 📁 Project Structure (planned)

```
macrochain/
├── webapp/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── locales/
│   │   └── App.jsx
├── scraper/
│   ├── cmd/
│   ├── pkg/
│   │   ├── fed/
│   │   ├── snb/
│   │   ├── ethereum/
│   │   ├── defi/
│   │   ├── .../
│   │   └── common/
│   └── main.go
├── api/
│   ├── ... package by feature
│   ├── ... package by feature
│   └── main.go
├── db/
│   ├── migrations/
│   └── schema.sql
├── docker-compose.yml
└── README.md
```

---

## 💡 Roadmap

- [x] UI mockup + initial HTML prototype
- [x] Defined frontend sections and layout
- [x] Project structure and service layout defined
- [ ] Implement React-based dashboard and navigation
- [ ] Set up PostgreSQL schema with TimescaleDB extension
- [ ] Implement Redis-based job queue (minimal pub/sub)
- [ ] Build core scrapers in Go for: FED, SNB, ETH on-chain, DEX, Staking
- [ ] Expose basic REST API for frontend consumption
- [ ] Deploy and run all services with Docker Compose
- [ ] Create Glossary and Learn static content pages
- [ ] Add i18n support for future multilingual expansion (EN MVP)

