# Phoceene Auto

This project is a simple Web App made with Golang to manage clients, cars and transactions for **Phoceene Auto**, a car dealer company.

> Authors:
> 
> - Antoine de Barbarin
> - Nicolas Moyon
> - Sabrina Eloundou

Task repartition:
- Organization & project setup: `Antoine`
- Database design: `Antoine`
- Database migrations: `Nicolas`
- Backend: `Antoine` & `Nicolas`
- Backend integration: `Antoine`
- Frontend: `Sabrina`
- README: `Antoine`

---

## Deployment (Linux)

- Create a PostgreSQL database (and a dedicated user).
- Copy the `.envrc.example` and name it `.envrc`, then modify the environment variables.
- Run the migrations using [`golang-migrate`](https://github.com/golang-migrate/migrate):
    ```bash
    make db/up
    ```

- Run the development server:
    ```bash
    make run
    ```

Open [http://localhost:8080](http://localhost:8080) (or whatever you configured) with your browser to see the result.
