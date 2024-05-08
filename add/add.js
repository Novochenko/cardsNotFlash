const addForm = document.getElementById('question-form');
const submitBtn = document.getElementById('submit-btn');
const responseAdd = document.getElementById('response');

let questionId = 1; // начальный ID вопроса

submitBtn.addEventListener('click', (e) => {
  e.preventDefault();
  const question = document.getElementById('question').value;
  const answer = document.getElementById('answer').value;

  // Создаем объект вопроса и ответа
  const questionData = {
    id: questionId,
    question: question,
    answer: answer
  };

  // Создаем JSON-файл
  const jsonData = JSON.stringify(questionData);

  // Отправляем JSON-файл на сервер
  fetch('127.0.0.1:9000/createcard', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: jsonData
  })
  .then(response => response.json())
  .then(data => {
    responseDiv.innerHTML = `Вопрос отправлен успешно ID вопроса: ${questionId}`;
    questionId++; // увеличиваем ID вопроса
  })
  .catch(error => {
    responseDiv.innerHTML = `Ошибка отправки вопроса: ${error}`;
  });
});