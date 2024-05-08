// Получаем элементы формы
const form = document.getElementById('register-form');
const registerButton = document.getElementById('register-btn');
//const data = document.getElementById('register-response');


// Добавляем обработчик события на кнопку регистрации
registerButton.addEventListener('click', (e) => {
  e.preventDefault();
  const nickname = document.getElementById('nickname').value;
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  const userData = {
    nickname,
    email,
    password
  };

  // Отправляем запрос на сервер
  fetch("http://127.0.0.1:9000/users", {
    method: 'POST',
    //mode: 'no-cors',
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData)
  })
  .then((response) => response.text())
  .then((data) => {
    document.getElementById('register-response').innerHTML = data;
    if (data.success) {
      data.innerHTML = 'Пользователь зарегистрирован успешно!';
      window.location.href = "../main/main.html"
    } 
    else {
      data.innerHTML = 'Ошибка регистрации: ' + data.error;
    }
  })
  .catch(error => {
    console.error(error);
  });
});
