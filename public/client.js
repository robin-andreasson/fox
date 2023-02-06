

window.addEventListener("DOMContentLoaded", async () => {


    const test = await fetch("/cookies")


    console.log(await test.json())
})