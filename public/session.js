
const send_robin = document.getElementById('send-robin')
const send_sean = document.getElementById('send-sean')

send_robin.addEventListener('click', async (e) => {
    send(e.target.dataset.variant)
})

send_sean.addEventListener('click', async (e) => {
    send(e.target.dataset.variant)
})

async function send(variant) {
    const isSean = variant == "sean"

    const firstname = isSean ? "Sean" : "Robin"
    const lastname = isSean ? "Mcphilemy" : "Andreasson"
    
    await fetch("/session/createSession", {
        method: "POST",
        headers: {
            "content-type": "application/json"
        },
        body: JSON.stringify({
            firstname: firstname,
            lastname: lastname
        })
    })
}