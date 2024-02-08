const out = document.querySelector("output")

console.log(out)

async function getData () {
    let response = await fetch("http://localhost:8080/get-order?order-uid=order-JYuWVENk")
    console.log(await response.json())
}
getData()
