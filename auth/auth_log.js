const loginForm = document.getElementById('login-form');
const loginButton = document.getElementById('login-btn');
const responseLogin = document.getElementById('login-response');
  
  
  // Добавляем обработчик события на кнопку регистрации
  loginButton.addEventListener('click', (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
  
    const userData = {
      email,
      password
    };
  
    // Отправляем запрос на сервер
    fetch("http://127.0.0.1:9000/sessions", {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(userData)
    })
    .then((response) => response.json())
    .then((data) => {
      if (data.success) {
        responseLogin.innerHTML = 'Пользователь вошел успешно!';
        window.location.href = "../main/main.html"
      } else {
        responseLogin.innerHTML = 'Ошибка регистрации: ' + data.error;
      }
    })
    .catch((error) => {
      console.error(error);
    });
  });

  /*function sendRequest() {
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
  }*/

  