const loadBtn = document.getElementById('load-btn');
const jsonFilesSelect = document.getElementById('json-files');
const updateBtn = document.getElementById('update-btn');
const jsonDataDiv = document.getElementById('json-data');

loadBtn.addEventListener('click', loadJsonFiles);
updateBtn.addEventListener('click', updateJsonFile);

function loadJsonFiles() {
  fetch('https://localhost:8080/editcard')
    .then(response => response.json())
    .then(data => {
      const options = data.map(file => {
        return `<option value="${file}">${file}</option>`;
      });
      jsonFilesSelect.innerHTML = options.join('');
    })
    .catch(error => console.error('Error:', error));
}

function updateJsonFile() {
  const selectedFile = jsonFilesSelect.value;
  fetch(`https://localhost:8080/${selectedFile}`)
    .then(response => response.json())
    .then(data => {
      jsonDataDiv.innerHTML = `ID: ${data.id}<br>Вопрос: ${data.question}<br>Ответ: ${data.answer}`;
      // Создаем объект для хранения изменений
      const changes = {};
      // ...
      // Отправляем обновленный файл на сервер
      fetch(`https://localhost:8080/${selectedFile}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(changes)
      })
      .then(response => response.json())
      .then(data => {
        console.log('JSON-файл обновлен успешно!');
      })
      .catch(error => console.error('Error:', error));
    })
    .catch(error => console.error('Error:', error));
}