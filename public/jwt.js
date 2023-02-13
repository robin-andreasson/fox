let token = ""

document.getElementById('login').addEventListener('submit', async (e) => {
    e.preventDefault()
    const res = await fetch("/auth/login", {
        method: "POST",
        headers: {
            "content-type": "application/x-www-form-urlencoded"
        },
        body: "user[person][username]=robin&user[person][password]=123"
    })


    const data = await res.json()

    console.log(data)
})

document.getElementById('check-token').addEventListener('click', async (e) => {
    e.preventDefault()
    const res = await fetch("/validate/token", {
        method: "POST",
        headers: {
            "authorization": "bearer: " + token
        }
    })


    const data = await res.json()

    console.log(data)
})