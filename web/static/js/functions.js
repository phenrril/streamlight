// Description: This file contains the
// functions that are used in the index.html file.


document.addEventListener('DOMContentLoaded', function() {
    var form = document.querySelector('form');
    if(form) {
        form.addEventListener('submit', function(event) {
            event.preventDefault();

            var data = {
                amount: document.getElementById('amount').value,
                customer_name: document.getElementById('customer_name').value,
                customer_email: document.getElementById('customer_email').value,
                order_id: document.getElementById('order_id').value
            };

            fetch('http://localhost:8080/pago', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            }).then(function(response) {
                if (!response.ok) {
                    throw new Error('HTTP error, status = ' + response.status);
                }
                window.location.href = '/exito';
            }).catch(function(error) {
                console.log(error);
            });
        });
    }
    console.log('functions.js loaded');
});