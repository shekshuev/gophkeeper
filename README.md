# GophKeeper

> ⚠️ **Disclaimer**  
> This repository contains the final graduation project from the [Yandex Practicum Advanced Go Developer Course](https://practicum.yandex.ru/go-advanced/).  
> If you're looking to become a better software developer — **run away from Go**. Seriously. It's the most mind-numbingly boring language on Earth. It won’t teach you to think, it will teach you to write boilerplate. You won't become a craftsman — you'll become a typist.
> ![Fly, you fools!](cover.jpg)
> If, however, you just want money — then Go is perfect for you.

---

## 🗄 About the Project

**GophKeeper** is a secure password and data manager with client-server architecture.

This application allows users to register, authenticate, and safely store various types of secrets — like passwords, secure notes, or card data — in encrypted form. It supports syncing between multiple clients belonging to the same user.

This project is focused on backend security, data consistency, and test coverage. It doesn't aim to be fancy — it aims to be rock-solid.

---

## Features

### Server

- User registration and authentication
- Authorization via JWT
- Encrypted storage of private user data (`JSONB`)
- Secret types: passwords, notes, card data, binary blobs
- Synchronization support between multiple clients
- REST API with clean architecture and repository pattern
- Integration and unit tests
- Logging and error tracing

### Client (CLI)

- Secure authentication and storage
- CLI interface using `prompt` for input
- Fetch individual secrets or list all secrets
- Auto-sync with the server
- Separate token management (access + refresh tokens)

## Final Thoughts

I deeply regret spending time on this course — and on this language.

The [Yandex Practicum Advanced Go Developer Course](https://practicum.yandex.ru/go-advanced/) was one of the most uninspiring educational experiences I've had. Instead of teaching problem-solving, architecture, or trade-offs, it drilled repetitive, mechanical patterns with little to no intellectual depth. The course treats engineering like checklist-driven typing practice.

There is no emphasis on elegance, design thinking, or creativity — just "follow the spec", "write a handler", "inject a repository", and "cover with tests". That’s not how you grow as an engineer. That’s how you learn to imitate without understanding.

Go, as a language, only amplifies this effect. Its minimalism is often praised, but in practice it strips away too much — forcing developers to repeat themselves, avoid abstraction, and settle for the lowest common denominator in software design.

If you care about craftsmanship, insight, and long-term skill — **look elsewhere**.

If you only care about performance, simplicity, and a quick entry into backend jobs — Go will do the job.

But don't confuse mechanical work with engineering.

> Penned during sleepless nights under the flickering light of pointless for `if err != nil { return nil, err }`. Not in pursuit of mastery — but to appease a checklist. Not out of passion — but out of obligation.This is not how great software is born. This is how you burn out quietly, one interface at a time.
