# Cockatrice PrintLock

Cockatrice PrintLock is a Go-based command-line tool that works with **Cockatrice deck files** to ensure consistent, deterministic card printing and identification by locking cards to their **exact UUIDs and printings**.

The goal is to **save time when you want to avoid manually choosing the card art you want for your new deck**.

---

## ✨ Features

- 📦 **Parse Cockatrice `.cod` deck files**
- 🧬 **Preserve and lock card UUIDs**
- 🗂️ Accurately handle:
  - Card names
  - Card sets
  - Card numbers
  - Zones (main deck, sideboard, etc.)
- 🚀 Lightweight, fast, and dependency-minimal

---

## 🎯 Project Goals

### MVP Goals
- Reliably **load Cockatrice deck files**
- Preserve **exact card identity** (UUID + printing metadata)
- Write decks back to disk **without data loss**
- Reduce friction when creating new decks by automatically preserving preferred card art

---

## 🧪 Current Status

### ✅ Completed
- Deck XML parsing
- Zone-aware card loading
- UUID-safe card representation
- Safe file handling and validation
- Deterministic deck writing

---

## 🔮 Upcoming Features

- 📁 Allow updating deck files inside subfolders
- 🖥️ Upgrade the application from a simple CLI to a TUI (terminal user interface)

---

## 🧠 Why This Exists

When building or iterating on decks in Cockatrice, manually re-selecting card art for every new deck is tedious and time-consuming.  
This tool exists to **remove that friction**, letting you focus on deckbuilding—not reconfiguring visuals.

---

## 📜 License

MIT License.  
See `LICENSE` for details.
