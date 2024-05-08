// Получаем элементы формы
const form = document.getElementById('register-form');
//const nicknameInput = document.getElementById("nickname");
//const emailInput = document.getElementById("email");
//const passwordInput = document.getElementById("password");
const registerButton = document.getElementById('register-btn');
const responseDiv = document.getElementById('register-response');


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
  .then((responseDiv) => responseDiv.json())
  .then((data) => {
    if (data.success) {
      responseDiv.innerHTML = 'Пользователь зарегистрирован успешно!';
      window.location.href = "../main/main.html"
    } else {
      responseDiv.innerHTML = 'Ошибка регистрации: ' + data.error;
    }
  })
  .catch((error) => {
    console.error(error);
  });
});


/*<span>Или используйте e-mail для регистрации</span>
<label for="nickname">Никнейм</label>
<input type="text" id="nickname" >
<label for="email">E-mail</label>
<input id="email" type="email" >
<label for="password">Пароль</label>
<input id="password" type="password">
<button id="register-btn">Создать</button>*/