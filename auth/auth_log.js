function sendRequest() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/sessions', true);
    xhr.setRequestHeader('Content-Type', 'application/json');

    const data = JSON.stringify({ email, password });
    xhr.send(data);

    xhr.onload = function() {
      if (xhr.status < 300) {
        const response = JSON.parse(xhr.responseText);
        if (response.authenticated) {
          window.location.href = '/main.html';
        } else {
          alert('Invalid e-mail or password');
        }
      } else {
        alert('Error: ' + xhr.status);
      }
    };
  }

