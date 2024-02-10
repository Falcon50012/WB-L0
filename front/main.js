const outputElement = document.querySelector('.output');
const form = document.querySelector('form');
const orderIDElem = document.getElementById('orderID')


form.addEventListener('submit', async function (event) {
    event.preventDefault();

    const orderID = document.getElementById('orderID').value;
    await getData(orderID);

    orderIDElem.value = '';
});

async function getData(orderID) {
    if (!orderID) {
        console.log('OrderID is required');
        return;
    }
    let response = await fetch(`http://localhost:8080/get-order?order-uid=${orderID}`);
    if (response.ok) {
        let data = await response.json();
        console.log(data);

        const formattedJsonString = JSON.stringify(data, null, 2);

        const safeHtml = formattedJsonString
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;")
            .replace(/ /g, '&nbsp;')
            .replace(/\n/g, '<br>');
        
        outputElement.innerHTML = `Order Data: <pre>${safeHtml}</pre>`;
    } else {
        console.error("HTTP-Error: " + response.status);
    }
}

