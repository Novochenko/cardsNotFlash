// Функция для получения списка JSON файлов из базы данных
async function getJsonFiles() {
  try {
    const response = await fetch('localhost:8080/editcard'); // замените на URL вашей базы данных
    const data = await response.json();
    const options = data.map(item => `<option value="${item.id}">${item.question}</option>`);
    document.getElementById('json-select').innerHTML = options.join('');
  } catch (error) {
    console.error(error);
  }
}

// Функция для получения выбранного JSON файла
async function getSelectedJson() {
  const selectedId = document.getElementById('json-select').value;
  try {
    const response = await fetch(`localhost:8080/editcard/${selectedId}`); // замените на URL вашей базы данных
    const data = await response.json();
    document.getElementById('json-data').innerHTML = `
      <p>Вопрос: ${data.question}</p>
      <p>Ответ: ${data.answer}</p>
    `;
  } catch (error) {
    console.error(error);
  }
}

// Функция для редактирования выбранного JSON файла
async function editJsonData() {
  const selectedId = document.getElementById('json-select').value;
  const question = prompt("Введите новый вопрос:");
  const answer = prompt("Введите новый ответ:");
  try {
    const response = await fetch(`localhost:8080/editcard/${selectedId}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ question, answer })
    }); // замените на URL вашей базы данных
    if (response.ok) {
      alert("JSON файл обновлен успешно!");
    } else {
      alert("Ошибка обновления JSON файла");
    }
  } catch (error) {
    console.error(error);
  }
}

// Инициализация
//getJsonFiles();
document.getElementById('load-btn').addEventListener('click',getJsonFiles)
document.getElementById('edit-btn').addEventListener('click', editJsonData);
document.getElementById('json-select').addEventListener('change', getSelectedJson);